package apiGO

import (
	"backend/internal/data"
	"backend/internal/helper"
	"encoding/json"
	"fmt"
	_ "image/png"
	"net/http"
)

func getStockImages() ([]string, error) {
	sqlString := `SELECT image_path FROM stockImages`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return []string{}, err
	}

	defer sqlStmt.Close()

	rows, err := sqlStmt.Query()
	var imgPaths []string
	if err != nil {
		fmt.Println("query error", err)
		return imgPaths, err
	}

	defer rows.Close()

	for rows.Next() {
		var imgPath string
		err = rows.Scan(&imgPath)
		if err != nil {
			fmt.Println(err)
		}
		imgPaths = append(imgPaths, imgPath)
	}
	return imgPaths, nil
}

func getUserImages(uuid int) ([]string, error) {
	var imgPaths []string

	sqlString := `SELECT image_path FROM userImages WHERE uuid = ?`

	sqlStmt, err := data.DB.Prepare(sqlString)
	if err != nil {
		return imgPaths, err
	}

	defer sqlStmt.Close()

	rows, err := sqlStmt.Query(uuid)
	if err != nil {
		fmt.Println("query error usrImg", err)
		return imgPaths, err
	}

	defer rows.Close()

	for rows.Next() {
		var imgPath string
		err = rows.Scan(&imgPath)
		if err != nil {
			fmt.Println(err)
		}
		imgPaths = append(imgPaths, imgPath)
	}
	return imgPaths, nil
}

type images struct {
	StockImages []string `json:"stock_images"`
	UserImages  []string `json:"user_images"`
	Status      string   `json:"status"`
}

func GetYourImages(w http.ResponseWriter, r *http.Request) {

	helper.EnableCors(&w)

	if r.Method == http.MethodPost {

		var imgData images
		uuid, err := helper.GetIdBySession(w, r)
		if err != nil {
			fmt.Println("session errir", err)
			helper.WriteResponse(w, "incorrect_session")
			return
		}
		imgData.UserImages, err = getUserImages(uuid)
		if err != nil {
			fmt.Println("userimg err", err)
		}
		imgData.StockImages, err = getStockImages()
		if err != nil {
			fmt.Println("stock img err", err)
		}
		imgData.Status = "success"
		imgDataJson, err := json.Marshal(imgData)
		if err != nil {
			fmt.Println("marshalling error", err)
			helper.WriteResponse(w, "marshalling_error")
		}
		//fmt.Println("success")
		w.Header().Set("Content-Type", "application/json")
		w.Write(imgDataJson)
	}

}
