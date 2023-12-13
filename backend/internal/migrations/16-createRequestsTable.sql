-- +migrate Up
CREATE TABLE IF NOT EXISTS `requests`
(
    `req_id` INTEGER,
    `req_sender` INTEGER NOT NULL,
    `req_receiver` INTEGER NOT NULL,
    `req_type` TEXT NOT NULL,
    `creation_date` TIMESTAMP NOT NULL,
    PRIMARY KEY(`req_id` AUTOINCREMENT)
);

-- +migrate Down
DROP TABLE `requests`;