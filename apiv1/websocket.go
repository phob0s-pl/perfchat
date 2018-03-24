package api

import (
	"time"

	"encoding/json"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type websocketClient struct {
	conn *websocket.Conn
	msg  chan *Message
	name string
}

// StartWebsocket starts receiving messages from websocket
func (a *API) StartWebsocket() {
	go func() {
		for {
			select {
			case msg := <-a.message:
				_, err := a.Engine.GetUserByName(msg.User)
				if err != nil {
					continue
				}
				room, err := a.Engine.GetRoomByName(msg.Room)
				if err != nil {
					continue
				}

				for _, roomUser := range room.Users {
					ws, ok := a.websocketClients[roomUser.Name]
					if !ok {
						continue
					}
					ws.msg <- msg
				}
			} // <- really nice bracketception xDDD
		}
	}()
}

func (a *API) writeClientMessage(client *websocketClient) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.msg:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			payload, err := json.Marshal(message)
			if err == nil {
				w.Write(payload)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}

}

func (a *API) readClientMessage(client *websocketClient) {
	defer func() {
		delete(a.websocketClients, client.name)
		client.conn.Close()
	}()
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			break
		}
		msg := &Message{}
		if err := json.Unmarshal(message, msg); err == nil {
			a.message <- msg
		}
	}
}
