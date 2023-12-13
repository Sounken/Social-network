package apiGO

import (
	"backend/internal/data"
	"backend/internal/helper"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type privateMessage struct {
	Sender   string    `json:"sender"`
	Receiver string    `json:"receiver"`
	Created  time.Time `json:"created"`
	Image    string    `json:"image"`
	Content  string    `json:"content"`
}

type privateMessages struct {
	Messages []privateMessage `json:"messages"`
	Status   string           `json:"status"`
}

func gatherPrivateMessagesfromDB(yourId int, otherUser int) ([]privateMessage, error) {
	var AllMessagesData []privateMessage
	sqlString := `SELECT COALESCE(u1.username, u1.email),
	COALESCE(u2.username, u2.email),
	pmg_content,
	IFNULL(pmg_image, 'http://localhost:8080/images/default.jpeg'),
	creation_date
	FROM privateMessages
	JOIN users AS 'u1'
	ON u1.uuid = pmg_sender
	JOIN users AS 'u2'
	ON u2.uuid = pmg_receiver
	WHERE (pmg_sender = ? AND pmg_receiver = ?)
	OR (pmg_sender = ? AND pmg_receiver = ?) 
	`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return AllMessagesData, err
	}

	defer sqlStmt.Close()

	rows, err := sqlStmt.Query(yourId, otherUser, otherUser, yourId)
	if err != nil {
		return AllMessagesData, err
	}

	defer rows.Close()

	for rows.Next() {
		var messageData privateMessage

		err = rows.Scan(&messageData.Sender, &messageData.Receiver, &messageData.Content, &messageData.Image, &messageData.Created)
		if err != nil {
			return AllMessagesData, err
		}
		AllMessagesData = append(AllMessagesData, messageData)
	}

	return AllMessagesData, err
}

func GetPrivateMessages(w http.ResponseWriter, r *http.Request) {
	helper.EnableCors(&w)

	if r.Method != http.MethodPost {
		return
	}

	var otherUser string
	err := json.NewDecoder(r.Body).Decode(&otherUser)
	if err != nil {
		fmt.Println(err)
		helper.WriteResponse(w, "decoding_error")
		return
	}
	uuid, err := helper.GetIdBySession(w, r)
	if err != nil {
		helper.WriteResponse(w, "session_error")
		return
	}
	var privateMessagesData privateMessages
	otherUseruuid, err := helper.GetuuidFromEmailOrUsername(otherUser)
	if err != nil {
		helper.WriteResponse(w, "user_does_not_exist")
		return
	}
	privateMessagesData.Messages, err = gatherPrivateMessagesfromDB(uuid, otherUseruuid)
	if err != nil {
		helper.WriteResponse(w, "database_error")
		return
	}

	privateMessagesData.Status = "success"
	privateMessagesDataJson, err := json.Marshal(privateMessagesData)
	if err != nil {
		fmt.Println("marshalling error", err)
		helper.WriteResponse(w, "marshalling_error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(privateMessagesDataJson)

}
