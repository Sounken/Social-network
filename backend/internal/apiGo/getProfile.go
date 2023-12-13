package apiGO

import (
	"backend/internal/data"
	"backend/internal/helper"
	"encoding/json"
	"fmt"
	"net/http"
)

type qInfo struct {
	Email string `json:"slug"`
}

type profile struct {
	Email        string   `json:"email"`
	Firstname    string   `json:"first_name"`
	Lastname     string   `json:"last_name"`
	DOB          string   `json:"dob"`
	Avatar       string   `json:"avatar"`
	UserName     string   `json:"username"`
	Bio          string   `json:"bio"`
	Privacy      string   `json:"privacy"`
	Followers    []string `json:"followers"`
	Following    []string `json:"following"`
	CreatedPosts []string `json:"created_posts"`
	Groups       []string `json:"groups"`
	Status       string   `json:"status"`
	You          bool     `json:"you"`
}
type profilePrivate struct {
	Email  string `json:"email"`
	Status string `json:"status"`
}

type headBarProf struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Status   string `json:"status"`
}

func getCreatedPosts(uuid int) ([]string, error) {
	var allUserCreatedPosts []string

	sqlString := `SELECT post_title
	FROM posts
	WHERE post_author = ?`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return allUserCreatedPosts, err
	}

	defer sqlStmt.Close()

	rows, err := sqlStmt.Query(uuid)
	if err != nil {
		return allUserCreatedPosts, err
	}

	defer rows.Close()

	for rows.Next() {
		var userCreatedPost string

		err = rows.Scan(&userCreatedPost)
		if err != nil {
			return allUserCreatedPosts, err
		}
		allUserCreatedPosts = append(allUserCreatedPosts, userCreatedPost)
	}

	return allUserCreatedPosts, err
}

func getEmailById(uuid int) (string, error) {
	sqlString := `SELECT email
	FROM users
	WHERE uuid = ?`

	var result string

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return result, err
	}

	defer sqlStmt.Close()

	err = sqlStmt.QueryRow(uuid).Scan(&result)

	return result, err
}

func getEmailByUsername(name string) (string, error) {
	sqlString := `SELECT email
	FROM users
	WHERE username = ?`

	var result string

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return result, err
	}

	defer sqlStmt.Close()

	err = sqlStmt.QueryRow(name).Scan(&result)

	return result, err
}

func getProfileFromDataBase(email string) (profile, int, error) {
	var usrProfile profile
	var uuid int
	sqlString := `SELECT uuid,
	email,
	first_name,
	last_name,
	DOB,
	IFNULL(avatar, 'http://localhost:8080/images/default.jpeg'),
	IFNULL(username, 'no username'),
	IFNULL(bio, 'no bio'),
	privacy
	FROM users
	WHERE email = ?`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return usrProfile, uuid, err
	}

	defer sqlStmt.Close()

	err = sqlStmt.QueryRow(email).Scan(&uuid, &usrProfile.Email, &usrProfile.Firstname, &usrProfile.Lastname, &usrProfile.DOB, &usrProfile.Avatar, &usrProfile.UserName, &usrProfile.Bio, &usrProfile.Privacy)

	return usrProfile, uuid, err
}

func GetHeadBar(w http.ResponseWriter, r *http.Request) {
	helper.EnableCors(&w)

	if r.Method == http.MethodPost {
		var usrDat headBarProf
		uuid, err := helper.GetIdBySession(w, r)
		if err != nil {
			fmt.Println("get head bar prof error", err)
			helper.WriteResponse(w, "session_error")
			return
		}

		sqlString := `SELECT IFNULL(avatar, 'http://localhost:8080/images/default.jpeg'),
		COALESCE(username, email)
		FROM users
		WHERE uuid = ?`

		sqlStmt, err := data.DB.Prepare(sqlString)
		if err != nil {
			helper.WriteResponse(w, "database_error")
			return
		}

		defer sqlStmt.Close()

		err = sqlStmt.QueryRow(uuid).Scan(&usrDat.Avatar, &usrDat.Username)
		if err != nil {
			fmt.Println("get headbar db error", err)
			helper.WriteResponse(w, "database_error")
		}

		usrDat.Status = "success"
		usrDatJson, err := json.Marshal(usrDat)
		if err != nil {
			fmt.Println("marshalling error", err)
			helper.WriteResponse(w, "marshalling_error")
		}
		//fmt.Println("profile get success")
		w.Header().Set("Content-Type", "application/json")
		w.Write(usrDatJson)
		fmt.Println("headbar donezo")
	}
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	helper.EnableCors(&w)

	if r.Method == http.MethodPost {
		var eInfo qInfo
		var yourProf = false
		yourId, err := helper.GetIdBySession(w, r)
		if err != nil {
			fmt.Println("get profile session error", err)
			helper.WriteResponse(w, "session_error")
		}
		//fmt.Println(r.URL.Path)
		err = json.NewDecoder(r.Body).Decode(&eInfo.Email)
		if err != nil {
			fmt.Println(err)
		}
		//voff
		if eInfo.Email == "voff" {
			eInfo.Email, err = getEmailById(yourId)
			if err != nil {
				helper.WriteResponse(w, "database_error")
				fmt.Println("get your profile db error", err)
				return
			}
			yourProf = true

		} else {
			emailExists := helper.CheckIfStringExist("users", "email", eInfo.Email)
			if !emailExists {
				if !helper.CheckIfStringExist("users", "username", eInfo.Email) {
					helper.WriteResponse(w, "user_does_not_exist")
					return
				}
				eInfo.Email, err = getEmailByUsername(eInfo.Email)
				if err != nil {
					helper.WriteResponse(w, "error_should_not_be_possible")
					fmt.Println("error should not be possible", err)
					return
				}

			}
		}
		emailuuid, err := helper.GetuuidFromEmailOrUsername(eInfo.Email)
		if err != nil {
			helper.WriteResponse(w, "email_does_not_exist")
			return
		}
		if yourId == emailuuid {
			yourProf = true
		}
		usrProfile, uuid, err := getProfileFromDataBase(eInfo.Email)
		if err != nil {
			fmt.Println("database err", err)
			helper.WriteResponse(w, "profile_error")
			return
		}
		if (usrProfile.Privacy == "private" && !checkDBIfFollowing(yourId, uuid)) && !yourProf {
			var privProf profilePrivate
			privProf.Email = usrProfile.Email
			privProf.Status = "private"
			privProfJson, err := json.Marshal(privProf)
			if err != nil {
				helper.WriteResponse(w, "marshalling_error")
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(privProfJson)
			return
		}
		usrProfile.Following, err = helper.GetFollowing(uuid)
		if err != nil {
			fmt.Println("following error", err)
			helper.WriteResponse(w, "following_error")
			return
		}

		usrProfile.Followers, err = helper.GetFollowers(uuid)
		if err != nil {
			fmt.Println("followers error", err)
			helper.WriteResponse(w, "followers_error")
			return
		}

		usrProfile.Groups, err = helper.GetYourGroups(uuid)
		if err != nil {
			fmt.Println("groups error", err)
			helper.WriteResponse(w, "groups_error")
			return
		}

		usrProfile.CreatedPosts, err = getCreatedPosts(uuid)
		if err != nil {
			fmt.Println("created posts error", err)
			helper.WriteResponse(w, "database_error")
			return
		}
		usrProfile.You = yourProf
		usrProfile.Status = "success"
		usrProfileJson, err := json.Marshal(usrProfile)
		if err != nil {
			fmt.Println("marshalling error", err)
			helper.WriteResponse(w, "marshalling_error")
		}
		//fmt.Println("profile get success")
		w.Header().Set("Content-Type", "application/json")
		w.Write(usrProfileJson)
	}
}
