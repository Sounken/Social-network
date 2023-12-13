-- +migrate Up
INSERT INTO notifications(notif_content,creation_date,uuid,sender_id,notif_type,notif_context)
VALUES
("A user would like to join your group!","2023-05-01",1,3,"group_join_request","The Bullies"),
("A user would like to join your group!","2023-05-01",1,4,"group_join_request","The Bullies");

-- +migrate Down
DELETE FROM notifications
WHERE (notif_id = 1)
OR (notif_id = 2);