-- +migrate Up
CREATE TABLE IF NOT EXISTS `groupPosts`
(
    `gpost_id` INTEGER,
    `gpost_author` INTEGER NOT NULL,
    `group_id` INTEGER NOT NULL,
    `gpost_image` TEXT,
    `creation_date` TIMESTAMP NOT NULL,
    `gpost_content` TEXT NOT NULL,
    `gpost_title` TEXT NOT NULL,
    PRIMARY KEY(`gpost_id` AUTOINCREMENT)
);

-- +migrate Down
DROP TABLE `groupPosts`;