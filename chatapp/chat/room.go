package main

import (
	"ch1/trace"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

type room struct {
	// forward keep msg.
	forward chan *message
	// join is channel for joining client
	join chan *client
	// leave is channel for leaving client
	leave chan *client
	// clients is joined clients
	clients map[*client]bool
	// tracer receive log
	tracer trace.Tracer
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			log.Println("join")
			r.clients[client] = true
			r.tracer.Trace("new client joined")
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("client left")
		case msg := <-r.forward:
			r.tracer.Trace("message received: ", msg.Message)
			for client := range r.clients {
				select {
				case client.send <- msg:
					r.tracer.Trace(" -- sent to client")
				default:
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace(" -- Failed to send msg. cleanup client")
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("serveHTTP:", err)
		return
	}
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get a cookie")
		return
	}
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
