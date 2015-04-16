package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"github.com/warmans/fluentd-api-client/monitoring"
)

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024, CheckOrigin: func(r *http.Request) bool { return true }}

//WebsocketHandler handles websocket connections
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
	client.StartListening()
}

func main() {

	// Setup configuration
	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")             //dev config (overrides live)
	viper.AddConfigPath("/etc/fluentd-monitor/") //live config
	viper.ReadInConfig()

	//start the hub
	hub := NewHub()
	go hub.Run()

	//convert raw host addresses to Host instances
	rawHosts := viper.GetStringSlice("hosts")
	hosts := make([]*monitoring.Host, len(rawHosts))
	for i, hostAddress := range rawHosts {
		hosts[i] = monitoring.NewHost(hostAddress)
	}

	//create the monitor
	monitor := NewMonitor(hub, hosts)
	//configure
	if viper.GetInt("history_size") > 0 {
		monitor.HistorySize = viper.GetInt("history_size")
	}
	if viper.GetInt("history_tick_seconds") > 0 {
		monitor.HistoryTickSeconds = viper.GetInt("history_tick_seconds")
	}
	if viper.GetInt("push_tick_seconds") > 0 {
		monitor.PushTickSeconds = viper.GetInt("push_tick_seconds")
	}
	//start monitoring
	go monitor.Run()

	//static assets
	var staticFileServer http.Handler
	if os.Getenv("DEV") == "true" {
		log.Print("Running in DEV mode")
		staticFileServer = http.FileServer(http.Dir("ui/static"))
	} else {
		//in production use embedded files
		staticFileServer = http.FileServer(FS(false))
	}
	http.Handle("/ui/", http.StripPrefix("/ui/", staticFileServer))

	//websocket
	http.Handle("/ws", &WebsocketHandler{Hub: hub})

	//redirect / to static dir
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			log.Print("redirect")
			http.Redirect(w, r, "/ui/", 301)
		}
	})

	log.Fatal("HTTP Server failed: ", http.ListenAndServe(viper.GetString("listen"), nil))
}
