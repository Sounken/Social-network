-- +migrate Up
INSERT INTO groupChat(group_id,gc_sender,gc_content,gc_image,creation_date)
VALUES
(1,1,":D","http://localhost:8080/images/feelsgoodman.png","2142-05-01"),
(1,2,">:DD","http://localhost:8080/images/feelsgoodman.png","2142-05-02"),
(1,5,":DDD","http://localhost:8080/images/feelsgoodman.png","2142-05-03");

-- +migrate Down
DELETE FROM groupChat
WHERE (creation_date = "2142-05-01")
OR (creation_date = "2142-05-02")
OR (creation_date = "2142-05-03");