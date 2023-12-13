-- +migrate Up
CREATE TABLE IF NOT EXISTS `eventOptionChoices`
(
    `event_id` INTEGER NOT NULL,
    `uuid` INTEGER NOT NULL,
    `choice` TEXT NOT NULL
);

-- +migrate Down
DROP TABLE `eventOptionChoices`;