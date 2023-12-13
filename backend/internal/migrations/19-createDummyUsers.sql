-- +migrate Up
INSERT INTO users (password, email, first_name, last_name, DOB, avatar, username, bio, session_token, privacy)
VALUES
('123', 'nobita@futaba.jp', 'Nobita', 'Nobi', '1973-08-07', NULL, 'Nobita73', 'Loves to sleep and daydream', '1', 'public'),
('123', 'doraemon@futaba.jp', 'Doraemon', '', '2112-09-03', NULL, 'Doraemon', 'A cat-like robot from the 22nd century', '2', 'public'),
('123', 'shizuka@futaba.jp', 'Shizuka', 'Minamoto', '1973-05-08', NULL, 'Shizuka73', 'Loves to take bath and is kind to everyone', '3', 'public'),
('123', 'gian@futaba.jp', 'Takeshi', 'Goda', '1973-06-15', NULL, 'Gian73', 'Strong and quick-tempered', '4', 'public'),
('123', 'suneo@futaba.jp', 'Suneo', 'Honekawa', '1973-02-29', NULL, 'Suneo73', 'Show-off and is fond of boasting', '5', 'private');

-- +migrate Down
DELETE FROM users
WHERE (first_name = 'Nobita')
OR (first_name = 'Doraemon')
OR (first_name = 'Shizuka')
OR (first_name = 'Takeshi')
OR (first_name = 'Suneo');