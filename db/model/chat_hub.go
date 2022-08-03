package model

import (
	"database/sql"
	"fmt"
	"time"
	"sort"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Connected DB.
	DB			*sql.DB

	// Registered clients.
	Clients 	map[*Client]bool

	// Inbound messages from the clients.
	Broadcast 	chan SentData

	// Register requests from the clients.
	Register 	chan *Client

	// Unregister requests from clients.
	Unregister 	chan *Client
}

func NewHub(db *sql.DB) *Hub {
	return &Hub{
		DB:			db,
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan SentData),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Helper functions
func (h *Hub) ChatStatusPayload(group Group) (ChatStatusPayload, error) {
	usersInGroup, err := GetAllUsersOfGroup(h.DB, group.ID)
	if err != nil { return ChatStatusPayload{}, err }

	usersInGroupMap := make(map[int64]User)
	for _, user := range usersInGroup {
		usersInGroupMap[user.ID] = user
	}

	MakeLess := func(users []User) func(int, int) bool {
		return func(i, j int) bool {
			return users[i].ID < users[j].ID
		}
	}
	inChatMembers := make([]User, 0) 
	onlineMembers := make([]User, 0)	
	offlineMembers := make([]User, 0)
	fmt.Printf("Websocket Clients UserID: [")
	for client := range h.Clients {
		fmt.Printf("%d, ", client.UserID)
		user, ok := usersInGroupMap[client.UserID]
		if ok {
			if client.GroupID == group.ID {
				inChatMembers = append(inChatMembers, user)
			} else {
				onlineMembers = append(onlineMembers, user)
			}
			delete(usersInGroupMap, client.UserID)
		}
	}
	fmt.Println("]")
	for _, user  := range usersInGroupMap {
		offlineMembers = append(offlineMembers, user)
	}
	sort.Slice(inChatMembers, MakeLess(inChatMembers))
	sort.Slice(onlineMembers, MakeLess(onlineMembers))
	sort.Slice(offlineMembers, MakeLess(offlineMembers))

	return ChatStatusPayload{
		Tag: ChatStatusTag, 
		GroupID: group.ID,
		GroupName: group.GroupName, 
		SortedInChatMembers: inChatMembers,
		SortedOnlineMembers: onlineMembers,
		SortedOfflineMembers: offlineMembers,
	}, nil
}

// Members are in only 1 of 3 states (in chat, online or offline)
// inChat means the member is on the given chat page
// online means the member is connected but on some other chat page
// offline means the member is not connected
func (h *Hub) SendOutToGroup(groupID int64, payload interface{}, cFilter func(*Client)bool) error {
	recipients, err := GetAllUsersOfGroup(h.DB, groupID)
	if err != nil { return err }
	recipientsMap := make(map[int64]bool)
	for _, user := range recipients {
		recipientsMap[user.ID] = true
	}
	if err != nil { return err }
	for client := range h.Clients {
		if recipientsMap[client.UserID] {
			if !cFilter(client) { continue }
			select {
				case client.Send <- payload:
				default:
					close(client.Send)
					delete(h.Clients, client)
			}
		}
	}
	return nil
}

func (h *Hub) SendOutChatStatus(userID int64) error {
	// userID is of user that induce the change in chat status
	groups, err := GetAllGroupsOfUser(h.DB, userID)
	if err != nil { return err }
	for _, group := range groups {
		chatStatusPayload, err := h.ChatStatusPayload(group)
		if err != nil { return err }
		err = h.SendOutToGroup(group.ID, chatStatusPayload, func(c *Client)bool{
			return c.GroupID == group.ID
		});
		if err != nil { return err }
	}
	return nil
}

func (h *Hub) SendOutServerMessage(groupID int64, serverMsg string) error {
	group, err := GetGroup(h.DB, groupID)
	if err != nil { return err }
	serverMessagePayload := MessagePayload{
		Tag: MessageTag,
		SenderName: "[WellNUS Server]",
		GroupName: group.GroupName,
		Message: Message{
			UserID: -1,
			GroupID: groupID,
			TimeAdded: time.Now(),
			Msg: serverMsg,
		},
	}
	err = h.SendOutToGroup(groupID, serverMessagePayload, func(c *Client)bool{
		return c.GroupID == group.ID
	});
	if err != nil { return err }
	return nil
}

// Main functions
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.HandleRegister(client)
		case client := <-h.Unregister:
			h.HandleUnregister(client)
		case sentData := <-h.Broadcast:
			switch sentData.Tag {
			case MessageTag:
				message, err := sentData.ToMessage()
				if err != nil {
					fmt.Printf("An error occured while converting sentData to Message. %v \n", err)
				}
				h.HandleMessage(message)
			default: 
				fmt.Printf("Unrecognised tag sent through websocket connection. Tag = %d \n", sentData.Tag)
			}
		}
	}
}

func (h *Hub) HandleRegister(client *Client) {
	h.Clients[client] = true

	err := h.SendOutChatStatus(client.UserID)
	if err != nil {
		fmt.Printf("An error occured during sending chat status payload. %v \n", err)
		return
	}

	clientName, err := client.UserName(h.DB)
	if err != nil {
		fmt.Printf("An error occured during retrieving first name of client. %v \n", err)
		return
	}
	
	err = h.SendOutServerMessage(client.GroupID, fmt.Sprintf("%s has joined the chat.", clientName))
	if err != nil {
		fmt.Printf("An error occured while sending out server message. %v \n", err)
		return
	}
}

func (h *Hub) HandleUnregister(client *Client) {
	if _, ok := h.Clients[client]; !ok { return }
	delete(h.Clients, client)
	close(client.Send)

	err := h.SendOutChatStatus(client.UserID)
	if err != nil {
		fmt.Printf("An error occured during sending chat status payload. %v \n", err)
		return
	}
	
	clientName, err := client.UserName(h.DB)
	if err != nil {
		fmt.Printf("An error occured during retrieving first name of client. %v \n", err)
		return
	}

	err = h.SendOutServerMessage(client.GroupID, fmt.Sprintf("%s has left the chat.", clientName))
	if err != nil {
		fmt.Printf("An error occured while sending out server message. %v \n", err)
		return
	}
}

func (h *Hub) HandleMessage(message Message) {
	if err := AddMessage(h.DB, message); err != nil {
		fmt.Printf("An error occured during adding to database. %v \n", err)
		return
	}

	messagePayload, err := message.Payload(h.DB)
	if err != nil {
		fmt.Printf("An error occured during loading. %v \n", err)
		return
	}
	err = h.SendOutToGroup(message.GroupID, messagePayload, func(c *Client)bool{
		return true
	})
	if err != nil {
		fmt.Printf("An error occured while getting recipient set. %v \n", err)
		return
	}
}