package ws

import (
	"wellnus/backend/db/query"
	"database/sql"
	"fmt"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Connected DB.
	DB			*sql.DB

	// Registered clients.
	Clients 	map[*Client]bool

	// Inbound messages from the clients.
	Broadcast 	chan Message

	// Register requests from the clients.
	Register 	chan *Client

	// Unregister requests from clients.
	Unregister 	chan *Client
}

func NewHub(db *sql.DB) *Hub {
	return &Hub{
		DB:			db,
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func getRecipientIDSet(db *sql.DB, groupID int64) (map[int64]bool, error) {
	recipients, err := query.GetAllUsersOfGroup(db, groupID)
	if err != nil { return nil, err }
	recipientsSet := make(map[int64]bool)
	for _, user := range recipients {
		recipientsSet[user.ID] = true
	}
	return recipientsSet, nil
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			err := query.AddMessage(h.DB, message)
			if err != nil {
				fmt.Printf("An error occured during adding to database. %v \n", err)
				continue
			}
			loadedMessage, err := query.LoadMessage(h.DB, message)
			if err != nil {
				fmt.Printf("An error occured during loading. %v \n", err)
				continue
			}
			recipientsSet, err := getRecipientIDSet(h.DB, message.GroupID)
			if err != nil {
				fmt.Printf("An error occured while getting recipient set. %v \n", err)
				continue
			} 
			// Add message to database here
			for client := range h.Clients {
				if recipientsSet[client.UserID] {
					select {
						case client.Send <- loadedMessage:
						default:
							close(client.Send)
							delete(h.Clients, client)
					}
				}
			}
		}
	}
}