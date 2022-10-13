package websockets

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	// Close the ws connection when the request finishes or fails
	go func() {
		select {
		case <-r.Context().Done():
			//conn.Close()
		}
	}()

	return conn, nil
}
