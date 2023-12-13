package apiGO

import (
	"backend/internal/helper"
	socket "backend/internal/websocket"
	"encoding/json"
	"fmt"
	"net/http"
)

type onlineUsers struct {
	Users  []string `json:"users"`
	Status string   `jsonL:"status"`
}

func GetOnlineGroupMembers(w http.ResponseWriter, r *http.Request) {
	helper.EnableCors(&w)

	if r.Method != http.MethodPost {
		return
	}
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
		fmt.Println(groupName, uuid)
		return
	}

	var allOnlineGroupUsers onlineUsers

	for userID := range socket.AllGroupChats[groupName] {
		username, err := helper.GetUsername(userID)
		if err != nil {
			delete(socket.AllGroupChats[groupName], userID)
		}

		allOnlineGroupUsers.Users = append(allOnlineGroupUsers.Users, username)
	}

	allOnlineGroupUsers.Status = "success"

	allOnlineGroupUsersJson, err := json.Marshal(allOnlineGroupUsers)
	if err != nil {
		fmt.Println("marshalling error", err)
		helper.WriteResponse(w, "marshalling_error")
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(allOnlineGroupUsersJson)
}
