-- +migrate Up
CREATE TABLE IF NOT EXISTS `users`
(
	`uuid`	INTEGER,
	`password`	TEXT NOT NULL,
	`email`	TEXT NOT NULL,
	`first_name`	TEXT NOT NULL,
	`last_name`	TEXT NOT NULL,
	`DOB`	TIMESTAMP NOT NULL,
	`avatar`	TEXT,
	`username`	TEXT,
	`bio`	TEXT,
	`session_token`	TEXT,
	`privacy` TEXT NOT NULL,
	PRIMARY KEY(`uuid` AUTOINCREMENT)
);

-- +migrate Down
DROP TABLE `users`;