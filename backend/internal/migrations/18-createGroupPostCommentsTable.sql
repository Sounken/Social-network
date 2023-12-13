-- +migrate Up
CREATE TABLE IF NOT EXISTS `groupPostComments`
(
    `gpcomment_id` INTEGER,
    `gpost_id` INTEGER NOT NULL,
    `group_id` INTEGER NOT NULL,
    `gpcomment_author` INTEGER NOT NULL,
    `gpcomment_content` TEXT,
    `gpcomment_image` TEXT,
    `creation_date` TIMESTAMP NOT NULL,
    PRIMARY KEY(`gpcomment_id` AUTOINCREMENT)
);

-- +migrate Down
DROP TABLE `groupPostComments`;