-- +goose Up
-- +goose StatementBegin
-- Users Table.
CREATE TABLE IF NOT EXISTS users (
    email       varchar(100)    PRIMARY KEY,
    created_by  varchar(100),
    passhash    varchar(60)     NOT NULL,
    is_active   boolean         NOT NULL DEFAULT true,
    is_admin    boolean         NOT NULL DEFAULT false,
    created_at  timestamp       NOT NULL DEFAULT (now()),
    updated_at  timestamp       NOT NULL DEFAULT (now()),
    last_login  timestamp
    );

CREATE INDEX users_admin_idx ON users(is_admin);
CREATE INDEX users_created_at_idx ON users(email, created_at);
CREATE INDEX users_updated_at_idx ON users(email, updated_at);
CREATE INDEX users_last_login_idx ON users(email, last_login);

-- Login Attempts Table.
CREATE TABLE IF NOT EXISTS login_attempts (
    id                  SERIAL primary key,
    email               varchar(100) NOT NULL,
    success             boolean NOT NULL,
    login_time          timestamp DEFAULT (now())
    );

CREATE INDEX IF NOT EXISTS login_attempts_email_idx ON login_attempts(email);
CREATE INDEX IF NOT EXISTS login_attempts_success_idx ON login_attempts(success);
CREATE INDEX IF NOT EXISTS login_attempts_time_idx ON login_attempts(login_time);
CREATE INDEX IF NOT EXISTS login_attempts_all_idx ON login_attempts(email, success, login_time);

-- User changes table.
CREATE TABLE IF NOT EXISTS user_changes (
    id                  SERIAL primary key,
    email               varchar(100) NOT NULL,
    updated_by          varchar(100),
    pwd                 boolean DEFAULT false,
    is_active           boolean DEFAULT false,
    is_admin            boolean DEFAULT false,
    created_at          timestamp DEFAULT (now())
    );

CREATE INDEX IF NOT EXISTS user_changes_email_idx ON user_changes(email);

-- User deletions table.
CREATE TABLE IF NOT EXISTS user_deletions (
    id                  SERIAL primary key,
    email               varchar(100) NOT NULL,
    deleted_by          varchar(100),
    deleted_at          timestamp DEFAULT (now())
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users, login_attempts, user_changes, user_deletions;
-- +goose StatementEnd
