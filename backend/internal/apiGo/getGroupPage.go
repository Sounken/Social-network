package apiGO

import (
	"backend/internal/data"
	"backend/internal/helper"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type event struct {
	Id            int       `json:"event_id"`
	Author        string    `json:"author"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Created       time.Time `json:"created"`
	Date          time.Time `json:"date"`
	Options       []string  `json:"options"`
	Answers       []int     `json:"answers"`
	AlreadyChosen bool      `json:"already_chosen"`
}

type chatMessage struct {
}

type groupMember struct {
	Username   string `json:"username"`
	MemberType string `json:"member_type"`
}

type groupPage struct {
	Members      []groupMember  `json:"members"`
	Events       []event        `json:"events"`
	ChatMessages []chatMessage  `json:"chat_messages"`
	GroupPosts   []posts        `json:"group_posts"`
	MemberType   string         `json:"member_type"`
	JoinRequests []notification `json:"join_requests"`
	Status       string         `json:"status"`
}

func gatherGroupMembers(groupName string) ([]groupMember, error) {
	var allGroupMembers []groupMember

	sqlStmt, err := data.DB.Prepare(`SELECT COALESCE(users.username,users.email),
	role
	FROM groupMembers
	JOIN users
	ON users.uuid = groupMembers.uuid
	WHERE groupMembers.group_id IN (
		SELECT groups.group_id
		FROM groups
		WHERE groups.group_name = ?
	)`)
	if err != nil {
		return allGroupMembers, err
	}

	defer sqlStmt.Close()

	rows, err := sqlStmt.Query(groupName)
	if err != nil {
		return allGroupMembers, err
	}

	defer rows.Close()

	for rows.Next() {
		var gMember groupMember
		err = rows.Scan(&gMember.Username, &gMember.MemberType)
		if err != nil {
			return allGroupMembers, err
		}
		allGroupMembers = append(allGroupMembers, gMember)
	}

	return allGroupMembers, err
}

func gatherEventChoices(allEvents []event) ([]event, error) {

	sqlStmt, err := data.DB.Prepare(`SELECT COUNT(*)
	FROM eventOptionChoices
	WHERE event_id = ? AND choice = ?`)
	if err != nil {
		return allEvents, err
	}

	defer sqlStmt.Close()
	for index, singleEvent := range allEvents {
		var allVotes []int
		for _, choice := range singleEvent.Options {
			var choiceVotes int
			err = sqlStmt.QueryRow(singleEvent.Id, choice).Scan(&choiceVotes)
			if err != nil {
				return allEvents, err
			}
			allVotes = append(allVotes, choiceVotes)
		}
		allEvents[index].Answers = allVotes
	}
	return allEvents, err
}

func gatherGroupEvents(groupName string, uuid int) ([]event, error) {
	var allEvents []event

	sqlStmt, err := data.DB.Prepare(`SELECT events.event_id,
	COALESCE(users.username, users.email),
	event_title,
	event_content,
	creation_date,
	event_date,
	options,
	COALESCE((SELECT 
		eventOptionChoices.uuid
		FROM eventOptionChoices
		WHERE events.event_id = eventOptionChoices.event_id AND eventOptionChoices.uuid = ?), -1)
	FROM events
	JOIN users
	ON events.event_author = users.uuid
	WHERE events.group_id IN (
		SELECT groups.group_id
		FROM groups
		WHERE groups.group_name = ?
	)`)
	if err != nil {
		fmt.Println(err)
		return allEvents, err
	}
	defer sqlStmt.Close()
	rows, err := sqlStmt.Query(uuid, groupName)
	if err != nil {
		fmt.Println(err)
		return allEvents, err
	}
	defer rows.Close()

	for rows.Next() {
		var gEvent event
		var optionsInString string
		var hasChoosen = -1
		err = rows.Scan(&gEvent.Id, &gEvent.Author, &gEvent.Title, &gEvent.Content, &gEvent.Created, &gEvent.Date, &optionsInString, &hasChoosen)
		if err != nil {
			fmt.Println(err)
			return allEvents, err
		}
		gEvent.Options = strings.Split(optionsInString, data.SEPERATOR)
		if hasChoosen != -1 {
			gEvent.AlreadyChosen = true
		}

		allEvents = append(allEvents, gEvent)
		//fmt.Println(gEvent)
	}
	return allEvents, err
}

func gatherGroupJoinRequests(groupName string) ([]notification, error) {
	var allJoinRequests []notification

	sqlStmt, err := data.DB.Prepare(`SELECT notif_content,
	creation_date,
	COALESCE(users.username, users.email),
	notif_type,
	notif_context
	FROM notifications
	JOIN users
	ON notifications.sender_id = users.uuid
	WHERE notifications.notif_type = 'group_join_request' AND notifications.notif_context = ?`)
	if err != nil {
		return allJoinRequests, err
	}
	defer sqlStmt.Close()

	rows, err := sqlStmt.Query(groupName)
	if err != nil {
		return allJoinRequests, err
	}
	defer rows.Close()

	for rows.Next() {
		var joinRequest notification

		err = rows.Scan(&joinRequest.Content, &joinRequest.Created, &joinRequest.Sender, &joinRequest.Type, &joinRequest.Context)
		if err != nil {
			return allJoinRequests, err
		}
		allJoinRequests = append(allJoinRequests, joinRequest)
	}
	return allJoinRequests, err
}

func gatherGroupPageData(groupName, memberType string, uuid int) (groupPage, error) {
	var groupPageData groupPage
	var err error
	groupPageData.GroupPosts, err = gatherGroupPosts(groupName)
	if err != nil {
		fmt.Println("Error in gatherGroupPosts: ", err)
		return groupPageData, err
	}
	groupPageData.Members, err = gatherGroupMembers(groupName)
	if err != nil {
		fmt.Println("Error in gatherGroupMembers: ", err)
		return groupPageData, err
	}
	groupPageData.Events, err = gatherGroupEvents(groupName, uuid)
	if err != nil {
		fmt.Println("Error in gatherGroupEvents: ", err)
		return groupPageData, err
	}
	groupPageData.Events, err = gatherEventChoices(groupPageData.Events)
	if memberType != "creator" {
		return groupPageData, err
	}
	groupPageData.JoinRequests, err = gatherGroupJoinRequests(groupName)
	if err != nil {
		fmt.Println("join request err", err)
	}

	return groupPageData, err
}

func gatherGroupPosts(groupName string) ([]posts, error) {
	var allGroupPosts []posts
	sqlStmt, err := data.DB.Prepare(`SELECT gpost_id,
	COALESCE(users.username, users.email),
	IFNULL(gpost_image, 'http://localhost:8080/images/default.jpeg'),
	creation_date,
	gpost_content,
	gpost_title
	FROM groupPosts
	JOIN users
	ON users.uuid = gpost_author
	WHERE groupPosts.group_id IN (
		SELECT groups.group_id
		FROM groups
		WHERE group_name = ?
	)`)
	if err != nil {
		return allGroupPosts, err
	}
	defer sqlStmt.Close()

	rows, err := sqlStmt.Query(groupName)
	if err != nil {
		return allGroupPosts, err
	}
	defer rows.Close()

	for rows.Next() {
		var groupPost posts

		err = rows.Scan(&groupPost.PostId, &groupPost.Author, &groupPost.Image, &groupPost.CreationDate, &groupPost.Content, &groupPost.Title)
		if err != nil {
			return allGroupPosts, err
		}
		allGroupPosts = append(allGroupPosts, groupPost)
	}

	return allGroupPosts, err
}

func GetGroupPage(w http.ResponseWriter, r *http.Request) {
	helper.EnableCors(&w)

	if r.Method == http.MethodPost {
		var groupName string
		err := json.NewDecoder(r.Body).Decode(&groupName)
		if err != nil {
			helper.WriteResponse(w, "decoding_error")
			fmt.Println("decoding", err)
			return
		}
		groupExist := helper.CheckIfStringExist("groups", "group_name", groupName)
		if !groupExist {
			helper.WriteResponse(w, "group_does_not_exist")
			fmt.Println("group not exist", err)
			return
		}
		uuid, err := helper.GetIdBySession(w, r)
		if err != nil {
			helper.WriteResponse(w, "session_error")
			fmt.Println(err)
			return
		}
		isMember, memberType := helper.CheckIfGroupMember(groupName, uuid)
		if !isMember {
			helper.WriteResponse(w, "not_a_member")
			fmt.Println("not a member group page", groupName, uuid)
			return
		}
		groupPageData, err := gatherGroupPageData(groupName, memberType, uuid)
		if err != nil {
			helper.WriteResponse(w, "database_error")
			fmt.Println("gathergroupPageData", err)
			return
		}
		groupPageData.Status = "success"
		groupPageData.MemberType = memberType
		pDataJson, err := json.Marshal(groupPageData)
		if err != nil {
			fmt.Println("marshalling error", err)
			helper.WriteResponse(w, "marshalling_error")
			fmt.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(pDataJson)

	}
}
