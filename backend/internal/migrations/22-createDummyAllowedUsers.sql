-- +migrate Up
INSERT INTO allowedUsers (post_id, uuid)
VALUES
(9, 1),
(9, 2),
(9, 3),
(10, 4),
(10, 5),
(10, 3);

-- +migrate Down
DELETE FROM allowedUsers
WHERE (post_id = 9)
OR (post_id = 10);
