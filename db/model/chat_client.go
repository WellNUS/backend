package model

import (
	"log"
	"encoding/json"
	"database/sql"

	"github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the Hub.
type Client struct {
	UserID		int64
	GroupID		int64
	Hub 		*Hub
	Conn 		*websocket.Conn
	Send 		chan interface{}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, jsentData, err := c.Conn.ReadMessage() // Read client's input field
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error when reading: %v", err)
			}
			break
		}
		var sentData SentData
		if err := json.Unmarshal(jsentData, &sentData); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error when unmarshalling: %v", err)
			}
			break
		}
		sentData.Client = c
		c.Hub.Broadcast <- sentData
	}
}

func (c *Client) WritePump() {
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
		jpayload, err := json.Marshal(payload)
		if err != nil { return }
		w, err := c.Conn.NextWriter(websocket.TextMessage)
		if err != nil { return }
		w.Write(jpayload)

		if err := w.Close(); err != nil { return }
	}
}

func (c Client) UserName(db *sql.DB) (string, error) {
	user, err := GetUser(db, c.UserID)
	if err != nil { return "", err }
	return user.FirstName, nil
}

func (c Client) GroupName(db *sql.DB) (string, error) {
	group, err := GetGroup(db, c.GroupID)
	if err != nil { return "", err }
	return group.GroupName, nil
}