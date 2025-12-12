package tests

import (
	"errors"
	"testing"
	"time"

	"github.com/brunobotter/notification-system/infra/web_socket"
	"github.com/brunobotter/notification-system/tests/fake"
	"github.com/brunobotter/notification-system/tests/mocks"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_ReadPump_MessageReceived(t *testing.T) {
	builder := mocks.NewSetup().WithConfig().WithLogger().WithHub().WithWebSocket().WithClient()

	msg := []byte("hello")
	builder.WebSocket.On("SetReadLimit", mock.Anything).Once()
	builder.WebSocket.On("SetReadDeadline", mock.Anything).Return(nil).Maybe()
	builder.WebSocket.On("SetPongHandler", mock.Anything).Once().Return()
	builder.WebSocket.On("ReadMessage").Return(1, msg, nil).Once()
	builder.WebSocket.On("ReadMessage").Return(0, ([]byte)(nil), errors.New("closed")).Once()
	builder.WebSocket.On("Close").Return(nil).Once()
	builder.Logger.On("InfoF", mock.Anything, mock.Anything)
	builder.Logger.On("ErrorF", mock.Anything, mock.Anything).Maybe()
	builder.Hub.On("Unregister", mock.Anything).Once()
	builder.Hub.On("Broadcast", msg).Once()

	client := web_socket.NewClient(builder.WebSocket, builder.Hub, builder.Logger)
	go client.ReadPump()
	time.Sleep(20 * time.Millisecond)

	builder.Hub.AssertCalled(t, "Broadcast", msg)
	builder.Hub.AssertCalled(t, "Unregister", mock.Anything)
	builder.WebSocket.AssertCalled(t, "Close")
	builder.Logger.AssertExpectations(t)
	builder.Hub.AssertExpectations(t)
	builder.WebSocket.AssertExpectations(t)
}

func TestClient_ReadPump_UnexpectedClose(t *testing.T) {
	builder := mocks.NewSetup().WithConfig().WithLogger().WithHub().WithWebSocket().WithClient()

	builder.WebSocket.On("SetReadLimit", mock.Anything).Once()
	builder.WebSocket.On("SetReadDeadline", mock.Anything).Return(nil).Maybe()
	builder.WebSocket.On("SetPongHandler", mock.Anything).Once()
	builder.WebSocket.On("ReadMessage").Return(0, ([]byte)(nil), errors.New("normal")).Once()
	builder.WebSocket.On("Close").Return(nil).Once()
	builder.Hub.On("Unregister", mock.Anything).Once()
	builder.Logger.On("InfoF", mock.Anything, mock.Anything)

	client := web_socket.NewClient(builder.WebSocket, builder.Hub, builder.Logger)
	go client.ReadPump()
	time.Sleep(20 * time.Millisecond)

	builder.Hub.AssertCalled(t, "Unregister", mock.Anything)
	builder.Logger.AssertCalled(t, "InfoF", mock.Anything, mock.Anything)
	builder.WebSocket.AssertCalled(t, "Close")
	builder.Logger.AssertExpectations(t)
	builder.Hub.AssertExpectations(t)
	builder.WebSocket.AssertExpectations(t)
}

func TestClient_ReadPump_NormalClose(t *testing.T) {
	builder := mocks.NewSetup().WithConfig().WithLogger().WithHub().WithWebSocket().WithClient()

	builder.WebSocket.On("SetReadLimit", mock.Anything).Once()
	builder.WebSocket.On("SetReadDeadline", mock.Anything).Return(nil).Maybe()
	builder.WebSocket.On("SetPongHandler", mock.Anything).Once()
	builder.WebSocket.On("ReadMessage").Return(0, ([]byte)(nil), errors.New("normal")).Once()
	builder.WebSocket.On("Close").Return(nil).Once()
	builder.Hub.On("Unregister", mock.Anything).Once()
	builder.Logger.On("InfoF", mock.Anything, mock.Anything)

	client := web_socket.NewClient(builder.WebSocket, builder.Hub, builder.Logger)
	go client.ReadPump()
	time.Sleep(20 * time.Millisecond)

	builder.Hub.AssertCalled(t, "Unregister", mock.Anything)
	builder.Logger.AssertCalled(t, "InfoF", mock.Anything, mock.Anything)
	builder.WebSocket.AssertCalled(t, "Close")
	builder.Logger.AssertExpectations(t)
	builder.Hub.AssertExpectations(t)
	builder.WebSocket.AssertExpectations(t)
}

func TestClient_WritePump_NextWriterError(t *testing.T) {
	builder := mocks.NewSetup().WithConfig().WithLogger().WithHub().WithWebSocket().WithClient()

	builder.WebSocket.On("SetWriteDeadline", mock.Anything).Once().Return(nil)
	builder.WebSocket.On("NextWriter", websocket.TextMessage).Return((*fake.DummyWriteCloser)(nil), errors.New("fail"))
	builder.WebSocket.On("Close").Return(nil).Once()
	builder.Logger.On("InfoF", mock.Anything, mock.Anything).Maybe()

	client := web_socket.NewClient(builder.WebSocket, builder.Hub, builder.Logger)
	client.Send([]byte("abc"))
	go client.WritePump()
	time.Sleep(20 * time.Millisecond)

	builder.WebSocket.AssertCalled(t, "NextWriter", websocket.TextMessage)
	builder.WebSocket.AssertExpectations(t)
}

func TestClient_WritePump_PingPeriod(t *testing.T) {
	builder := mocks.NewSetup().WithConfig().WithLogger().WithHub().WithWebSocket().WithClient()
	builder.WebSocket.On("SetWriteDeadline", mock.Anything).Return(nil)
	builder.WebSocket.On("WriteMessage", websocket.PingMessage, []byte(nil)).Return(nil)
	builder.WebSocket.On("WriteMessage", websocket.CloseMessage, mock.Anything).Return(nil)
	builder.WebSocket.On("Close").Return(nil)
	builder.Logger.On("InfoF", mock.Anything, mock.Anything).Maybe()
	client := web_socket.NewClient(builder.WebSocket, builder.Hub, builder.Logger, func(c *web_socket.ClientImpl) {
		c.PingPeriod = 10 * time.Millisecond
	})
	go client.WritePump()
	time.Sleep(50 * time.Millisecond)
	client.CloseSendChannel()
	time.Sleep(10 * time.Millisecond)
	builder.WebSocket.AssertCalled(t, "WriteMessage", websocket.PingMessage, []byte(nil))
	builder.WebSocket.AssertCalled(t, "Close")
	builder.WebSocket.AssertExpectations(t)
}

func TestClient_Receive_Ok(t *testing.T) {
	builder := mocks.NewSetup().WithConfig().WithLogger().WithHub().WithWebSocket().WithClient()
	builder.WebSocket.On("ReadMessage").Return(1, []byte("msg"), nil)
	client := web_socket.NewClient(builder.WebSocket, builder.Hub, builder.Logger)
	msg, err := client.Receive()
	assert.NoError(t, err)
	assert.Equal(t, []byte("msg"), msg)
}

func TestClient_Receive_Error(t *testing.T) {
	builder := mocks.NewSetup().WithConfig().WithLogger().WithHub().WithWebSocket().WithClient()
	builder.WebSocket.On("ReadMessage").Return(0, []byte(nil), errors.New("fail"))
	client := web_socket.NewClient(builder.WebSocket, builder.Hub, builder.Logger)
	msg, err := client.Receive()
	assert.Error(t, err)
	assert.Nil(t, msg)
}

func TestClient_Close(t *testing.T) {
	builder := mocks.NewSetup().WithConfig().WithLogger().WithHub().WithWebSocket().WithClient()
	builder.Hub.On("Unregister", mock.Anything).Once()
	builder.WebSocket.On("Close").Return(nil).Once()

	client := web_socket.NewClient(builder.WebSocket, builder.Hub, builder.Logger)
	err := client.Close()
	assert.NoError(t, err)
	builder.Hub.AssertCalled(t, "Unregister", client)
	builder.WebSocket.AssertCalled(t, "Close")
}
func TestClient_Send_ChannelFull(t *testing.T) {
	builder := mocks.NewSetup().WithConfig().WithLogger().WithHub().WithWebSocket().WithClient()
	builder.Hub.On("Unregister", mock.Anything).Once()
	builder.WebSocket.On("Close").Return(nil).Once()

	client := web_socket.NewClient(builder.WebSocket, builder.Hub, builder.Logger)
	// Preenche o canal
	for i := 0; i < cap(client.Sends); i++ {
		client.Sends <- []byte("msg")
	}
	client.Send([]byte("overflow"))
	// Espera o Close ser chamado
	time.Sleep(10 * time.Millisecond)
	builder.WebSocket.AssertCalled(t, "Close")
}

func TestClient_WritePump_MultipleMessages(t *testing.T) {
	builder := mocks.NewSetup().WithConfig().WithLogger().WithHub().WithWebSocket().WithClient()
	writer := &fake.DummyWriteCloser{}
	builder.WebSocket.On("SetWriteDeadline", mock.Anything).Return(nil)
	builder.WebSocket.On("NextWriter", websocket.TextMessage).Return(writer, nil)
	builder.WebSocket.On("Close").Return(nil)
	builder.Logger.On("InfoF", mock.Anything, mock.Anything).Maybe()

	client := web_socket.NewClient(builder.WebSocket, builder.Hub, builder.Logger)
	client.Send([]byte("msg1"))
	client.Send([]byte("msg2"))
	go client.WritePump()
	time.Sleep(20 * time.Millisecond)
	builder.WebSocket.AssertNumberOfCalls(t, "NextWriter", 2)
}
