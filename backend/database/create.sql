CREATE TABLE user (
 id INTEGER NOT NULL PRIMARY KEY,
 privacy INTEGER NOT NULL,
 username VARCHAR(30) NOT NULL,
 passwrd VARCHAR(100) NOT NULL,
 email VARCHAR(30) NOT NULL,
 fname VARCHAR(30) NOT NULL,
 lname VARCHAR(30) NOT NULL,
 age INTEGER NOT NULL,
 avatar VARCHAR(100) NOT NULL,
 created_at DATETIME NOT NULL,
 about_me VARCHAR(1000) NOT NULL,
 last_session_id VARCHAR(100) NOT NULL,
 last_session_start DATETIME NOT NULL,
 last_session_end DATETIME NOT NULL,
 last_session_duration INTEGER NOT NULL
);

INSERT INTO user (id, privacy, username, passwrd, email, fname, lname, age, avatar, created_at, about_me, last_session_id, last_session_start, last_session_end, last_session_duration)
VALUES
    (1, 1, 'admin', '123', 'admin@hotmale.com', 'admin', 'admin', 99, 'https://imgpile.com/images/9NDeGC.jpg', DateTime('now', 'localtime'), 'I am the admin', '0', DateTime('now', 'localtime'), strftime('%s', DateTime('now', 'localtime')) + 100000, 100000),
    (2, 1, 'batman', '123', 'batman@gmail.com', 'John', 'Doe', 30, 'https://imgpile.com/images/9NDeGC.jpg', DateTime('now', 'localtime'), 'Lucas je vais te tracker', '1', DateTime('now', 'localtime'), strftime('%s', DateTime('now', 'localtime')) + 150000, 150000),
    (3, 1, 'wolverine', '123', 'logan@gmail.com', 'Jane', 'Doe', 28, 'https://imgpile.com/images/9NDeGC.jpg', DateTime('now', 'localtime'), 'Woof Woof Woof', '2', DateTime('now', 'localtime'), strftime('%s', DateTime('now', 'localtime')) + 120000, 120000),
    (4, 0, 'ironman', '123', 'ironman@gmail.com', 'Hidden', 'User', 35, 'https://imgpile.com/images/9NDeGC.jpg', DateTime('now', 'localtime'), 'efijnejifns', '3', DateTime('now', 'localtime'), strftime('%s', DateTime('now', 'localtime')) + 90000, 90000);

-- Table post


CREATE TABLE post (
 id INTEGER NOT NULL PRIMARY KEY,
 user_id INTEGER NOT NULL,
 privacy INTEGER NOT NULL,
 content VARCHAR(1000) NOT NULL,
 created_at DATETIME NOT NULL,
 FOREIGN KEY (user_id) REFERENCES user(id)
);

INSERT INTO post (id, user_id, privacy, content, created_at)

VALUES
    (1, 1, 0, 'Hello World!', DateTime('now', 'localtime')),
    (2, 1, 1, 'mon premier post :p', DateTime('now', 'localtime')),
    (3, 1, 2, ' admin', DateTime('now', 'localtime')),
    (4, 2, 1, 'En sah la jeunesse', DateTime('now', 'localtime')),
    (5, 2, 2, 'Wesh alors', DateTime('now', 'localtime')),
