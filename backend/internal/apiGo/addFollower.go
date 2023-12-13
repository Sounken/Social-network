package apiGO

import (
	"backend/internal/data"
	"backend/internal/helper"
	socket "backend/internal/websocket"
	"encoding/json"
	"fmt"
	"net/http"
)

func removeYourFollowRequest(yourId, theirId int) error {
	sqlString := `DELETE FROM notifications
	WHERE (sender_id = ? AND uuid = ? AND notif_type = 'follow_request')`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return err
	}

	defer sqlStmt.Close()

	_, err = sqlStmt.Exec(yourId, theirId)

	return err
}

func isPrivate(uuid int) bool {
	sqlString := `SELECT privacy
	FROM users
	WHERE uuid = ?`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return false
	}

	defer sqlStmt.Close()

	var res string
	sqlStmt.QueryRow(uuid).Scan(&res)
	return res == "private"
}

func addFollower(followerId int, followedId int) error {

	sqlStmt, err := data.DB.Prepare(`INSERT INTO followers (uuid,follower_id)
	VALUES(?,?)`)
	if err != nil {
		fmt.Println("AUUUUUGHHH", err)
		return err
	}
	defer sqlStmt.Close()

	_, err = sqlStmt.Exec(followedId, followerId)

	return err
}

func removeFollower(followerId int, followedId int) error {
	sqlStmt, err := data.DB.Prepare(`DELETE FROM followers
	WHERE followers.follower_id = ? AND followers.uuid = ?;`)
	if err != nil {
		fmt.Println("gosh darn it", err)
		return err
	}
	defer sqlStmt.Close()

	_, err = sqlStmt.Exec(followerId, followedId)
	return err
}

func AddOrRemoveFollow(w http.ResponseWriter, r *http.Request) {
	helper.EnableCors(&w)

	if r.Method == http.MethodPost {
		var email string
		err := json.NewDecoder(r.Body).Decode(&email)
		if err != nil {
			fmt.Println("add/remove follow decoding error", err)
			fmt.Println(err)
		}
		yourID, err := helper.GetIdBySession(w, r)
		if err != nil {
			helper.WriteResponse(w, "session_error")
			return
		}

		theirId, err := helper.GetuuidFromEmailOrUsername(email)
		if err != nil || theirId == yourID {
			fmt.Println("add follow uuid shit", err, email)
			helper.WriteResponse(w, "incorrect_email")
			return
		}

		alreadyFollowing := checkDBIfFollowing(yourID, theirId)
		fmt.Println("to follow / unfollow email", email)
		fmt.Println("your id", yourID)
		//fmt.Println("your email", yourEmail)
		if alreadyFollowing {
			err = removeFollower(yourID, theirId)
			fmt.Println("already following")
		} else {
			if isPrivate(theirId) {
				fmt.Println("private lol")

				addedNotification, err := helper.AddNotificationToDB("You have a new Follow request!", "follow_request", "", theirId, yourID)
				if err != nil {
					helper.WriteResponse(w, "database_error")
					fmt.Println("addnotif err", err)
					return
				}
				if addedNotification {
					socket.SendNotificationToAUser(yourID, theirId, "You have a new Follow request", "", "follow_request")
				}
			} else {
				err = removeYourFollowRequest(yourID, theirId)
				if err != nil {
					helper.WriteResponse(w, "database_error")
					fmt.Println("deletenotif err", err)
					return
				}
				err = addFollower(yourID, theirId)
			}
		}
		if err != nil {
			helper.WriteResponse(w, "database_error")
			return
		}
		helper.WriteResponse(w, "success")
	}
}
