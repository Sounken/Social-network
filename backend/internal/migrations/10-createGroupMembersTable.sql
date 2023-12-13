-- +migrate Up
CREATE TABLE IF NOT EXISTS `groupMembers`
(
    `group_id` INTEGER NOT NULL,
    `uuid` INTEGER NOT NULL,
    `role` TEXT NOT NULL
);

-- +migrate Down
DROP TABLE `groupMembers`;