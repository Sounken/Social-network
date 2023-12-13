package apiGO

import (
	"backend/internal/data"
	"backend/internal/helper"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type notification struct {
	Id       int       `json:"id"`
	Content  string    `json:"content"`
	Sender   string    `json:"sender"`
	Receiver string    `json:"receiver"`
	Type     string    `json:"type"`
	Context  string    `json:"context"`
	Created  time.Time `json:"created"`
}

type allNotifications struct {
	Notifications []notification `json:"notifications"`
	Status        string         `json:"status"`
}

func getAllNotificationsFromDB(uuid int) ([]notification, error) {
	var allNotifications []notification
	sqlString := `SELECT notif_id,
	notif_content,
	creation_date,
	COALESCE(u2.username, u2.email),
	COALESCE(u1.username, u1.email),
	notif_type,
	notif_context
	FROM notifications
	JOIN users AS u1
	ON u1.uuid = notifications.sender_id
	JOIN users AS u2
	ON u2.uuid = ?
	WHERE notifications.uuid = ?`
	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return allNotifications, err
	}

	defer sqlStmt.Close()

	rows, err := sqlStmt.Query(uuid, uuid)
	if err != nil {
		return allNotifications, err
	}

	defer rows.Close()

	for rows.Next() {
		var userNotification notification

		err = rows.Scan(&userNotification.Id, &userNotification.Content, &userNotification.Created, &userNotification.Receiver, &userNotification.Sender, &userNotification.Type, &userNotification.Context)
		if err != nil {
			return allNotifications, err
		}

		allNotifications = append(allNotifications, userNotification)
	}

	return allNotifications, err
}

func GetAllNotifications(w http.ResponseWriter, r *http.Request) {
	helper.EnableCors(&w)

	if r.Method != http.MethodPost {
		return
	}

	uuid, err := helper.GetIdBySession(w, r)
	if err != nil {
		helper.WriteResponse(w, "session_error")
		return
	}
	var allUserNotifications allNotifications
	allUserNotifications.Notifications, err = getAllNotificationsFromDB(uuid)
	if err != nil {
		fmt.Println("get notifications error", err)
		helper.WriteResponse(w, "database_error")
		return
	}

	allUserNotifications.Status = "success"

	allUserNotificationsJson, err := json.Marshal(allUserNotifications)
	if err != nil {
		fmt.Println("marshalling error", err)
		helper.WriteResponse(w, "marshalling_error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(allUserNotificationsJson)
}
