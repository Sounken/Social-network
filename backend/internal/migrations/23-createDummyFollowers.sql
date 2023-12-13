-- +migrate Up
INSERT INTO followers(uuid, follower_id)
VALUES
(1, 4),
(1, 5),
(2, 1),
(2, 3),
(3, 1),
(3, 2);

-- +migrate Down
DELETE FROM followers
WHERE (uuid = 1)
OR (uuid = 2)
OR (uuid = 3);