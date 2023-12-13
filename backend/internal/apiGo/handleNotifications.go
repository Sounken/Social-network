package apiGO

import (
	"backend/internal/data"
	"backend/internal/helper"
	"encoding/json"
	"fmt"
	"net/http"
)

type notificationResponse struct {
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	GroupName string `json:"group_name"`
	Type      string `json:"type"`
	Response  string `json:"response"`
}

func addMemberToGroup(member, groupName string) error {
	sqlString := `INSERT INTO groupMembers(group_id,uuid,role)
	SELECT group_id,
	users.uuid,
	'member' AS role
	FROM groups
	JOIN users
	ON users.username = ? OR users.email = ?
	WHERE group_name = ?`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return err
	}

	defer sqlStmt.Close()

	_, err = sqlStmt.Exec(member, member, groupName)

	return err
}

func checkIfAlreadyGroupMember(username, groupName string) bool {
	sqlString := `SELECT
	group_id
	FROM groupMembers
	WHERE uuid IN (
		SELECT u.uuid
		FROM users AS u
		WHERE u.email = ? OR u.username = ?
	)
	AND group_id IN (
		SELECT g.group_id
		FROM groups AS g
		WHERE g.group_name = ?
	)`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return false
	}
	var dummy int

	defer sqlStmt.Close()

	err = sqlStmt.QueryRow(username, username, groupName).Scan(&dummy)
	return err == nil
}

func addFollowerToDB(uuid int, follower string) error {
	sqlString := `INSERT INTO followers(uuid,follower_id)
	SELECT ? AS uuid,
	u.uuid AS follower_id
	FROM users AS u
	WHERE u.username = ?
	OR u.email = ?`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return err
	}

	defer sqlStmt.Close()

	_, err = sqlStmt.Exec(uuid, follower, follower)

	return err
}

func deleteGroupNotificationsFromDB(sender, receiver string) error {
	sqlString := `DELETE FROM notifications
	WHERE (sender_id IN (
		SELECT u.uuid
		FROM users AS u
		WHERE u.username = ? OR u.email = ?
	) AND notif_type = 'group_join_request')
	OR (uuid IN (
		SELECT u2.uuid
		FROM users AS u2
		WHERE u2.username = ? OR u2.email = ?
	) AND notif_type = 'group_invite')`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return err
	}

	defer sqlStmt.Close()

	_, err = sqlStmt.Exec(sender, sender, receiver, receiver)

	return err
}

func deleteFollowRequestsFromDB(sender string, receiver int) error {
	sqlString := `DELETE FROM notifications
	WHERE (uuid = ? AND sender_id IN (
		SELECT u.uuid
		FROM users AS u
		WHERE u.email = ? OR u.username = ?
	) AND notif_type = 'follow_request')`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return err
	}

	defer sqlStmt.Close()

	_, err = sqlStmt.Exec(receiver, sender, sender)
	return err
}

func HandleGroupInvite(w http.ResponseWriter, r *http.Request) {

	helper.EnableCors(&w)

	if r.Method != http.MethodPost {
		return
	}

	var notifInfo notificationResponse
	err := json.NewDecoder(r.Body).Decode(&notifInfo)
	if err != nil {
		helper.WriteResponse(w, "decoding_error")
		fmt.Println("decoding", err)
		return
	}

	groupExist := helper.CheckIfStringExist("groups", "group_name", notifInfo.GroupName)
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

	receiveruuid, err := helper.GetuuidFromEmailOrUsername(notifInfo.Receiver)
	if err != nil || uuid != receiveruuid {
		helper.WriteResponse(w, "fucked")
		fmt.Println("user not matching with notification", err)
		fmt.Println("info", notifInfo.Receiver, notifInfo.Sender)
		return
	}

	isAlreadyMember := checkIfAlreadyGroupMember(notifInfo.Receiver, notifInfo.GroupName)
	if isAlreadyMember {
		fmt.Println("already in the group")
		deleteGroupNotificationsFromDB(notifInfo.Sender, notifInfo.Receiver)
		helper.WriteResponse(w, "success")
		return
	}

	if notifInfo.Response == "accept" {
		err = addMemberToGroup(notifInfo.Receiver, notifInfo.GroupName)
		if err != nil {
			helper.WriteResponse(w, "database_error")
			fmt.Println("failed joining group", err)
			return
		}
		err = deleteGroupNotificationsFromDB(notifInfo.Sender, notifInfo.Receiver)
		if err != nil {
			helper.WriteResponse(w, "database_error")
			return
		}
	}
	if notifInfo.Response == "decline" {
		err = deleteGroupNotificationsFromDB(notifInfo.Sender, notifInfo.Receiver)
		if err != nil {
			helper.WriteResponse(w, "database_error")
			return
		}
	}

	helper.WriteResponse(w, "success")
}

func HandleFollowRequest(w http.ResponseWriter, r *http.Request) {

	helper.EnableCors(&w)

	if r.Method != http.MethodPost {
		return
	}

	var notifInfo notificationResponse
	err := json.NewDecoder(r.Body).Decode(&notifInfo)
	if err != nil {
		helper.WriteResponse(w, "decoding_error")
		fmt.Println("decoding", err)
		return
	}

	uuid, err := helper.GetIdBySession(w, r)
	if err != nil {
		helper.WriteResponse(w, "session_error")
		fmt.Println(err)
		return
	}

	receiveruuid, err := helper.GetuuidFromEmailOrUsername(notifInfo.Receiver)
	if err != nil || receiveruuid != uuid {
		fmt.Println("liar detected", notifInfo.Receiver, notifInfo.Sender, receiveruuid, uuid)
		helper.WriteResponse(w, "stop_lying_pls")
		return
	}

	if notifInfo.Response != "accept" && notifInfo.Response != "decline" {
		helper.WriteResponse(w, "incorrect_response")
		return
	}

	if notifInfo.Response == "accept" {
		err = addFollowerToDB(uuid, notifInfo.Sender)
		if err != nil {
			fmt.Println("database error add follower", err)
			helper.WriteResponse(w, "database_error")
			return
		}
	}
	err = deleteFollowRequestsFromDB(notifInfo.Sender, uuid)
	if err != nil {
		helper.WriteResponse(w, "error_deleting_notification")
		return
	}

	helper.WriteResponse(w, "success")

}

func HandleGroupJoinRequest(w http.ResponseWriter, r *http.Request) {

	helper.EnableCors(&w)

	if r.Method != http.MethodPost {
		return
	}

	var notifInfo notificationResponse
	err := json.NewDecoder(r.Body).Decode(&notifInfo)
	if err != nil {
		helper.WriteResponse(w, "decoding_error")
		fmt.Println("decoding", err)
		return
	}

	groupExist := helper.CheckIfStringExist("groups", "group_name", notifInfo.GroupName)
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

	isMember, memberType := helper.CheckIfGroupMember(notifInfo.GroupName, uuid)
	if !isMember || memberType != "creator" {
		helper.WriteResponse(w, "not_group_creator")
		return
	}

	isAlreadyMember := checkIfAlreadyGroupMember(notifInfo.Sender, notifInfo.GroupName)
	if isAlreadyMember {
		fmt.Println("already in the group")
		deleteGroupNotificationsFromDB(notifInfo.Sender, notifInfo.Receiver)
		helper.WriteResponse(w, "success")
		return
	}

	if notifInfo.Response != "accept" && notifInfo.Response != "decline" {
		helper.WriteResponse(w, "incorrect_response")
		return
	}

	if notifInfo.Response == "accept" {

		_, err := helper.GetuuidFromEmailOrUsername(notifInfo.Sender)
		if err != nil {
			helper.WriteResponse(w, "user_does_not_exist")
			return
		}

		err = addMemberToGroup(notifInfo.Sender, notifInfo.GroupName)
		if err != nil {
			fmt.Println("add member err", err)
			fmt.Println("add member err", err)
			helper.WriteResponse(w, "database_error")
			return
		}
	}

	err = deleteGroupNotificationsFromDB(notifInfo.Sender, notifInfo.Receiver)
	if err != nil {
		fmt.Println("delete notif error", err)
		helper.WriteResponse(w, "database_error")
		return
	}

	helper.WriteResponse(w, "success")

}

func DeleteNotification(w http.ResponseWriter, r *http.Request) {
	helper.EnableCors(&w)

	if r.Method == http.MethodPost {
		var notifInfo notificationResponse
		err := json.NewDecoder(r.Body).Decode(&notifInfo)
		if err != nil {
			helper.WriteResponse(w, "decoding_error")
			fmt.Println("decoding error delete notification", err)
			return
		}

		uuid, err := helper.GetIdBySession(w, r)
		if err != nil {
			helper.WriteResponse(w, "session_error")
		}

		receiveruuid, err := helper.GetuuidFromEmailOrUsername(notifInfo.Receiver)
		if err != nil || receiveruuid != uuid {
			helper.WriteResponse(w, "stop_lying_pls")
			return
		}

		senderId, err := helper.GetuuidFromEmailOrUsername(notifInfo.Sender)
		if err != nil {
			helper.WriteResponse(w, "stop_lying_pls")
			return
		}
		err = helper.DeleteNotificationFromDB(senderId, receiveruuid, "")
		if err != nil {
			helper.WriteResponse(w, "database_error")
			return
		}

		helper.WriteResponse(w, "success")

	}
}
