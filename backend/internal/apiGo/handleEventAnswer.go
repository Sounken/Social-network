package apiGO

import (
	"backend/internal/data"
	"backend/internal/helper"
	"encoding/json"
	"fmt"
	"net/http"
)

type eventResponse struct {
	GroupName string `json:"group_name"`
	Answer    string `json:"answer"`
	EventId   int    `json:"event_id"`
	Sender    string `json:"sender"`
}

func addEventResponseToDB(eventId, uuid int, groupName, choice string) (bool, error) {
	sqlString := `INSERT INTO eventOptionChoices(event_id,uuid,choice)
	SELECT ? AS event_id,
	u.uuid AS uuid,
	? AS choice
	FROM users AS u
	WHERE (u.uuid = ?)
	AND NOT EXISTS(
		SELECT 1 FROM eventOptionChoices AS e
		WHERE e.event_id = ?
		AND e.uuid = ?
	)`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return false, err
	}

	res, err := sqlStmt.Exec(eventId, choice, uuid, eventId, uuid)
	if err != nil {
		return false, err
	}

	rowsAffected, err := res.RowsAffected()

	return rowsAffected == 1, err
}

func HandleEventAnswer(w http.ResponseWriter, r *http.Request) {
	helper.EnableCors(&w)

	if r.Method == http.MethodPost {
		var response eventResponse

		err := json.NewDecoder(r.Body).Decode(&response)
		if err != nil {
			helper.WriteResponse(w, "decoding_error")
			fmt.Println("decoding event answer", err)
			return
		}
		fmt.Println("resp", response)

		groupExist := helper.CheckIfStringExist("groups", "group_name", response.GroupName)
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
		isMember, _ := helper.CheckIfGroupMember(response.GroupName, uuid)
		if !isMember {
			helper.WriteResponse(w, "not_a_member")
			return
		}

		addedResponse, err := addEventResponseToDB(response.EventId, uuid, response.GroupName, response.Answer)
		if err != nil {
			helper.WriteResponse(w, "database_error")
			return
		}
		if !addedResponse {
			helper.WriteResponse(w, "already_answered")
			return
		}

		helper.WriteResponse(w, "success")
	}
}
