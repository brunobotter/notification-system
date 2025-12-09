package websocket

import (
	"log"
	"sync"
)

type Hub interface {
	Run()
	Register(client Client)
	Unregister(client Client)
	Broadcast(message []byte)
}

// Hub é a implementação concreta da interface Hub.
// Ele gerencia todos os clientes conectados e faz o broadcast das mensagens.
type hubImpl struct {
	clients    map[*clientImpl]bool
	broadcast  chan []byte
	register   chan *clientImpl
	unregister chan *clientImpl
	mu         sync.RWMutex
}

// NewHub cria uma nova instância do Hub.
func NewHub() Hub {
	return &hubImpl{
		clients:    make(map[*clientImpl]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *clientImpl),
		unregister: make(chan *clientImpl),
	}
}

// Run inicia o loop principal do Hub, gerenciando registro, remoção e broadcast.
func (h *hubImpl) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Println("Novo cliente registrado")
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Println("Cliente removido")
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Register adiciona um novo cliente ao Hub.
func (h *hubImpl) Register(c Client) {
	if cli, ok := c.(*clientImpl); ok {
		h.register <- cli
	}
}

// Unregister remove um cliente do Hub.
func (h *hubImpl) Unregister(c Client) {
	if cli, ok := c.(*clientImpl); ok {
		h.unregister <- cli
	}
}

// Broadcast envia uma mensagem para todos os clientes conectados.
func (h *hubImpl) Broadcast(message []byte) {
	h.broadcast <- message
}
