package socket

import (
	"backend/internal/helper"
	"encoding/json"
	"fmt"
	"net/http"
)

func AddGroupMemberToChat(w http.ResponseWriter, r *http.Request) {

	helper.EnableCors(&w)

	if r.Method == http.MethodPost {

		var groupName string
		err := json.NewDecoder(r.Body).Decode(&groupName)
		//fmt.Println("post data", pData)
		if err != nil {
			helper.WriteResponse(w, "decoding_error")
			fmt.Println(err)
			return
		}

		uuid, err := helper.GetIdBySession(w, r)

		if err != nil {
			helper.WriteResponse(w, "session_error")
			return
		}

		groupExist := helper.CheckIfStringExist("groups", "group_name", groupName)
		if !groupExist {
			helper.WriteResponse(w, "group_does_not_exist")
			return
		}

		isMember, _ := helper.CheckIfGroupMember(groupName, uuid)
		if !isMember {
			helper.WriteResponse(w, "not_a_member")
			return
		}

		err = AddGroupChatUser(groupName, uuid)
		if err != nil {
			helper.WriteResponse(w, "error_adding_to_chat")
			return
		}

		helper.WriteResponse(w, "success")
	}
}
