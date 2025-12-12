package web_socket

import (
	"io"
	"time"
)

type WebSocketConn interface {
	ReadMessage() (int, []byte, error)
	WriteMessage(messageType int, data []byte) error
	NextWriter(messageType int) (io.WriteCloser, error)
	SetReadLimit(limit int64)
	SetReadDeadline(t time.Time) error
	SetPongHandler(h func(string) error)
	Close() error
	SetWriteDeadline(t time.Time) error
}
