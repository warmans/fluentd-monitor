package main

import (
	"github.com/gorilla/websocket"
	"log"
)

//Hub connects clients
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

//Register adds a client to the hub
func (h *Hub) Register(conn *Client) {
	h.register <- conn
}

//Unregister removes a client from the hub
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

//Broadcast sends a message to all hub clients
func (h *Hub) Broadcast(msg []byte) {
	h.broadcast <- msg
}

//Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		//new client
		case conn := <-h.register:
			h.clients[conn] = true

		//client disconnects
		case conn := <-h.unregister:
			if _, ok := h.clients[conn]; ok {
				//do remove
				delete(h.clients, conn)
				close(conn.sendQueue)
			}

		//new message
		case msg := <-h.broadcast:
			for client := range h.clients {
				client.EnqueueMessage(msg)
			}
		}
	}
}

//NewHub create a new hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 64),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

//Client represents a connected instance of the front end
type Client struct {
	conn      *websocket.Conn
	sendQueue chan []byte
	hub       *Hub
}

func (client *Client) JoinHub(hub *Hub) {
	client.hub = hub
	client.hub.Register(client)
}

func (client *Client) LeaveHub() {
	if client.hub != nil {
		client.hub.Unregister(client)
	}
}

func (client *Client) Close() {
	client.LeaveHub()
	client.conn.Close()
}

func (client *Client) EnqueueMessage(msg []byte) {
	select {
	case client.sendQueue <- msg:
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
