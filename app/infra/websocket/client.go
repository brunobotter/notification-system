package websocket

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client interface {
	Send(message []byte)
	Receive() ([]byte, error)
	Close() error
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

// Client representa uma conexão WebSocket com um usuário.
type clientImpl struct {
	hub  Hub
	conn *websocket.Conn
	send chan []byte
}

// NewClient cria uma nova instância de clientImpl.
func NewClient(conn *websocket.Conn, hub Hub) *clientImpl {
	return &clientImpl{
		conn: conn,
		hub:  hub,
		send: make(chan []byte, 256),
	}
}

// ReadPump escuta mensagens do WebSocket e repassa para o Hub.
func (c *clientImpl) ReadPump() {
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("erro inesperado: %v", err)
			}
			break
		}
		c.hub.Broadcast(message)
	}
}

// WritePump envia mensagens do canal 'send' para o WebSocket.
func (c *clientImpl) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// O canal foi fechado.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Close encerra a conexão WebSocket, remove o cliente do Hub e fecha o canal de envio.
func (c *clientImpl) Close() error {
	// Remove do Hub (caso ainda não tenha sido removido)
	c.hub.Unregister(c)
	// Fecha o canal de envio (WritePump irá encerrar)
	close(c.send)
	// Fecha a conexão WebSocket
	return c.conn.Close()
}

func (c *clientImpl) Send(message []byte) {
	select {
	case c.send <- message:
		// Mensagem enviada
	default:
		// Canal cheio ou cliente desconectado
		c.Close()
	}
}

// Receive lê uma mensagem do WebSocket.
func (c *clientImpl) Receive() ([]byte, error) {
	_, message, err := c.conn.ReadMessage()
	return message, err
}
