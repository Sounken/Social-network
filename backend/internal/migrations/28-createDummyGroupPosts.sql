-- +migrate Up
INSERT INTO groupPosts (gpost_author,group_id, creation_date, gpost_content, gpost_title)
VALUES
(1, 1,  '2023-05-01', 'The Adventure Begins', 'First Adventure'),
(2, 1,  '2023-05-02', 'Future Gadgets Introduction', 'Second Adventure'),
(3, 1,  '2023-05-03', 'Time Machine', 'Third Adventure'),
(4, 1,  '2023-05-04', 'Invisible Cape', 'Fourth Adventure'),
(5, 1,  '2023-05-05', 'Anywhere Door', 'Fifth Adventure'),
(1, 1,  '2023-05-06', 'Mini-Dora', 'Sixth Adventure'),
(2, 1,  '2023-05-07', 'Pass Loop', 'Seventh Adventure'),
(3, 1,  '2023-05-08', 'Time Furoshiki', 'Eighth Adventure'),
(4, 1,  '2023-05-09', 'Hopter', 'Ninth Adventure'),
(5, 1,  '2023-05-10', 'Translation Gummy', 'Tenth Adventure');

-- +migrate Down
DELETE FROM groupPosts
WHERE (gpost_title = 'First Adventure')
OR (gpost_title = 'Second Adventure')
OR (gpost_title = 'Third Adventure')
OR (gpost_title = 'Fourth Adventure')
OR (gpost_title = 'Fifth Adventure')
OR (gpost_title = 'Sixth Adventure')
OR (gpost_title = 'Seventh Adventure')
OR (gpost_title = 'Eighth Adventure')
OR (gpost_title = 'Ninth Adventure')
OR (gpost_title = 'Tenth Adventure');
