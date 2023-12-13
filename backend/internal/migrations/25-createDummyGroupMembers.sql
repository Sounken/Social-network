-- +migrate Up
INSERT OR IGNORE INTO groupMembers(group_id, uuid, role)
VALUES
(1, 1, 'creator'),
(1, 2, 'member'),
(1, 5, 'member'),
(2, 3, 'creator'),
(2, 4, 'member'),
(3, 5, 'creator'),
(3, 4, 'member');

-- +migrate Down
DELETE FROM groupMembers
WHERE (group_id = 1)
OR (group_id = 2)
OR (group_id = 3);