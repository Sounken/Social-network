-- +migrate Up
CREATE TABLE IF NOT EXISTS `groupChat`
(
    `gc_id` INTEGER,
    `group_id` INTEGER NOT NULL,
    `gc_sender` INTEGER NOT NULL,
    `gc_content` TEXT NOT NULL,
    `gc_image` TEXT,
    `creation_date` TIMESTAMP NOT NULL,
    PRIMARY KEY(`gc_id` AUTOINCREMENT)
);

-- +migrate Down
DROP TABLE `groupChat`;