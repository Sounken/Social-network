-- +migrate Up
INSERT INTO groups(group_creator, group_name, group_description, creation_date)
VALUES 
(1, 'The Bullies', 'nerds need not apply', '2023-04-19'),
(3, 'The Bullied', 'do not come to school tomorrow', '2023-04-20'),
(5, 'The Mediator', 'peace and love', '2022-12-22');

-- +migrate Down
DELETE FROM groups
WHERE (group_name = 'The Bullies')
OR (group_name = 'The Bullied')
OR (group_name = 'The Mediator');