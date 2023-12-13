package apiGO

import (
	"backend/internal/helper"
	"encoding/json"
	"fmt"
	"net/http"
)

type settings struct {
	Image    string `json:"avatar"`
	Privacy  string `json:"privacy"`
	Nickname string `json:"nickname"`
	AboutMe  string `json:"aboutMe"`
}

func UpdateSettings(w http.ResponseWriter, r *http.Request) {
	helper.EnableCors(&w)
	if r.Method == http.MethodPost {
		uid, err := helper.GetIdBySession(w, r)
		if err != nil {
			helper.WriteResponse(w, "cookie")
			return
		}

		var settData settings

		err = json.NewDecoder(r.Body).Decode(&settData)
		fmt.Println("data to update", settData)

		if err != nil {
			fmt.Println(err)
			helper.WriteResponse(w, "json")
			return
		}

		if len(settData.Nickname) > 0 {
			nameTaken := helper.CheckIfStringExist("users", "username", settData.Nickname)
			if nameTaken {
				helper.WriteResponse(w, "name_taken")
				return
			}
			err = helper.UpdateTableColumnStringById("users", settData.Nickname, "username", uid)
			if err != nil {
				helper.WriteResponse(w, "username_error")
				fmt.Println(err)
				return
			}
		}

		if len(settData.AboutMe) > 0 {
			err = helper.UpdateTableColumnStringById("users", settData.AboutMe, "bio", uid)
			if err != nil {
				helper.WriteResponse(w, "aboutMe")
				return
			}
		}

		if len(settData.Image) > 0 {
			err = helper.UpdateTableColumnStringById("users", settData.Image, "avatar", uid)
			if err != nil {
				helper.WriteResponse(w, "image")
				return
			}
		}

		if settData.Privacy == "private" || settData.Privacy == "public" {
			err = helper.UpdateTableColumnStringById("users", settData.Privacy, "privacy", uid)
			if err != nil {
				helper.WriteResponse(w, "privacy")
				return
			}
		}

		helper.WriteResponse(w, "success")
	}
}
