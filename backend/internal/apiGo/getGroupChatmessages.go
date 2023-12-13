package apiGO

import (
	"backend/internal/data"
	"backend/internal/helper"
	"encoding/json"
	"fmt"
	"net/http"
)

func gatherGroupChatMessagesFromDB(groupName string) ([]privateMessage, error) {
	var allGroupChatMessages []privateMessage

	sqlString := `SELECT COALESCE(users.username, users.email),
	gc_content,
	IFNULL(gc_image, 'no_image'),
	creation_date
	FROM groupChat
	JOIN users
	ON gc_sender = users.uuid
	WHERE group_id IN (
		SELECT group_id
		FROM groups
		WHERE group_name = ?
	)`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return allGroupChatMessages, err
	}

	defer sqlStmt.Close()

	rows, err := sqlStmt.Query(groupName)
	if err != nil {
		return allGroupChatMessages, err
	}

	defer rows.Close()

	for rows.Next() {
		var groupChatMessage privateMessage

		err = rows.Scan(&groupChatMessage.Sender, &groupChatMessage.Content, &groupChatMessage.Image, &groupChatMessage.Created)
		if err != nil {
			return allGroupChatMessages, err
		}

		allGroupChatMessages = append(allGroupChatMessages, groupChatMessage)
	}

	return allGroupChatMessages, err
}

func GetGroupChatMessages(w http.ResponseWriter, r *http.Request) {

	helper.EnableCors(&w)

	if r.Method == http.MethodPost {

		var groupName string
		err := json.NewDecoder(r.Body).Decode(&groupName)
		if err != nil {
			helper.WriteResponse(w, "decoding_error")
			fmt.Println("decoding", err)
			return
		}
		groupExist := helper.CheckIfStringExist("groups", "group_name", groupName)
		if !groupExist {
			helper.WriteResponse(w, "group_does_not_exist")
			fmt.Println("group not exist", err)
			return
		}
		uuid, err := helper.GetIdBySession(w, r)
		if err != nil {
			helper.WriteResponse(w, "session_error")
			fmt.Println(err)
			return
		}
		isMember, _ := helper.CheckIfGroupMember(groupName, uuid)
		if !isMember {
			helper.WriteResponse(w, "not_a_member")
			fmt.Println("not a member groupchat", groupName, uuid)
			return
		}

		var groupMessagesData privateMessages

		groupMessagesData.Messages, err = gatherGroupChatMessagesFromDB(groupName)
		if err != nil {
			helper.WriteResponse(w, "database_error")
			return
		}

		groupMessagesData.Status = "success"
		groupMessagesDataJson, err := json.Marshal(groupMessagesData)
		if err != nil {
			fmt.Println("marshalling error", err)
			helper.WriteResponse(w, "marshalling_error")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(groupMessagesDataJson)
	}
}
