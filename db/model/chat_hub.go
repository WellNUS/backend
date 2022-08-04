package model

import (
	"database/sql"
	"fmt"
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







// General Helper functions
func (h *Hub) SendOutToAll(payload interface{}, cFilter func(*Client)bool) {
	for client := range h.Clients {
		if !cFilter(client) { continue }
		select {
			case client.Send <- payload:
			default:
				close(client.Send)
				delete(h.Clients, client)
		}
	}
}

func (h *Hub) SendOutChatStatus(userID int64) error {
	// userID is of user that induce the change in chat status
	// Send out group status
	err := h.SendOutGroupStatusPayload(userID)
	if err != nil { return err }
	// send out direct status
	err = h.SendOutUserStatusPayload(userID)
	if err != nil { return err }
	return nil
}






// Direct related functions
func (h *Hub) SendOutToSelectedClients(userIDs []int64, payload interface{}, cFilter func(*Client)bool) error {
	recipientsMap := make(map[int64]bool)
	for _, id := range userIDs {
		recipientsMap[id] = true
	}
	h.SendOutToAll(payload, func(c *Client)bool{
		return recipientsMap[c.UserID] && cFilter(c)
	})
	return nil
}

func (h *Hub) SendOutUserStatusPayload(userID int64) error {
	user, err := GetUser(h.DB, userID)
	if err != nil { return err }

	var targetClient *Client = nil
	allClients := make(map[int64]*Client, 0)
	observingClients := make([]*Client, 0)

	for client := range h.Clients {
		allClients[client.UserID] = client

		if client.TargetIsGroup { continue }
		if client.UserID == userID {
			targetClient = client
		} else if client.TargetID == userID {
			// To update others observing user
			observingClients = append(observingClients, client)
		}
	}

	// Sending observing clients
	for _, client := range observingClients {
		status := OfflineStatus
		if targetClient != nil {
			status = OnlineStatus
			if !targetClient.TargetIsGroup && targetClient.TargetID == client.UserID {
				status = InChatStatus
			}
		}
		client.Send <- UserStatusPayload{
			Tag: ChatStatusTag,
			Label: UserStatusLabel,
			User: user,
			Status: status,
		}
	}

	// Sending target client
	if targetClient != nil {
		observedUser, err := GetUser(h.DB, targetClient.TargetID)
		if err != nil { return err }

		status := OfflineStatus
		observedClient, ok := allClients[targetClient.TargetID]
		if ok {
			status = OnlineStatus
			if !observedClient.TargetIsGroup && observedClient.TargetID == userID {
				status = InChatStatus
			}
		}
		targetClient.Send <- UserStatusPayload{
			Tag: ChatStatusTag,
			Label: UserStatusLabel,
			User: observedUser,
			Status: status,
		}
	}
	return nil
}




// Group related functions
func (h *Hub) SendOutToGroup(groupID int64, payload interface{}, cFilter func(*Client)bool) error {
	recipients, err := GetAllUsersOfGroup(h.DB, groupID)
	if err != nil { return err }
	recipientsMap := make(map[int64]bool)
	for _, user := range recipients {
		recipientsMap[user.ID] = true
	}
	if err != nil { return err }
	h.SendOutToAll(payload, func(c *Client)bool{
		return recipientsMap[c.UserID] && cFilter(c)
	})
	return nil
}

func (h *Hub) SendOutServerMessageToGroupInChat(groupID int64, serverMsg string) error {
	err := h.SendOutToGroup(
		groupID, 
		ServerMessagePayload{
			Tag: MessageTag,
			Label: ServerMessageLabel,
			Msg: serverMsg,
		}, 
		func(c *Client)bool{
			return c.TargetIsGroup && c.TargetID == groupID
		});
	if err != nil { return err }
	return nil
}

func (h *Hub) GroupStatusPayload(group Group) (GroupStatusPayload, error) {
	usersInGroup, err := GetAllUsersOfGroup(h.DB, group.ID)
	if err != nil { return GroupStatusPayload{}, err }

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
	for client := range h.Clients {
		user, ok := usersInGroupMap[client.UserID]
		if ok {
			if client.TargetIsGroup && client.TargetID == group.ID {
				inChatMembers = append(inChatMembers, user)
			} else {
				onlineMembers = append(onlineMembers, user)
			}
			delete(usersInGroupMap, client.UserID)
		}
	}
	for _, user  := range usersInGroupMap {
		offlineMembers = append(offlineMembers, user)
	}
	sort.Slice(inChatMembers, MakeLess(inChatMembers))
	sort.Slice(onlineMembers, MakeLess(onlineMembers))
	sort.Slice(offlineMembers, MakeLess(offlineMembers))

	return GroupStatusPayload{
		Tag: ChatStatusTag, 
		GroupID: group.ID,
		GroupName: group.GroupName, 
		SortedInChatMembers: inChatMembers,
		SortedOnlineMembers: onlineMembers,
		SortedOfflineMembers: offlineMembers,
	}, nil
}

func (h *Hub) SendOutGroupStatusPayload(userID int64) error {
	groups, err := GetAllGroupsOfUser(h.DB, userID)
	if err != nil { return err }
	for _, group := range groups {
		groupStatusPayload, err := h.GroupStatusPayload(group)
		if err != nil { return err }
		err = h.SendOutToGroup(group.ID, groupStatusPayload, func(c *Client)bool{
			return c.TargetIsGroup && c.TargetID == group.ID
		})
		if err != nil { return err }
	}
	return err
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
				if sentData.Client.TargetIsGroup {
					groupMessage, err := sentData.ToGroupMessage()
					if err != nil {
						fmt.Printf("An error occured while converting sentData to GroupMessage. %v \n", err)
					}
					h.HandleGroupMessage(groupMessage)
				} else {
					directMessage, err := sentData.ToDirectMessage()
					if err != nil {
						fmt.Printf("An error occured while converting sentData to DirectMessage %v \n", err)
					}
					h.HandleDirectMessage(directMessage)
				}
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
}

func (h *Hub) HandleGroupMessage(groupMessage GroupMessage) {
	if err := AddGroupMessage(h.DB, groupMessage); err != nil {
		fmt.Printf("An error occured during adding to database. %v \n", err)
		return
	}

	groupMessagePayload, err := groupMessage.Payload(h.DB)
	if err != nil {
		fmt.Printf("An error occured during loading. %v \n", err)
		return
	}

	err = h.SendOutToGroup(groupMessage.GroupID, groupMessagePayload, func(c *Client)bool{
		return true
	})
	if err != nil {
		fmt.Printf("An error occured while sending message to group. %v \n", err)
		return
	}
}

func (h *Hub) HandleDirectMessage(directMessage DirectMessage) {
	if err := AddDirectMessage(h.DB, directMessage); err != nil {
		fmt.Printf("An error occured during adding to database. %v \n", err)
		return
	}

	directMessagePayload, err := directMessage.Payload(h.DB)
	if err != nil {
		fmt.Printf("An error occured during loading. %v \n", err)
		return
	}

	err = h.SendOutToSelectedClients(
		[]int64{directMessage.SenderID, directMessage.RecipientID}, 
		directMessagePayload, 
		func(c *Client)bool{
			return true
		},
	)
	if err != nil {
		fmt.Printf("An error occured while sending message to selected clients. %v \n", err)
		return
	}
}