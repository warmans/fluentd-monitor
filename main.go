package main

import (
    "github.com/gorilla/websocket"
    "log"
    "net/http"
    "flag"
    "time"
    "encoding/json"
    "sync"
    "bytes"
    "text/template"
    "io/ioutil"
    "os"
)

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024, CheckOrigin: func(r *http.Request) bool { return true; }}

//-----------------------------------------------
// Hub
//-----------------------------------------------

type Hub struct {
    clients     map[*Client]bool
    broadcast   chan []byte
    register    chan *Client
    unregister  chan *Client
}

func (h *Hub) Register(conn *Client) {
    h.register <- conn
}

func (h *Hub) Unregister(client *Client) {
    h.unregister <- client
}

func (h *Hub) Broadcast(msg []byte) {
    h.broadcast <- msg
}

func (h *Hub) Run() {
    for {
        select {

        //new client
        case conn := <-h.register:
            h.clients[conn] = true
            log.Printf("client joined (%v active)", len(h.clients))

        //client disconnects
        case conn := <-h.unregister:
            if _, ok := h.clients[conn]; ok {
                //do remove
                delete(h.clients, conn)
                close(conn.sendQueue)
                log.Printf("client left (%v active)", len(h.clients))
            }

        //new message
        case msg := <-h.broadcast:
            for client := range h.clients {
                client.EnqueueMessage(msg)
            }
        }
    }
}

func NewHub() *Hub {
    return &Hub{
        clients:     make(map[*Client]bool),
        broadcast:   make(chan []byte, 64),
        register:    make(chan *Client),
        unregister:  make(chan *Client),
    }
}

//-----------------------------------------------
// Client
//-----------------------------------------------

type Client struct {
    conn *websocket.Conn
    sendQueue chan []byte
    hub *Hub
}

func (client *Client) JoinHub(hub *Hub) {
    client.hub = hub
    client.hub.Register(client)
}

func (client *Client) LeaveHub() {
    if (client.hub != nil) {
        client.hub.Unregister(client)
    }
}

func (client *Client) Close() {
    client.LeaveHub()
    client.conn.Close()
}

func (client *Client) EnqueueMessage(msg []byte) {
    select {
    case  client.sendQueue <- msg:
    default:
        log.Print("client send buffer is full - they're probably dead so they've been disconnected")
        client.Close()
    }
}

func (client *Client) StartSending() {

    defer client.Close()

    for msg := range client.sendQueue {
        if err := client.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
            break
        }
    }
}

// AcceptMessages accepts messages from client and re-broadcasts to other clients
func (client *Client) StartListening() {

    defer client.Close()

    for {
        _, _, err := client.conn.ReadMessage()
        if err != nil {
            break
        }
        //don't do anything with messages for now
    }
}


//-----------------------------------------------
// Fuentd Status Monitor
//-----------------------------------------------

type FMonitor struct {
    Hub         *Hub
    Hosts       []*Host
}

func (mon *FMonitor) Run() {
    updateTicker := time.NewTicker(time.Second * 10)
    for range updateTicker.C {

        for _, host := range mon.Hosts {
            host.UpdatePlugins()
        }

        //Create a JSON payload
        buffer := &bytes.Buffer{}
        encoder := json.NewEncoder(buffer)

        if err := encoder.Encode(mon.Hosts); err != nil {
            log.Print("Failed to encode host stats payload");
            continue
        }

        //Broadcast to all clients
        mon.Hub.Broadcast(buffer.Bytes())
    }
}

type Plugins struct {
    Timestamp   int64
    Plugins     []Plugin    `json:"plugins"`
}

type Plugin struct {
    PluginId                string  `json:"plugin_id"`
    PluginCategory          string  `json:"plugin_category"`
    Type                    string  `json:"type"`
    OutputPlugin            bool    `json:"output_plugin"`
    BufferQueueLength       int     `json:"buffer_queue_length"`
    BufferTotalQueuedSize   int     `json:"buffer_total_queued_size"`
    RetryCount              int     `json:"retry_count"`
}

type Host struct {
    Address         string
    Online          bool
    LastError       string
    PluginHistory   []Plugins
    UpdateLock      sync.RWMutex
}

func (h *Host) HandleUpdateError(err error) {
    log.Printf("Error querying host %s: %s", h.Address, err.Error())
    h.Online = false
    h.LastError = err.Error()
}

func (h *Host) ClearUpdateError() {
    h.Online = true
    h.LastError = ""
}

func (h *Host) UpdatePlugins() {

    response, err := http.Get("http://"+h.Address+"/api/plugins.json");

    h.UpdateLock.Lock()

    if err != nil {
        h.HandleUpdateError(err)
        h.UpdateLock.Unlock()
        return
    }

    plugins := Plugins{Timestamp: time.Now().Unix()}
    decoder := json.NewDecoder(response.Body)
    if err := decoder.Decode(&plugins); err != nil {
        h.HandleUpdateError(err)
        h.UpdateLock.Unlock()
        return
    }

    h.ClearUpdateError()

    //trim history until it's the correct size
    for len(h.PluginHistory) >= 60 {
        h.PluginHistory = h.PluginHistory[1:]
    }
    //append new value
    h.PluginHistory = append(h.PluginHistory, plugins)

    h.UpdateLock.Unlock()
}

func (h *Host) GetPluginHistory() []Plugins {

    h.UpdateLock.RLock()
    history := h.PluginHistory
    h.UpdateLock.RUnlock()
    return history
}

func NewHost(Address string) *Host {
    return &Host{Address: Address, Online: false, PluginHistory: make([]Plugins, 0)}
}

//-----------------------------------------------
// HTTP Handler
//-----------------------------------------------

type WebsocketHandler struct {
    Hub *Hub
}

func (handler *WebsocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print("socket upgrade failed: ", err)
        return
    }

    //open new client connection
    client := &Client{sendQueue: make(chan []byte, 256), conn: conn}

    //join the hub
    client.JoinHub(handler.Hub)

    //unregister on connection close
    defer client.Close()

    //start relaying messages to client
    go client.StartSending()

    //start listening for client messages
    client.StartListening();
}

func InDevMode() bool {
    if os.Getenv("DEV") == "true" {
        return true
    }
    return false
}

func main() {
    flag.Parse()

    if InDevMode() {
        log.Print("Running in DEV mode")
    }

    hub := NewHub()
    go hub.Run()

    monitor := FMonitor{Hub: hub, Hosts: []*Host{NewHost("127.0.0.1:24220")}}
    go monitor.Run()

    //static assets
    var staticFileServer http.Handler
    if (InDevMode()) {
        //in dev mode serve raw files
        staticFileServer = http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
    } else {
        //in production use embedded files
        staticFileServer = http.FileServer(FS(false))
    }
    http.Handle("/static/", staticFileServer)

    //websocket
    http.Handle("/ws", &WebsocketHandler{Hub: hub})

    //redirect / to static dir
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/" {
            log.Print("redirect")
            http.Redirect(w, r, "/static/", 301)
        }
    });

    if err := http.ListenAndServe(*flag.String("addr", ":8080", "http service address"), nil); err != nil {
        log.Fatal("HTTP Server failed: ", err)
    }
}
