package apiGO

import (
	"backend/internal/data"
	"backend/internal/helper"
	"encoding/json"
	"fmt"
	"net/http"
)

type isFollowing struct {
	Following bool   `json:"following"`
	Status    string `json:"status"`
}

func checkDBIfFollowing(yourId, theirId int) bool {
	sqlString := `SELECT followers.uuid
	FROM followers
	WHERE followers.follower_id = ? AND followers.uuid = ?;`
	var lole int

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return false
	}

	defer sqlStmt.Close()

	err = sqlStmt.QueryRow(yourId, theirId).Scan(&lole)

	return err == nil
}

func CheckFollowing(w http.ResponseWriter, r *http.Request) {
	helper.EnableCors(&w)
	if r.Method == http.MethodPost {
		var followInfo isFollowing
		uuid, err := helper.GetIdBySession(w, r)
		if err != nil {
			helper.WriteResponse(w, "incorrect_session")
			return
		}
		var email string
		err = json.NewDecoder(r.Body).Decode(&email)
		if err != nil {
			fmt.Println(err)
			helper.WriteResponse(w, "decoding_error")
		}
		theirId, err := helper.GetuuidFromEmailOrUsername(email)
		if err != nil {
			helper.WriteResponse(w, "user_does_not_exist")
			return
		}
		followInfo.Following = checkDBIfFollowing(uuid, theirId)
		followInfo.Status = "success"

		followInfoJson, err := json.Marshal(followInfo)
		if err != nil {
			fmt.Println("marshalling error", err)
			helper.WriteResponse(w, "marshalling_error")
		}
		fmt.Println("checkFollow success")
		w.Header().Set("Content-Type", "application/json")
		w.Write(followInfoJson)
	}
}
