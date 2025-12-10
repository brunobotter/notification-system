package tests

import (
	"testing"
	"time"

	"github.com/brunobotter/notification-system/infra/web_socket"
	"github.com/brunobotter/notification-system/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHub_RegisterAndUnregister(t *testing.T) {
	builder := mocks.NewSetup().WithConfig().WithLogger().WithWebSocket().WithHub().WithClient()
	builder.Logger.On("Info", mock.Anything)
	builder.Logger.On("InfoF", mock.Anything, mock.Anything)

	hub := web_socket.NewHub(builder.Logger)
	go hub.Run()

	client := web_socket.NewClient(builder.WebSocket, hub, builder.Logger)
	hub.Register(client)
	time.Sleep(10 * time.Millisecond)
	hub.Unregister(client)
	time.Sleep(10 * time.Millisecond)

	builder.Logger.AssertCalled(t, "InfoF", mock.Anything, mock.Anything)
}

func TestHub_BroadcastToClients(t *testing.T) {
	builder := mocks.NewSetup().WithConfig().WithLogger().WithWebSocket().WithHub()
	builder.Logger.On("Info", mock.Anything)
	builder.Logger.On("InfoF", mock.Anything, mock.Anything).Maybe()
	builder.Logger.On("ErrorF", mock.Anything, mock.Anything).Maybe()

	hub := web_socket.NewHub(builder.Logger)
	go hub.Run()

	// Cria dois clients reais usando o builder
	client1 := web_socket.NewClient(builder.WithWebSocket().WebSocket, hub, builder.Logger)
	client2 := web_socket.NewClient(builder.WithWebSocket().WebSocket, hub, builder.Logger)

	hub.Register(client1)
	hub.Register(client2)
	time.Sleep(10 * time.Millisecond)

	msg := []byte("hello")
	hub.Broadcast(msg)
	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, msg, <-client1.Sends)
	assert.Equal(t, msg, <-client2.Sends)
}
func TestHub_Broadcast_ChannelFull_RemovesClient(t *testing.T) {
	builder := mocks.NewSetup().WithConfig().WithLogger().WithWebSocket().WithHub()
	builder.Logger.On("Info", mock.Anything)
	builder.Logger.On("InfoF", mock.Anything, mock.Anything)
	builder.Logger.On("ErrorF", mock.Anything, mock.Anything)
	hub := web_socket.NewHub(builder.Logger)
	go hub.Run()

	client := web_socket.NewClient(builder.WebSocket, hub, builder.Logger)
	hub.Register(client)
	time.Sleep(20 * time.Millisecond) // Aguarda registro

	// Esvazie o canal antes de testar fechamento
	for i := 0; i < cap(client.Sends); i++ {
		client.Sends <- []byte("preenchido")
	}
	hub.Broadcast([]byte("overflow"))
	time.Sleep(100 * time.Millisecond) // Aguarde processamento

	// Esvazie o buffer antes de verificar se está fechado
	for i := 0; i < cap(client.Sends); i++ {
		<-client.Sends
	}
	_, ok := <-client.Sends
	assert.False(t, ok) // Canal deve estar fechado

	builder.Logger.AssertNumberOfCalls(t, "ErrorF", 1)
}
