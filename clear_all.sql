PRAGMA foreign_keys = OFF;

BEGIN;

DELETE FROM users;
DELETE FROM posts;
DELETE FROM messages;

COMMIT;

VACUUM;
