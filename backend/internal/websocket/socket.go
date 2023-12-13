package socket

import (
	"backend/internal/data"
	"backend/internal/helper"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var AllClients = make(map[int]*Client)

var AllGroupChats = make(map[string]map[int]bool)

type Client struct {
	Session        string
	Username       string
	Uuid           int
	Socket         *websocket.Conn
	Join           chan bool
	Leave          chan bool
	GroupMessage   chan data.UserMessage
	PrivateMessage chan data.UserMessage
	Notification   chan data.UserMessage
}

func NewClient(conn *websocket.Conn, username string, userId int, session string) *Client {

	return &Client{
		Session:        session,
		Username:       username,
		Uuid:           userId,
		Socket:         conn,
		Join:           make(chan bool),
		Leave:          make(chan bool),
		GroupMessage:   make(chan data.UserMessage),
		PrivateMessage: make(chan data.UserMessage),
		Notification:   make(chan data.UserMessage),
	}
}

func AddGroupChatUser(groupName string, userId int) error {
	for group := range AllGroupChats {
		if _, ok := AllGroupChats[group][userId]; ok {
			delete(AllGroupChats[group], userId)
		}
	}
	if _, ok := AllGroupChats[groupName]; !ok {
		AllGroupChats[groupName] = make(map[int]bool)
	}
	AllGroupChats[groupName][userId] = true
	fmt.Println("new group member in chat", AllGroupChats[groupName])
	return nil
}

func deleteUserFromAllGroups(userId int) {
	for group := range AllGroupChats {
		if _, ok := AllGroupChats[group][userId]; ok {
			delete(AllGroupChats[group], userId)
		}
	}
}

func updateUserlistGroupChat(groupName string) {

	msg := data.UserMessage{}
	msg.Sender = groupName
	msg.Type = "update_group_users"
	msg.Created = time.Now()

	for groupMemberId := range AllGroupChats[groupName] {

		username, err := helper.GetUsername(groupMemberId)
		if err != nil {
			delete(AllGroupChats[groupName], groupMemberId)
		}

		msg.Receiver = username

		if err := AllClients[groupMemberId].Socket.WriteJSON(msg); err != nil {
			fmt.Println("error sending groupchat users update", err)
			AllClients[groupMemberId].Socket.Close()
		}
	}
}

func updateUserlist() {
	msg := data.UserMessage{}
	msg.Sender = "server"
	msg.Receiver = "everyone"
	msg.Created = time.Now()
	msg.Type = "group_message"
	for _, client := range AllClients {
		if err := client.Socket.WriteJSON(msg); err != nil {
			fmt.Println("error at updateUserlist: ", err)
			client.Socket.Close()
		}
	}
}

func AddClient(c *Client) {
	for v := range AllClients {
		if v == c.Uuid {
			discon := AllClients[c.Uuid]
			delete(AllClients, c.Uuid)
			fmt.Println("removed")
			discon.Socket.Close()
		}
	}
	AllClients[c.Uuid] = c
	fmt.Printf("current userlist %v \n", AllClients)
	fmt.Println("current groupchat list", AllGroupChats)
}

func (c *Client) Read() {
	defer func() {
		c.Leave <- true
	}()
	for {
		//fmt.Println("users", AllClients)
		msg := data.UserMessage{}
		err := c.Socket.ReadJSON(&msg)
		if err != nil {
			fmt.Println("problem reading message", err, msg, c.Leave, c.Username)
			break
		}
		msg.Created = time.Now()
		msg.Sender = c.Username
		if msg.Type == "group_message" {
			err = helper.AddGroupMessageToDB(msg)
			if err != nil {
				return
			}
			c.GroupMessage <- msg
		}
		if msg.Type == "private_message" {
			err = helper.AddPrivateMessageToDB(msg)
			if err != nil {
				return
			}
			c.PrivateMessage <- msg
		}

	}
	c.Socket.Close()
}

func (c *Client) Write(msg data.UserMessage) {
	if msg.Type == "typing_update" || msg.Type == "update_group_users" {
		return
	}
	err := c.Socket.WriteJSON(msg)
	if err != nil {
		fmt.Println(err)
		c.Socket.Close()
	}
}

func FindAndSend(msg data.UserMessage) {
	recipientId, err := helper.GetuuidFromEmailOrUsername(msg.Receiver)
	if err != nil {
		return
	}
	if val, ok := AllClients[recipientId]; ok {
		//fmt.Printf("receiver found %v %v", val, ok)
		if err := val.Socket.WriteJSON(msg); err != nil {
			fmt.Println("problem sending message", err)
			val.Socket.Close()
		}
	}
}

func FindAndSendToGroup(msg data.UserMessage) {

	if _, ok := AllGroupChats[msg.Receiver]; ok {
		for memberId := range AllGroupChats[msg.Receiver] {
			if _, ok := AllClients[memberId]; ok {
				if err := AllClients[memberId].Socket.WriteJSON(msg); err != nil {
					AllClients[memberId].Socket.Close()
				}
			} else {
				delete(AllGroupChats[msg.Receiver], memberId)
			}
		}
	}
}

func (c *Client) UpdateClient() {
	for {
		select {
		case <-c.Join:
			AddClient(c)

		case <-c.Leave:
			//fmt.Println("user disconnected")
			deleteUserFromAllGroups(c.Uuid)
			delete(AllClients, c.Uuid)

			fmt.Println("current groupchat list after closing socket", AllGroupChats)
		case msg := <-c.GroupMessage:
			fmt.Println("new message", msg)
			FindAndSendToGroup(msg)
		case msg := <-c.PrivateMessage:
			FindAndSend(msg)
			c.Write(msg)
		case msg := <-c.Notification:
			FindAndSend(msg)
		}
	}
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("websocket reached")
	fmt.Println(r.URL)

	helper.EnableCors(&w)

	session, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println("websocket; cookie does not exist", err)
		return
	}
	sessionExists := helper.CheckSessionExist(session.Value)
	if !sessionExists {
		fmt.Println("websocket; session not found", err)
		return
	}
	userID, _ := helper.GetIdBySession(w, r)

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err, "error with upgrader")
		return
	}
	username, err := helper.GetUsernameBySession(session.Value)
	if err != nil {
		fmt.Println("socket name error", err)
		return
	}
	c := NewClient(ws, username, userID, session.Value)
	go c.UpdateClient()
	c.Join <- true
	go c.Read()

}