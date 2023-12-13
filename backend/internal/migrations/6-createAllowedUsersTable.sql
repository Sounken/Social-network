-- +migrate Up
CREATE TABLE IF NOT EXISTS `allowedUsers`
(
    `post_id` INTEGER NOT NULL,
    `uuid` INTEGER NOT NULL
);

-- +migrate Down
DROP TABLE `allowedUsers`;