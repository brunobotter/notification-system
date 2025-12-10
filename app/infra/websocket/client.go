package websocket

import (
	"time"

	"github.com/brunobotter/notification-system/infra/logger"
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
	hub    Hub
	conn   *websocket.Conn
	send   chan []byte
	logger logger.Logger
}

// NewClient cria uma nova instância de clientImpl.
func NewClient(conn *websocket.Conn, hub Hub, logger logger.Logger) *clientImpl {
	return &clientImpl{
		conn:   conn,
		hub:    hub,
		send:   make(chan []byte, 256),
		logger: logger,
	}
}

// ReadPump escuta mensagens do WebSocket e repassa para o Hub.
func (c *clientImpl) ReadPump() {
	c.logger.InfoF("ReadPump iniciado para cliente: %p", c)

	defer func() {
		c.hub.Unregister(c)
		c.conn.Close()
		c.logger.InfoF("ReadPump encerrado para cliente: %p", c)

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
				c.logger.ErrorF("Erro inesperado na conexão do cliente %p: %v", c, err)
			} else {
				c.logger.InfoF("Conexão fechada pelo cliente %p: %v", c, err)
			}
			break
		}
		c.logger.InfoF("Mensagem recebida do cliente %p: %s", c, string(message))

		c.hub.Broadcast(message)
	}
}

// WritePump envia mensagens do canal 'send' para o WebSocket.
func (c *clientImpl) WritePump() {
	c.logger.InfoF("WritePump iniciado para cliente: %p", c)

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
		c.logger.InfoF("WritePump encerrado para cliente: %p", c)

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
