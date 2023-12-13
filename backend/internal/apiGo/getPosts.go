package apiGO

import (
	"backend/internal/data"
	"backend/internal/helper"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type allposts struct {
	Posts            []posts `json:"posts"`
	PrivatePosts     []posts `json:"private_posts"`
	SemiPrivatePosts []posts `json:"semi_private_posts"`
	Status           string  `json:"status"`
}
type posts struct {
	PostId       int       `json:"post_id"`
	Author       string    `json:"author"`
	Image        string    `json:"image"`
	CreationDate time.Time `json:"creation_date"`
	Content      string    `json:"content"`
	Title        string    `json:"title"`
}

func gatherPosts() ([]posts, error) {
	var pData []posts
	sqlString := `SELECT post_id,
COALESCE(users.username, users.email),
IFNULL(post_image, 'http://localhost:8080/images/default.jpeg'),
creation_date,
post_content,
post_title
FROM posts
JOIN users
ON users.uuid = posts.post_author
WHERE post_privacy = 'public';`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return pData, err
	}

	defer sqlStmt.Close()

	rows, err := sqlStmt.Query()
	if err != nil {
		fmt.Println("query error", err)
		return pData, err
	}

	defer rows.Close()

	for rows.Next() {
		var normPost posts
		//var authId int
		err = rows.Scan(&normPost.PostId, &normPost.Author, &normPost.Image, &normPost.CreationDate, &normPost.Content, &normPost.Title)
		if err != nil {
			fmt.Println("scan normal Post err", err)
		}
		pData = append(pData, normPost)
	}

	return pData, nil
}

func gatherSemiPrivatePosts(uuid int) ([]posts, error) {
	var pData []posts
	sqlString := `SELECT post_id,
COALESCE(users.username, users.email),
IFNULL(post_image, 'http://localhost:8080/images/default.jpeg'),
creation_date,
post_content,
post_title
FROM posts
JOIN users
ON users.uuid = posts.post_author
WHERE post_privacy = 'semi-private' AND post_id IN (
SELECT post_id
FROM allowedUsers
WHERE uuid = ?
)
OR (posts.post_author = ? AND post_privacy = 'semi-private');`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return pData, err
	}

	defer sqlStmt.Close()

	rows, err := sqlStmt.Query(uuid, uuid)
	if err != nil {
		fmt.Println("query error", err)
		return pData, err
	}

	defer rows.Close()

	for rows.Next() {
		var semiPost posts
		//var authId int
		err = rows.Scan(&semiPost.PostId, &semiPost.Author, &semiPost.Image, &semiPost.CreationDate, &semiPost.Content, &semiPost.Title)
		if err != nil {
			fmt.Println("scan semi private Post err", err)
		}
		pData = append(pData, semiPost)
	}

	return pData, nil
}

func gatherPrivatePosts(uuid int) ([]posts, error) {
	var pData []posts
	sqlString := `SELECT post_id,
COALESCE(users.username, users.email),
IFNULL(post_image, 'http://localhost:8080/images/default.jpeg'),
creation_date,
post_content,
post_title
FROM posts
JOIN users
ON users.uuid = posts.post_author
WHERE post_privacy = 'private' AND post_author IN (
SELECT uuid
FROM followers
WHERE follower_id = ?
)
OR (posts.post_author = ? AND post_privacy = 'private');`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return pData, err
	}

	defer sqlStmt.Close()

	rows, err := sqlStmt.Query(uuid, uuid)
	if err != nil {
		fmt.Println("query error", err)
		return pData, err
	}

	defer rows.Close()

	for rows.Next() {
		var privPost posts
		//var authId int
		err = rows.Scan(&privPost.PostId, &privPost.Author, &privPost.Image, &privPost.CreationDate, &privPost.Content, &privPost.Title)
		if err != nil {
			fmt.Println("scan privatePost err", err)
		}
		pData = append(pData, privPost)
	}

	return pData, nil
}


func GetPosts(w http.ResponseWriter, r *http.Request) {
	helper.EnableCors(&w)
	if r.Method == http.MethodPost {
		uuid, err := helper.GetIdBySession(w, r)
		if err != nil {
			helper.WriteResponse(w, "incorrect_session")
			return
		}
		var pData allposts
		pData.Posts, err = gatherPosts()
		if err != nil {
			helper.WriteResponse(w, "gatherPosts_error")
			return
		}
		pData.PrivatePosts, err = gatherPrivatePosts(uuid)
		if err != nil {
			helper.WriteResponse(w, "gatherPrivatePosts_error")
			return
		}
		pData.SemiPrivatePosts, err = gatherSemiPrivatePosts(uuid)
		if err != nil {
			helper.WriteResponse(w, "gatherSemiPrivatePosts_error")
			return
		}
		pData.Status = "success"
		pDataJson, err := json.Marshal(pData)
		if err != nil {
			fmt.Println("marshalling error", err)
			helper.WriteResponse(w, "marshalling_error")
			return
		}
		fmt.Println("success")
		w.Header().Set("Content-Type", "application/json")
		w.Write(pDataJson)
	}
}
