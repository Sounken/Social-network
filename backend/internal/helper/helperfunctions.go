package helper

import (
	"backend/internal/data"
	"fmt"
	"net/http"
	"time"
)

func WriteResponse(w http.ResponseWriter, status string) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"` + status + `"}`))
}

func EnableCors(w *http.ResponseWriter) {
	//(*w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, withCredentials, credentials")
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}

func UpdateTableColumnStringById(table, newData, column string, uid int) error {
	sqlStmt, err := data.DB.Prepare("UPDATE " + table + " SET " + column + " = ? WHERE uuid = ?;")
	if err != nil {
		return err
	}
	defer sqlStmt.Close()

	_, err = sqlStmt.Exec(newData, uid)
	if err != nil {
		return err
	}
	return nil
}

func UpdateTableColumnByteById(table string, newData []byte, column string, uid int) error {
	sqlStmt, err := data.DB.Prepare("UPDATE " + table + " SET " + column + " = ? WHERE uuid = ?;")
	if err != nil {
		return err
	}
	defer sqlStmt.Close()

	_, err = sqlStmt.Exec(newData, uid)
	if err != nil {
		return err
	}
	return nil
}

func GetIdBySession(w http.ResponseWriter, r *http.Request) (int, error) {
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println("haaaa?", err)
		return -1, err
	}

	var uid int

	sqlString := "SELECT uuid FROM users WHERE session_token = ?;"

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return uid, err
	}

	defer sqlStmt.Close()

	err = sqlStmt.QueryRow(sessionToken.Value).Scan(&uid)
	if err != nil {
		return -1, err
	}

	return uid, nil
}

func CheckIfStringExist(table, column, tableData string) bool {
	sqlString := "SELECT " + column + " FROM " + table + " WHERE " + column + " = ?;"
	var dummy string

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return false
	}

	defer sqlStmt.Close()

	err = sqlStmt.QueryRow(tableData).Scan(&dummy)

	return err == nil
}

func GetUsername(uuid int) (string, error) {
	sqlString := `SELECT COALESCE(username,email)
	FROM users
	WHERE uuid = ?`

	var username string

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return "", err
	}

	defer sqlStmt.Close()

	err = sqlStmt.QueryRow(uuid).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, err
}

func GetFollowing(uuid int) ([]string, error) {
	sqlString := `SELECT followers.uuid,
COALESCE(users.username, users.email)
FROM followers
JOIN users
ON users.uuid = followers.uuid
WHERE follower_id = ?;`
	var bleh = make([]string, 0)

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return bleh, err
	}

	defer sqlStmt.Close()

	rows, err := sqlStmt.Query(uuid)
	if err != nil {
		fmt.Println("query error", err)
		return bleh, err
	}

	defer rows.Close()

	for rows.Next() {
		var d1 string
		var d2 string
		err = rows.Scan(&d1, &d2)
		if err != nil {
			fmt.Println("scan error", err)
			return bleh, err
		}
		//fmt.Println("var1", d1, "var2", d2)
		bleh = append(bleh, d2)
	}

	return bleh, err
}

func GetFollowers(uuid int) ([]string, error) {
	sqlString := `SELECT COALESCE(users.username, users.email)
	FROM followers
	JOIN users
	ON users.uuid = followers.follower_id
	WHERE followers.uuid = ?;`
	var result []string

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return result, err
	}

	defer sqlStmt.Close()

	rows, err := sqlStmt.Query(uuid)
	if err != nil {
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		var usr string

		err = rows.Scan(&usr)
		if err != nil {
			return result, err
		}
		result = append(result, usr)
	}

	return result, err
}

func GetYourGroups(uuid int) ([]string, error) {
	sqlString := `SELECT groupMembers.group_id,
	groups.group_name
	FROM groupMembers
	JOIN groups
	ON groups.group_id = groupMembers.group_id
	WHERE uuid = ?`
	var res = make([]string, 0)

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return res, err
	}

	defer sqlStmt.Close()

	rows, err := sqlStmt.Query(uuid)
	if err != nil {
		fmt.Println("query getYourGroups error", err)
	}

	defer rows.Close()

	for rows.Next() {
		var temp string
		var useless int
		err = rows.Scan(&useless, &temp)
		if err != nil {
			fmt.Println("scan getYourGroups error", err)
			return res, err
		}
		res = append(res, temp)
	}
	return res, err
}

func AddNotificationToDB(content, nType, context string, usr, sender int) (bool, error) {
	sqlString := `INSERT INTO notifications (notif_content,creation_date,uuid,sender_id,notif_type,notif_context)
	SELECT ? AS notif_content,
	? AS creation_date,
	? AS uuid,
	? AS sender_id,
	? AS notif_type,
	? AS notif_context
	WHERE NOT EXISTS(
		SELECT 1 FROM notifications
		WHERE uuid = ? AND sender_id = ? AND notif_type = ?
	);`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return true, err
	}

	defer sqlStmt.Close()

	res, err := sqlStmt.Exec(content, time.Now(), usr, sender, nType, context, usr, sender, nType)
	if err != nil {
		return true, err
	}

	affectedRows, err := res.RowsAffected()

	return affectedRows == 1, err
}

func GetuuidByString(tabType, value string) (int, error) {
	sqlString := `SELECT uuid
	FROM users
	WHERE ` + tabType + ` = ?`

	var result int

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return result, err
	}

	defer sqlStmt.Close()

	err = sqlStmt.QueryRow(value).Scan(&result)

	return result, err
}

func GetuuidFromEmailOrUsername(user string) (int, error) {
	sqlString := `SELECT uuid
	FROM users
	WHERE email = ?
	OR username = ?`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return 0, err
	}

	defer sqlStmt.Close()

	var uuid int
	err = sqlStmt.QueryRow(user, user).Scan(&uuid)

	return uuid, err
}

func CheckSessionExist(session_token string) bool {
	sqlString := `SELECT uuid
	FROM users
	WHERE session_token = ?`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return false
	}

	defer sqlStmt.Close()

	_, err = sqlStmt.Exec(session_token)

	return err == nil
}

func GetUsernameBySession(session string) (string, error) {
	sqlString := `SELECT COALESCE(username,email)
	FROM users
	WHERE session_token = ?`

	var username string

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return username, err
	}

	defer sqlStmt.Close()

	err = sqlStmt.QueryRow(session).Scan(&username)

	return username, err
}

func CheckIfGroupMember(groupName string, uuid int) (bool, string) {
	sqlString := `SELECT group_id,
	role
	FROM groupMembers
	WHERE groupMembers.uuid = ? AND group_id IN (
		SELECT groups.group_id
		FROM groups
		WHERE groups.group_name = ?
	)`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return false, "lol"
	}

	defer sqlStmt.Close()

	var res string
	var dummy int

	err = sqlStmt.QueryRow(uuid, groupName).Scan(&dummy, &res)

	return err == nil, res
}

func AddGroupMessageToDB(msg data.UserMessage) error {
	//fmt.Println("stuff")
	sqlString := `INSERT INTO groupChat (group_id,gc_sender,gc_content,gc_image,creation_date)
	SELECT groups.group_id AS group_id,
	users.uuid AS gc_sender,
	? AS gc_content,
	? AS gc_image,
	? AS creation_date
	FROM groups
	JOIN users
	ON users.username = ? OR users.email = ?
	WHERE group_name = ?`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return err
	}

	defer sqlStmt.Close()

	_, err = sqlStmt.Exec(msg.Content, msg.Image, time.Now(), msg.Sender, msg.Sender, msg.Receiver)

	return err
}

func AddPrivateMessageToDB(msg data.UserMessage) error {

	sqlString := `INSERT INTO privateMessages (pmg_sender,pmg_receiver,pmg_content,pmg_image,creation_date)
	SELECT a.uuid,
	b.uuid,
	? AS pmg_content,
	? AS pmg_image,
	? AS creation_date
	FROM users AS a
	JOIN users AS b
	ON b.email = ? OR b.username = ?
	WHERE a.email = ? OR a.username = ?`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return err
	}

	defer sqlStmt.Close()

	_, err = sqlStmt.Exec(msg.Content, msg.Image, time.Now(), msg.Receiver, msg.Receiver, msg.Sender, msg.Sender)

	return err
}

func DeleteNotificationFromDB(sender, receiver int, context string) error {

	sqlString := `DELETE FROM notifications
	WHERE (uuid = ? AND sender_id = ? AND notif_context = ?)`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return err
	}

	defer sqlStmt.Close()

	_, err = sqlStmt.Exec(receiver, sender, context)

	return err
}
