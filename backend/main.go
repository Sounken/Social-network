package main

import (
	apiGO "backend/internal/apiGo"
	"backend/internal/data"
	socket "backend/internal/websocket"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
)

func runMigrations() {
	migrations := &migrate.FileMigrationSource{
		Dir: "internal/migrations",
	}
	fmt.Println("mig", *migrations)
	n, err := migrate.Exec(data.DB, "sqlite3", migrations, migrate.Up)
	if err != nil {
		fmt.Println("failed migration", err)
	}
	fmt.Printf("Applied %d migrations!\n", n)
}

func testhandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Wah bien joué</h1>")
	fmt.Fprintf(w, "<title>Go</title>")
	fmt.Fprintf(w, "<img src=images/feelsgoodman.png>")
}

func setupApi() {
	http.HandleFunc("/api/register", apiGO.Register)
	http.HandleFunc("/api/login", apiGO.Login)
	http.HandleFunc("/api/updateSettings", apiGO.UpdateSettings)
	http.HandleFunc("/api/createPost", apiGO.PostApi)
	http.HandleFunc("/api/cookie", apiGO.CheckCookie)
	http.HandleFunc("/api/followCheck", apiGO.CheckFollowing)

	http.HandleFunc("/api/getYourImages", apiGO.GetYourImages)
	http.HandleFunc("/api/getGroupnames", apiGO.GetGroupnames)
	http.HandleFunc("/api/getUsernames", apiGO.GetUsernames)
	http.HandleFunc("/api/getPosts", apiGO.GetPosts)
	http.HandleFunc("/api/getProfile", apiGO.GetProfile)
	http.HandleFunc("/api/getNotifications", apiGO.GetAllNotifications)
	http.HandleFunc("/api/getPrivateMessages", apiGO.GetPrivateMessages)
	http.HandleFunc("/api/getGroupChat", apiGO.GetGroupChatMessages)
	http.HandleFunc("/api/getComments", apiGO.GetComments)
	http.HandleFunc("/api/getHeadbar", apiGO.GetHeadBar)
	http.HandleFunc("/api/getGroupPage", apiGO.GetGroupPage)

	http.HandleFunc("/api/addOrRemoveFollow", apiGO.AddOrRemoveFollow)
	http.HandleFunc("/api/addGroup", apiGO.AddGroup)
	http.HandleFunc("/api/addEvent", apiGO.AddEvent)
	http.HandleFunc("/api/addComment", apiGO.AddComment)
	http.HandleFunc("/api/addGroupMemberToChat", socket.AddGroupMemberToChat)
	http.HandleFunc("/api/addPost", apiGO.PostApi)
	http.HandleFunc("/api/addImage", apiGO.AddImage)

	http.HandleFunc("/api/handleGroupJoinRequest", apiGO.HandleGroupJoinRequest)
	http.HandleFunc("/api/handleGroupInvite", apiGO.HandleGroupInvite)
	http.HandleFunc("/api/handleFollowRequest", apiGO.HandleFollowRequest)
	http.HandleFunc("/api/sendGroupJoinRequest", apiGO.SendGroupJoinRequestNotification)
	http.HandleFunc("/api/sendGroupInvite", apiGO.SendGroupInviteNotification)
	http.HandleFunc("/api/handleEventAnswer", apiGO.HandleEventAnswer)
	http.HandleFunc("/api/deleteNotification", apiGO.DeleteNotification)

	http.HandleFunc("/api/ws", socket.WsEndpoint)

	img := http.FileServer(http.Dir("internal/images"))
	http.Handle("/images/", http.StripPrefix("/images/", img))
	http.HandleFunc("/test", testhandler)
}

func main() {
	//var db *sql.DB
	_, err := os.Stat("./database/network.db")
	if err != nil {
		fmt.Println("dir not found", err)
		return
	}
	data.DB, err = sql.Open("sqlite3", "./database/network.db")
	if err != nil {
		fmt.Println("wrong bumbo", err)
		return
	}
	defer data.DB.Close()
	runMigrations()

	srv := &http.Server{
		Addr: ":8080",
	}
	setupApi()

	fmt.Println("Starting application on port " + srv.Addr)
	if srv.ListenAndServe() != nil {
		log.Fatalf("%v - Internal Server Error", http.StatusInternalServerError)
	}
}
