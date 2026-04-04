-- +migrate Down
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS taskpriority;
DROP TYPE IF EXISTS taskstatus;
DROP TYPE IF EXISTS projectrole;
DROP TYPE IF EXISTS teamrole;
DROP TYPE IF EXISTS workspacerole;
DROP TYPE IF EXISTS workspacetype;
