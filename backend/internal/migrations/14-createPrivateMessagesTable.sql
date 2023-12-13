-- +migrate Up
CREATE TABLE IF NOT EXISTS `privateMessages`
(
    `pmg_id` INTEGER,
    `pmg_sender` INTEGER NOT NULL,
    `pmg_receiver` INTEGER NOT NULL,
    `pmg_content` TEXT NOT NULL,
    `pmg_image` TEXT,
    `creation_date` TIMESTAMP NOT NULL,
    PRIMARY KEY(`pmg_id` AUTOINCREMENT)
);

-- +migrate Down
DROP TABLE `privateMessages`;