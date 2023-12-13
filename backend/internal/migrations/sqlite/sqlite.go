package sqlite

import "backend/internal/data"

func CheckSessionToken(sessionString string) (int, error) {
	stmt := "CHECK uid FROM users WHERE session_token = ?;"
	var uid int
	err := data.DB.QueryRow(stmt, sessionString).Scan(&uid)
	return uid, err
}

