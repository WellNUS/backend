package ws

import (
	"wellnus/backend/db/model"

	"bytes"
	"log"
	"net/http"
	"time"
	"encoding/json"
	"database/sql"

	"github.com/gorilla/websocket"
)

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
	Send 		chan interface{}
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
		payload, ok := <-c.Send
		if !ok {
			// The Hub closed the channel.
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		w, err := c.Conn.NextWriter(websocket.TextMessage)
		if err != nil { return }
		jpayload, err := json.Marshal(payload)
		if err != nil { return }
		w.Write(jpayload)

		if err := w.Close(); err != nil { return }
	}
}

func (c Client) UserName(db *sql.DB) (string, error) {
	user, err := model.GetUser(db, c.UserID)
	if err != nil { return "", err }
	return user.FirstName, nil
}

func (c Client) GroupName(db *sql.DB) (string, error) {
	group, err := model.GetGroup(db, c.GroupID)
	if err != nil { return "", err }
	return group.GroupName, nil
}

func ServeWs(Hub *Hub, w http.ResponseWriter, r *http.Request, userID int64, groupID int64) {
	Conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{ UserID: userID, GroupID: groupID, Hub: Hub, Conn: Conn, Send: make(chan interface{}, loadedMessageBuffer)}
	client.Hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}