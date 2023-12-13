package apiGO

import (
	"backend/internal/data"
	"backend/internal/helper"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type post struct {
	Type      string   `json:"type"`
	Privacy   string   `json:"privacy"`
	Allowed   []string `json:"allowed_users"`
	GroupName string   `json:"group_name"`
	Image     string   `json:"image"`
	Content   string   `json:"content"`
	Title     string   `json:"title"`
	Author    int
}

func addAllowedUsersToDB(users []string, postTitle string) error {
	sqlString := `INSERT INTO allowedUsers(post_id,uuid)
	SELECT posts.post_id AS post_id,
	users.uuid
	FROM posts
	JOIN users
	ON users.username = ? OR users.email = ?
	WHERE post_title = ?
	ORDER BY posts.post_id DESC
	LIMIT 1`
	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return err
	}

	defer sqlStmt.Close()

	for _, user := range users {
		sqlStmt.Exec(user, user, postTitle)
	}
	return err
}

func addPostToTable(postData post) error {
	sqlStmt, err := data.DB.Prepare(`INSERT INTO posts(post_author,post_privacy,post_image,creation_date,post_content,post_title) values(?,?,?,?,?,?)`)
	if err != nil {
		return err
	}
	defer sqlStmt.Close()

	_, err = sqlStmt.Exec(postData.Author, postData.Privacy, postData.Image, time.Now(), postData.Content, postData.Title)
	return err
}

func addGroupPostToTable(postData post) error {
	sqlString := `INSERT INTO groupPosts(gpost_author,group_id,gpost_image,creation_date,gpost_content,gpost_title)
	SELECT ? AS gpost_author,
	group_id,
	? AS gpost_image,
	? AS creation_date,
	? AS gpost_content,
	? AS gpost_title
	FROM groups
	WHERE group_name = ?`
	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return err
	}
	defer sqlStmt.Close()

	_, err = sqlStmt.Exec(postData.Author, postData.Image, time.Now(), postData.Content, postData.Title, postData.GroupName)
	return err
}

// api that takes in a json object and tries to add it in either the group_post or the post table
func PostApi(w http.ResponseWriter, r *http.Request) {
	helper.EnableCors(&w)
	if r.Method == http.MethodPost {
		var pData post
		err := json.NewDecoder(r.Body).Decode(&pData)
		//fmt.Println("post data", pData)
		if err != nil {
			helper.WriteResponse(w, "decoding_error")
			fmt.Println(err)
			return
		}
		uuid, err := helper.GetIdBySession(w, r)
		if err != nil {
			helper.WriteResponse(w, "incorrect_session")
			return
		}
		pData.Author = uuid
		if err != nil {
			helper.WriteResponse(w, "user_not_found")
			return
		}
		if pData.Type == "post" {
			err = addPostToTable(pData)
			if err != nil {
				helper.WriteResponse(w, "incorrect_input")
				return
			}
			if pData.Privacy == "semi-private" {
				fmt.Println("got to semi private")
				err = addAllowedUsersToDB(pData.Allowed, pData.Title)
				if err != nil {
					helper.WriteResponse(w, "could_not_add_users")
					return
				}
			}
		} else if pData.Type == "group_post" {
			isMember, _ := helper.CheckIfGroupMember(pData.GroupName, uuid)
			if !isMember {
				helper.WriteResponse(w, "not_a_member")
				return
			}
			err = addGroupPostToTable(pData)
			if err != nil {
				helper.WriteResponse(w, "incorrect_input")
				return
			}
		}
		helper.WriteResponse(w, "success")

	}
}
