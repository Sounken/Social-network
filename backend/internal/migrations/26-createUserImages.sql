-- +migrate Up
CREATE TABLE IF NOT EXISTS `userImages`
(
    `uuid` INTEGER NOT NULL,
    `image_path` TEXT NOT NULL,
    PRIMARY KEY(`uuid`,`image_path`)
);

-- +migrate Down
DROP TABLE `userImages`;