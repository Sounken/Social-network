-- +migrate Up
INSERT INTO comments (comment_author, post_id, comment_content, comment_image, creation_date)
VALUES
(1, 1, 'The Adventure of Nerds', NULL, '2023-05-02'),
(2, 1, 'The Adventure of Noobs', NULL, '2023-05-03'),
(3, 1, 'The Adventure of Weebs', NULL, '2023-05-04'),
(4, 1, 'The Adventure of Cracks', NULL, '2023-05-05'),
(5, 1, 'The Adventure of Gods', NULL, '2023-05-06');

-- +migrate Down
DELETE FROM comments
WHERE (comment_author = 1)
OR (comment_author = 2)
OR (comment_author = 3)
OR (comment_author = 4)
OR (comment_author = 5);