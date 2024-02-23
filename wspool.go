package main

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// WsPool is a collection of websocket connections
type WsPool struct {
	clients   sync.Map
	broadcast chan string
}

func NewWsPool() *WsPool {
	return &WsPool{
		broadcast: make(chan string),
	}
}

func (p *WsPool) Add(conn *websocket.Conn) {
	p.clients.Store(conn, true)
}

func (p *WsPool) Remove(conn *websocket.Conn) {
	p.clients.Delete(conn)
}

func (p *WsPool) Broadcast(ip string) {
	p.broadcast <- ip
}

func (p *WsPool) Start() {
	for {
		ip := <-p.broadcast
		p.clients.Range(func(client, _ interface{}) bool {
			err := client.(*websocket.Conn).WriteMessage(websocket.TextMessage, []byte(ip))
			if err != nil {
				log.Println(err)
				p.Remove(client.(*websocket.Conn))
			}
			return true
		})
	}
}
