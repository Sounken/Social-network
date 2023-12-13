-- +migrate Up
CREATE TABLE IF NOT EXISTS `comments`
(
    `comment_id` INTEGER,
    `comment_author` INTEGER NOT NULL,
    `post_id` INTEGER NOT NULL,
    `comment_content` TEXT,
    `comment_image` TEXT,
    `creation_date` TIMESTAMP NOT NULL,
    PRIMARY KEY(`comment_id` AUTOINCREMENT)
);

-- +migrate Down
DROP TABLE `comments`;