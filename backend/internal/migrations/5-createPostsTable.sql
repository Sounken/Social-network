-- +migrate Up
CREATE TABLE IF NOT EXISTS `posts`
(
    `post_id` INTEGER,
    `post_author` INTEGER NOT NULL,
    `post_privacy` TEXT NOT NULL,
    `post_image` TEXT,
    `creation_date` TIMESTAMP NOT NULL,
    `post_content` TEXT NOT NULL,
    `post_title` TEXT NOT NULL,
    PRIMARY KEY(`post_id` AUTOINCREMENT)
);

-- +migrate Down
DROP TABLE `posts`;