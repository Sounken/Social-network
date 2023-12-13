-- +migrate Up
CREATE TABLE IF NOT EXISTS `followers`
(
    `uuid` INTEGER NOT NULL,
    `follower_id` INTEGER NOT NULL
);

-- +migrate Down
DROP TABLE `followers`;