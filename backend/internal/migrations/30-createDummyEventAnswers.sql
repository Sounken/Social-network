-- +migrate Up
INSERT INTO eventOptionChoices (event_id,uuid,choice)
VALUES
(1, 1, "going"),
(1, 2, "not going"),
(1, 5, "i hate events"),
(2, 2, "not going"),
(2, 5, "i hate events");

-- +migrate Down
DELETE FROM eventOptionChoices
WHERE (event_id = 1)
OR (event_id = 2);