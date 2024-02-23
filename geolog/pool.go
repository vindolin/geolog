package main

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// WsPool is a collection of websocket connections

type WsPool struct {
	mu        sync.Mutex
	clients   map[*websocket.Conn]bool
	broadcast chan string
}

func NewWsPool() *WsPool {
	return &WsPool{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan string),
	}
}

func (p *WsPool) Add(conn *websocket.Conn) {
	p.mu.Lock()
	p.clients[conn] = true
	p.mu.Unlock()
}

func (p *WsPool) Remove(conn *websocket.Conn) {
	p.mu.Lock()
	delete(p.clients, conn)
	p.mu.Unlock()
}

func (p *WsPool) Broadcast(ip string) {
	p.broadcast <- ip
}

func (p *WsPool) Start() {
	for {
		ip := <-p.broadcast
		for client := range p.clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(ip))
			if err != nil {
				log.Println(err)
				p.Remove(client)
			}
		}
	}
}
