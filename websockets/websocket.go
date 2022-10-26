package websockets

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jettdc/switchboard/u"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WSConn struct {
	*websocket.Conn
	ReadLock   *sync.Mutex
	WriteLock  *sync.Mutex
	Ctx        context.Context
	CancelFunc context.CancelFunc
}

func HandleConnection(w http.ResponseWriter, r *http.Request) (*WSConn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	var readLock sync.Mutex
	var writeLock sync.Mutex

	ctx, cancelFunc := context.WithCancel(context.Background())

	return &WSConn{conn, &readLock, &writeLock, ctx, cancelFunc}, nil
}

// TODO: might not actually need locks on these...
func (w *WSConn) ReadMessageSafe() (messageType int, p []byte, err error) {
	w.ReadLock.Lock()
	defer w.ReadLock.Unlock()
	return w.ReadMessage()
}

func (w *WSConn) ReadJSONSafe(v interface{}) error {
	w.ReadLock.Lock()
	defer w.ReadLock.Unlock()
	return w.ReadJSON(v)
}

func (w *WSConn) WriteJSONSafe(j interface{}) {
	w.WriteLock.Lock()
	defer w.WriteLock.Unlock()
	if err := w.WriteJSON(j); err != nil {
		u.Logger.Error(fmt.Sprintf("Error writing json message to WS: %s", j))
	}
}

func (w *WSConn) CloseAndCancel() {
	u.Logger.Debug("Closing")
	if err := w.Close(); err != nil {
		u.Logger.Error(fmt.Sprintf("Error closing websocket connection: %s", err.Error()))
	}

	w.CancelFunc()
}

// CancelCtxOnReadErr should only be used on connections where no other reads are being made. No more than one goroutine
// may read from the connection at the same time, so this should be used strictly on write only channels. In other cases,
// context cancelling should be included with the message processing logic.
func (w *WSConn) CancelCtxOnReadErr() {
	go func() {
		for {
			msgType, _, err := w.ReadMessageSafe()
			if msgType == websocket.CloseMessage || err != nil {
				w.CloseAndCancel()
				return
			}
		}
	}()
}
