package mocks

import (
	"github.com/brunobotter/notification-system/infra/web_socket"
	"github.com/stretchr/testify/mock"
)

type MockHub struct {
	mock.Mock
}

func (m *MockHub) Run() {
	m.Called()
}
func (m *MockHub) Register(client web_socket.Client) {
	m.Called(client)
}
func (m *MockHub) Unregister(client web_socket.Client) {
	m.Called(client)
}
func (m *MockHub) Broadcast(message []byte) {
	m.Called(message)
}
