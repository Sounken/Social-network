package socket

import (
	"backend/internal/data"
	"backend/internal/helper"
	"time"
)

func SendNotificationToAUser(sender, receiver int, content, context, notifType string) error {

	senderName, err := helper.GetUsername(sender)
	if err != nil {
		return err
	}
	receiverName, err := helper.GetUsername(receiver)
	if err != nil {
		return err
	}

	newMsg := data.UserMessage{
		Receiver: receiverName,
		Sender:   senderName,
		Created:  time.Now(),
		Context:  context,
		Content:  content,
		Type:     notifType,
	}
	FindAndSend(newMsg)

	return err
}
