package web_socket

import (
	"sync"

	"github.com/brunobotter/notification-system/infra/logger"
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
	clients    map[*ClientImpl]bool
	broadcast  chan []byte
	register   chan *ClientImpl
	unregister chan *ClientImpl
	mu         sync.RWMutex
	logger     logger.Logger
}

// NewHub cria uma nova instância do Hub.
func NewHub(l logger.Logger) Hub {
	return &hubImpl{
		clients:    make(map[*ClientImpl]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *ClientImpl),
		unregister: make(chan *ClientImpl),
		logger:     l,
	}
}

// Run inicia o loop principal do Hub, gerenciando registro, remoção e broadcast.
func (h *hubImpl) Run() {
	h.logger.Info("hub iniciado")
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			h.logger.InfoF("Novo cliente registrado: %p", client)
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Sends)
				h.logger.InfoF("Cliente removido: %p", client)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.logger.InfoF("Broadcast de mensagem para %d clientes: %s", len(h.clients), string(message))

			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Sends <- message:
				default:
					close(client.Sends)
					delete(h.clients, client)
					h.logger.ErrorF("Canal cheio ou cliente desconectado: %p", client)

				}
			}
			h.mu.RUnlock()
		}
	}
}

// Register adiciona um novo cliente ao Hub.
func (h *hubImpl) Register(c Client) {
	if cli, ok := c.(*ClientImpl); ok {
		h.register <- cli
	}
}

// Unregister remove um cliente do Hub.
func (h *hubImpl) Unregister(c Client) {
	if cli, ok := c.(*ClientImpl); ok {
		h.unregister <- cli
	}
}

// Broadcast envia uma mensagem para todos os clientes conectados.
func (h *hubImpl) Broadcast(message []byte) {
	h.broadcast <- message
}
