-- +migrate Up
-- 1. Enums
-- +migrate StatementBegin
DO $$ BEGIN
    CREATE TYPE workspacetype AS ENUM ('PERSONAL', 'TEAM');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;
-- +migrate StatementEnd

-- +migrate StatementBegin
DO $$ BEGIN
    CREATE TYPE workspacerole AS ENUM ('ADMIN', 'MEMBER', 'VIEWER');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;
-- +migrate StatementEnd

-- +migrate StatementBegin
DO $$ BEGIN
    CREATE TYPE teamrole AS ENUM ('ADMIN', 'MEMBER');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;
-- +migrate StatementEnd

-- +migrate StatementBegin
DO $$ BEGIN
    CREATE TYPE projectrole AS ENUM ('ADMIN', 'MEMBER', 'VIEWER');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;
-- +migrate StatementEnd

-- +migrate StatementBegin
DO $$ BEGIN
    CREATE TYPE taskstatus AS ENUM ('TODO', 'IN_PROGRESS', 'DONE', 'BACKLOG');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;
-- +migrate StatementEnd

-- +migrate StatementBegin
DO $$ BEGIN
    CREATE TYPE taskpriority AS ENUM ('P0', 'P1', 'P2', 'P3');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;
-- +migrate StatementEnd

-- 2. Users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    full_name VARCHAR(255),
    hashed_password VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS ix_users_email ON users(email);
