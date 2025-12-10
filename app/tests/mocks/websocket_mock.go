package mocks

import (
	"io"
	"time"

	"github.com/brunobotter/notification-system/infra/web_socket"
	"github.com/stretchr/testify/mock"
)

var _ web_socket.WebSocketConn = (*MockWebSocketConn)(nil)

type MockWebSocketConn struct {
	mock.Mock
	pongHandler func(string) error
}

func (m *MockWebSocketConn) ReadMessage() (int, []byte, error) {
	args := m.Called()
	return args.Int(0), args.Get(1).([]byte), args.Error(2)
}

func (m *MockWebSocketConn) WriteMessage(messageType int, data []byte) error {
	args := m.Called(messageType, data)
	return args.Error(0)
}

func (m *MockWebSocketConn) NextWriter(messageType int) (io.WriteCloser, error) {
	args := m.Called(messageType)
	return args.Get(0).(io.WriteCloser), args.Error(1)
}

func (m *MockWebSocketConn) SetReadLimit(limit int64) {
	m.Called(limit)
}

func (m *MockWebSocketConn) SetReadDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockWebSocketConn) SetWriteDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockWebSocketConn) SetPongHandler(h func(string) error) {
	m.Called(h)
	m.pongHandler = h
}

func (m *MockWebSocketConn) Close() error {
	args := m.Called()
	return args.Error(0)
}
