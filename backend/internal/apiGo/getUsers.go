package apiGO

import (
	"backend/internal/data"
	"backend/internal/helper"
	"encoding/json"
	"fmt"
	"net/http"
)

type user struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type users struct {
	User   []user `json:"users"`
	Status string `json:"status"`
}

func getUsernameAndEmail() (users, error) {
	var userDat users

	sqlString := `SELECT email,
	IFNULL(username, email) FROM users`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return userDat, err
	}

	defer sqlStmt.Close()

	rows, err := sqlStmt.Query()
	if err != nil {
		return userDat, err
	}

	defer rows.Close()

	for rows.Next() {
		var usr user
		err = rows.Scan(&usr.Email, &usr.Username)
		if err != nil {
			return userDat, err
		}
		userDat.User = append(userDat.User, usr)
	}
	return userDat, nil
}

func GetUsernames(w http.ResponseWriter, r *http.Request) {
	helper.EnableCors(&w)
	if r.Method == http.MethodPost {
		userData, err := getUsernameAndEmail()
		if err != nil {
			helper.WriteResponse(w, "database_error")
			fmt.Println("database error get usernames", err)
			return
		}
		userData.Status = "success"
		userDataJson, err := json.Marshal(userData)
		if err != nil {
			fmt.Println("marshalling error", err)
			helper.WriteResponse(w, "marshalling_error")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(userDataJson)
	}
}
