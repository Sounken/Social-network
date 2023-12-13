package apiGO

import (
	"backend/internal/data"
	"backend/internal/helper"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
)

func CreateSessionToken(w http.ResponseWriter) string {
	sessionToken := uuid.Must(uuid.NewV4()).String()
	return sessionToken
}

func updateSessionToken(token string, uid int) error {
	sqlStmt, err := data.DB.Prepare("UPDATE users SET session_token = ? WHERE uuid = ?;")
	if err != nil {
		return err
	}
	defer sqlStmt.Close()

	_, err = sqlStmt.Exec(token, uid)
	if err != nil {
		return err
	}
	return nil
}

type loginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func checkLoginDetails(email, password string) (bool, int) {
	sqlString := "SELECT uuid FROM users WHERE email = ? AND password = ?;"
	var dum int

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return false, dum
	}

	defer sqlStmt.Close()

	err = sqlStmt.QueryRow(email, password).Scan(&dum)

	return err == nil, dum
}

func Login(w http.ResponseWriter, r *http.Request) {
	helper.EnableCors(&w)

	if r.Method == http.MethodPost {
		var logDat loginData
		err := json.NewDecoder(r.Body).Decode(&logDat)
		if err != nil {
			fmt.Println("failed decoding", err)
			return
		}
		if !helper.CheckIfStringExist("users", "email", logDat.Email) {
			helper.WriteResponse(w, "user_not_exist")
			fmt.Println("failed usercheck", err)
			return
		}
		credentialsMatch, uid := checkLoginDetails(logDat.Email, logDat.Password)
		if !credentialsMatch {
			helper.WriteResponse(w, "incorrect_password")
			fmt.Println("failed password check", err, credentialsMatch, uid)
			return
		}
		token := CreateSessionToken(w)
		updateSessionToken(token, uid)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"success", "token":"` + token + `"}`))

	}
}
