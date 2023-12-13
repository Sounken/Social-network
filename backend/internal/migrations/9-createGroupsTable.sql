-- +migrate Up
CREATE TABLE IF NOT EXISTS `groups`
(
    `group_id` INTEGER,
    `group_creator` INTEGER NOT NULL,
    `group_name` TEXT NOT NULL,
    `group_description` TEXT NOT NULL,
    `creation_date` TIMESTAMP NOT NULL,
    PRIMARY KEY(`group_id` AUTOINCREMENT)
);

-- +migrate Down
DROP TABLE `groups`;