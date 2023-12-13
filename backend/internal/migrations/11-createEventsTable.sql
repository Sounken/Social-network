-- +migrate Up
CREATE TABLE IF NOT EXISTS `events`
(
    `event_id` INTEGER,
    `group_id` INTEGER NOT NULL,
    `event_author` INTEGER NOT NULL,
    `event_title` TEXT NOT NULL,
    `event_content` TEXT NOT NULL,
    `creation_date` TIMESTAMP NOT NULL,
    `event_date` TIMESTAMP NOT NULL,
    `options` TEXT NOT NULL,
    PRIMARY KEY(`event_id` AUTOINCREMENT)
);

-- +migrate Down
DROP TABLE `events`;