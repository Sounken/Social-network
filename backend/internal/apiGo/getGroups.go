package apiGO

import (
	"backend/internal/data"
	"backend/internal/helper"
	"encoding/json"
	"fmt"
	"net/http"
)

type groups struct {
	Groups []string `json:"groups"`
	Status string   `json:"status"`
}

func getGroupnamesFromDB() (groups, error) {
	var groupDat groups
	sqlString := `SELECT group_name FROM groups`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return groupDat, err
	}

	defer sqlStmt.Close()

	rows, err := sqlStmt.Query()
	if err != nil {
		return groupDat, err
	}

	defer rows.Close()

	for rows.Next() {
		var grp string
		err = rows.Scan(&grp)
		if err != nil {
			return groupDat, err
		}
		groupDat.Groups = append(groupDat.Groups, grp)
	}
	return groupDat, nil
}

func GetGroupnames(w http.ResponseWriter, r *http.Request) {
	helper.EnableCors(&w)
	if r.Method == http.MethodPost {
		groupData, err := getGroupnamesFromDB()
		if err != nil {
			helper.WriteResponse(w, "database_error")
			return
		}
		groupData.Status = "success"
		groupDataJson, err := json.Marshal(groupData)
		if err != nil {
			fmt.Println("marshalling error", err)
			helper.WriteResponse(w, "marshalling_error")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(groupDataJson)
	}
}
