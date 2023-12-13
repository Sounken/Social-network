-- +migrate Up
INSERT INTO events (group_id,event_author,event_title,event_content,creation_date,event_date,options)
VALUES
(1,1,"first event","better fucking work",'2023-05-01', "2450-01-05", "going*_*not going*_*i hate events"),
(1,1,"second event","better fucking work",'2023-05-01', "2450-01-05", "going*_*not going*_*i hate events");

-- +migrate Down
DELETE FROM events
WHERE (event_id = 1)
OR (event_id = 2);