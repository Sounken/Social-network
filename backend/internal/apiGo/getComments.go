package apiGO

import (
	"backend/internal/data"
	"backend/internal/helper"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type comment struct {
	Author  string    `json:"author"`
	Content string    `json:"content"`
	Image   string    `json:"image_path"`
	Created time.Time `json:"created"`
}
type comments struct {
	AllComments []comment `json:"comments"`
	Status      string    `json:"status"`
}

type requestInfo struct {
	PostId    int    `json:"post_id"`
	Type      string `json:"type"`
	GroupName string `json:"group_name"`
}

func gatherGroupCommentsFromDB(postId int, groupName string) (comments, error) {
	var commentsData comments
	sqlString := `SELECT COALESCE(users.username, users.email),
	gpcomment_content,
	IFNULL(gpcomment_image, 'http://localhost:8080/images/default.jpeg'),
	creation_date
	FROM groupPostComments
	JOIN users
	ON gpcomment_author = users.uuid
	WHERE gpost_id = ? AND group_id IN (
		SELECT g.group_id
		FROM groups AS g
		WHERE g.group_name = ?
	)`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return commentsData, err
	}

	defer sqlStmt.Close()

	rows, err := sqlStmt.Query(postId, groupName)
	if err != nil {
		return commentsData, err
	}

	defer rows.Close()

	for rows.Next() {
		var cData comment
		err = rows.Scan(&cData.Author, &cData.Content, &cData.Image, &cData.Created)
		if err != nil {
			return commentsData, err
		}
		commentsData.AllComments = append(commentsData.AllComments, cData)
	}
	return commentsData, err
}

func gatherCommentsFromDB(postId int) (comments, error) {
	var commentsData comments
	sqlString := `SELECT COALESCE(users.username, users.email),
	comment_content,
	IFNULL(comment_image, 'http://localhost:8080/images/default.jpeg'),
	creation_date
	FROM comments
	JOIN users
	ON comment_author = users.uuid
	WHERE post_id = ?`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return commentsData, err
	}

	defer sqlStmt.Close()

	rows, err := sqlStmt.Query(postId)
	if err != nil {
		return commentsData, err
	}

	defer rows.Close()

	for rows.Next() {
		var cData comment
		err = rows.Scan(&cData.Author, &cData.Content, &cData.Image, &cData.Created)
		if err != nil {
			return commentsData, err
		}
		commentsData.AllComments = append(commentsData.AllComments, cData)
	}
	return commentsData, err
}

func GetComments(w http.ResponseWriter, r *http.Request) {
	helper.EnableCors(&w)
	if r.Method == http.MethodPost {
		var idAndType requestInfo
		err := json.NewDecoder(r.Body).Decode(&idAndType)
		if err != nil {
			helper.WriteResponse(w, "decoding_error")
			fmt.Println("decoding error", err)
			return
		}
		if idAndType.Type != "normal_comments" && idAndType.Type != "group_comments" {
			helper.WriteResponse(w, "no_type_selected")
			return
		}

		var commentData comments

		if idAndType.Type == "normal_comments" {
			commentData, err = gatherCommentsFromDB(idAndType.PostId)
			if err != nil {
				fmt.Println("database error gather comments", err)
				helper.WriteResponse(w, "database_error")
				return
			}
		}
		if idAndType.Type == "group_comments" {
			groupExists := helper.CheckIfStringExist("groups", "group_name", idAndType.GroupName)
			if !groupExists {
				helper.WriteResponse(w, "group_does_not_exist")
				return
			}
			uuid, err := helper.GetIdBySession(w, r)
			if err != nil {
				helper.WriteResponse(w, "session_error")
				return
			}
			isGroupMember, _ := helper.CheckIfGroupMember(idAndType.GroupName, uuid)
			if !isGroupMember {
				helper.WriteResponse(w, "not_a_member")
				return
			}
			commentData, err = gatherGroupCommentsFromDB(idAndType.PostId, idAndType.GroupName)
			if err != nil {
				fmt.Println("database error gather comments", err)
				helper.WriteResponse(w, "database_error")
				return
			}
		}

		commentData.Status = "success"
		commentDataJson, err := json.Marshal(commentData)
		if err != nil {
			fmt.Println("marshalling error", err)
			helper.WriteResponse(w, "marshalling_error")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(commentDataJson)
	}
}
