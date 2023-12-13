package data

import (
	"database/sql"
	"time"
)

var DB *sql.DB

var SEPERATOR = "*_*"

type UserMessage struct {
	Sender   string    `json:"sender"`
	Receiver string    `json:"receiver"`
	Created  time.Time `json:"created"`
	Image    string    `json:"image"`
	Content  string    `json:"content"`
	Type     string    `json:"type"`
	Context  string    `json:"context"`
}
