package ws

import (
	"wellnus/backend/model"

	"bytes"
	"log"
	"net/http"
	"time"
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Message = model.Message
type LoadedMessage = model.LoadedMessage

const (
	loadedMessageBuffer = 256
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the Hub.
type Client struct {
	UserID		int64
	GroupID		int64
	Hub 		*Hub
	Conn 		*websocket.Conn
	Send 		chan LoadedMessage
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, msg, err := c.Conn.ReadMessage() // Read client's input field
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
		message := Message{ UserID: c.UserID, GroupID: c.GroupID, TimeAdded: time.Now(), Msg: string(msg) }
		c.Hub.Broadcast <- message
	}
}

func (c *Client) writePump() {
	defer func() {
		c.Conn.Close()
	}()
	for {
		loadedMessage, ok := <-c.Send
		if !ok {
			// The Hub closed the channel.
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		w, err := c.Conn.NextWriter(websocket.TextMessage)
		if err != nil { return }
		jLoadedMessage, err := json.Marshal(loadedMessage)
		if err != nil { return }
		w.Write(jLoadedMessage)

		if err := w.Close(); err != nil { return }
	}
}

func ServeWs(Hub *Hub, w http.ResponseWriter, r *http.Request, userID int64, groupID int64) {
	Conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{ UserID: userID, GroupID: groupID, Hub: Hub, Conn: Conn, Send: make(chan LoadedMessage, loadedMessageBuffer)}
	client.Hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}