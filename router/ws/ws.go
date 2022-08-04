package ws

import (
	"wellnus/backend/db/model"
	"wellnus/backend/config"
	"wellnus/backend/router/http_helper"
	"wellnus/backend/router/http_helper/http_error"
	
	"log"
	"fmt"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Client = model.Client
type Hub = model.Hub

const (
	loadedMessageBuffer = 256
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == config.FRONTEND_ADDRESS || origin == config.BACKEND_ADDRESS
	},
}

func ServeWs(Hub *Hub, w http.ResponseWriter, r *http.Request, userID int64, targetID int64, targetIsGroup bool) {
	Conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		UserID: userID, 
		TargetID: targetID,
		TargetIsGroup: targetIsGroup,
		Hub: Hub, 
		Conn: Conn, 
		Send: make(chan interface{}, loadedMessageBuffer),
	}
	client.Hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()
}

func ConnectToWSHandler(wsHub *Hub, db *sql.DB, targetIsGroup bool) func(*gin.Context) {
	return func(c *gin.Context) {
		http_helper.SetHeaders(c)

		targetID, err := http_helper.GetIDParams(c)
		if err != nil {
			fmt.Printf("An error occured when retrieving group ID params. %v \n", err)
			return
		}
		userID, err := http_helper.GetUserIDFromSessionCookie(db, c)
		if err != nil {
			fmt.Printf("An error occured when retrieving user ID cookies. %v \n", err)
			return
		}

		if targetIsGroup {
			isMember, err := model.IsUserInGroup(db, userID, targetID)
			if err != nil {
				fmt.Printf("An error occured when checking if user is in group. %v \n", err)
				return
			}
			if !isMember {
				err = http_error.UnauthorizedError
				fmt.Printf("User is not part of group. %v \n", err)
				return
			}
		}
		
		ServeWs(wsHub, c.Writer, c.Request, userID, targetID, targetIsGroup)
	}
}