-- +migrate Up
INSERT INTO privateMessages(pmg_sender,pmg_receiver,pmg_content,pmg_image,creation_date)
VALUES
(2,1,"xd","http://localhost:8080/images/elmo.gif","2142-05-01"),
(2,1,"xdd","http://localhost:8080/images/feelsgoodman.png","2142-06-02"),
(3,1,"*_*","http://localhost:8080/images/theydobebirds.jpg","2142-07-03");

-- +migrate Down
DELETE FROM privateMessages
WHERE (pmg_id = 1)
OR (pmg_id = 2)
OR (pmg_id = 3);