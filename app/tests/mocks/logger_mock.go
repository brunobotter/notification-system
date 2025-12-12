package mocks

import (
	"github.com/brunobotter/notification-system/infra/logger"
	"github.com/stretchr/testify/mock"
)

var _ logger.Logger = (*MockLogger)(nil)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) SetCommonField(commonFields map[string]any) {
	m.Called(commonFields)
}
func (m *MockLogger) InfoF(format string, args ...interface{}) {
	m.Called(format, args)
}
func (m *MockLogger) Info(args ...interface{}) {
	m.Called(args)
}
func (m *MockLogger) ErrorF(format string, args ...interface{}) {
	m.Called(format, args)
}
func (m *MockLogger) Error(format string, args ...interface{}) {
	m.Called(args)
}
func (m *MockLogger) Log(msg string) {
	m.Called(msg)
}
func (m *MockLogger) Sync() {
	m.Called()
}
