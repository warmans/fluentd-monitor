package main

import (
    "github.com/gorilla/websocket"
    "log"
    "net/http"
    "flag"
    "os"
    "github.com/warmans/fluentd-api-client/monitoring"
)

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024, CheckOrigin: func(r *http.Request) bool { return true; }}

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

    monitor := NewMonitor(hub, []*monitoring.Host{monitoring.NewHost("127.0.0.1:24220")})
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
