package mocks

import "github.com/brunobotter/notification-system/main/config"

type Setup struct {
	Config    *config.Config
	Logger    *MockLogger
	Hub       *MockHub
	WebSocket *MockWebSocketConn
	Client    *MockClient
}

func NewSetup() *Setup {
	return &Setup{}
}

func (s *Setup) WithConfig() *Setup {
	s.Config = &config.Config{}
	return s
}

func (s *Setup) WithLogger() *Setup {
	s.Logger = &MockLogger{}
	return s
}

func (s *Setup) WithHub() *Setup {
	s.Hub = &MockHub{}
	return s
}

func (s *Setup) WithWebSocket() *Setup {
	s.WebSocket = &MockWebSocketConn{}
	return s
}

func (s *Setup) WithClient() *Setup {
	s.Client = &MockClient{}
	return s
}
