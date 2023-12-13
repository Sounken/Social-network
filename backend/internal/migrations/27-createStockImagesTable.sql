-- +migrate Up
CREATE TABLE IF NOT EXISTS `stockImages`
(
    `uid` INTEGER,
    `image_path` TEXT,
    PRIMARY KEY(`uid` AUTOINCREMENT)
);
-- +migrate Up
INSERT INTO stockImages(image_path)
VALUES
("http://localhost:8080/images/feelsgoodman.png"),
("http://localhost:8080/images/letthefirepurgeusall.png"),
("http://localhost:8080/images/pepehands.png"),
("http://localhost:8080/images/philosoraptor.png"),
("http://localhost:8080/images/elmo.gif"),
("http://localhost:8080/images/theydobebirds.jpg");

-- +migrate Down
DROP TABLE `stockImages`;