-- +migrate Up
INSERT INTO posts (post_author, post_privacy, post_image, creation_date, post_content, post_title)
VALUES
(1, 'public', NULL, '2023-05-01', 'The Adventure Begins', 'First Adventure'),
(2, 'public', NULL, '2023-05-02', 'Future Gadgets Introduction', 'Second Adventure'),
(3, 'public', NULL, '2023-05-03', 'Time Machine', 'Third Adventure'),
(4, 'public', NULL, '2023-05-04', 'Invisible Cape', 'Fourth Adventure'),
(5, 'public', NULL, '2023-05-05', 'Anywhere Door', 'Fifth Adventure'),
(1, 'private', NULL, '2023-05-06', 'Mini-Dora', 'Sixth Adventure'),
(2, 'private', NULL, '2023-05-07', 'Pass Loop', 'Seventh Adventure'),
(3, 'private', NULL, '2023-05-08', 'Time Furoshiki', 'Eighth Adventure'),
(4, 'semi-private', NULL, '2023-05-09', 'Hopter', 'Ninth Adventure'),
(5, 'semi-private', NULL, '2023-05-10', 'Translation Gummy', 'Tenth Adventure');

-- +migrate Down
DELETE FROM posts
WHERE (post_title = 'First Adventure')
OR (post_title = 'Second Adventure')
OR (post_title = 'Third Adventure')
OR (post_title = 'Fourth Adventure')
OR (post_title = 'Fifth Adventure')
OR (post_title = 'Sixth Adventure')
OR (post_title = 'Seventh Adventure')
OR (post_title = 'Eighth Adventure')
OR (post_title = 'Ninth Adventure')
OR (post_title = 'Tenth Adventure');
