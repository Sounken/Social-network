package apiGO

import (
	"backend/internal/data"
	"backend/internal/helper"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type addComment struct {
	Content   string `json:"content"`
	Image     string `json:"image"`
	PostId    int    `json:"post_id"`
	PostType  string `json:"post_type"`
	GroupName string `json:"group_name"`
}

func addCommentToDB(uuid, postId int, content, image string) error {

	sqlString := `INSERT INTO comments(comment_author,post_id,comment_content,comment_image,creation_date)
	VALUES(?,?,?,?,?)`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return err
	}

	defer sqlStmt.Close()

	_, err = sqlStmt.Exec(uuid, postId, content, image, time.Now())
	return err
}

func addGroupCommentToDB(uuid, postId int, content, image, groupName string) error {
	sqlString := `INSERT INTO groupPostComments(gpost_id,group_id,gpcomment_author,gpcomment_content,gpcomment_image,creation_date)
	SELECT ? AS gpost_id,
	g.group_id AS group_id,
	? AS gpcomment_author,
	? AS gpcomment_content,
	? AS gpcomment_image,
	? AS creation_date
	FROM groups AS g
	WHERE ? = g.group_name`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return err
	}

	defer sqlStmt.Close()

	_, err = sqlStmt.Exec(postId, uuid, content, image, time.Now(), groupName)
	return err
}

func AddComment(w http.ResponseWriter, r *http.Request) {

	helper.EnableCors(&w)

	if r.Method != http.MethodPost {
		return
	}

	var commentData addComment
	err := json.NewDecoder(r.Body).Decode(&commentData)
	if err != nil {
		helper.WriteResponse(w, "decoding_error")
		return
	}

	uuid, err := helper.GetIdBySession(w, r)
	if err != nil {
		helper.WriteResponse(w, "session_error")
	}
	if commentData.PostType != "post" && commentData.PostType != "group_post" {
		helper.WriteResponse(w, "incorrect_post_type")
		return
	}

	if commentData.PostType == "post" {

		err = addCommentToDB(uuid, commentData.PostId, commentData.Content, commentData.Image)
	}
	if commentData.PostType == "group_post" {

		groupExists := helper.CheckIfStringExist("groups", "group_name", commentData.GroupName)
		if !groupExists {
			fmt.Println("group does not exist (adding comment)", commentData, err)
			helper.WriteResponse(w, "group_does_not_exist")
			return
		}

		isGroupMember, _ := helper.CheckIfGroupMember(commentData.GroupName, uuid)
		if !isGroupMember {
			fmt.Println("is not a group member (adding comment)", err)
			helper.WriteResponse(w, "not_a_member")
			return
		}

		err = addGroupCommentToDB(uuid, commentData.PostId, commentData.Content, commentData.Image, commentData.GroupName)

	}

	if err != nil {
		helper.WriteResponse(w, "database_error")
		return
	}

	helper.WriteResponse(w, "success")

}
