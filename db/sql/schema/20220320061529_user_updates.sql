-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_request_user(user_name text) RETURNS void AS $$
BEGIN
    PERFORM set_config('request.user', user_name, false);
END ;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION audit_user_delete() RETURNS TRIGGER AS $$
DECLARE
    deleted_by VARCHAR ;
BEGIN
    BEGIN
        deleted_by := COALESCE (current_setting('request.user', true), user );
    END;
    IF (tg_op = 'DELETE') THEN
        INSERT INTO user_deletions (email, deleted_by, deleted_at)
        VALUES (OLD.email, deleted_by, now());
    END IF;

    RETURN OLD;
end;
$$LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS user_delete_audit ON users;
CREATE TRIGGER user_delete_audit
    AFTER DELETE ON users
    FOR EACH ROW EXECUTE FUNCTION audit_user_delete();

CREATE OR REPLACE FUNCTION audit_user_update() RETURNS TRIGGER AS $$
DECLARE
    changed_by VARCHAR ;
BEGIN
    BEGIN
        changed_by := COALESCE (current_setting('request.user', true), user );
    END;

    IF (tg_op = 'UPDATE') THEN
        IF (
            old.passhash != coalesce(new.passhash, old.passhash) OR
            old.is_active != coalesce(new.is_active, old.is_active) OR
            old.is_admin != coalesce(new.is_admin, old.is_admin)
            ) THEN
            NEW.updated_at = now();
            INSERT INTO user_changes (email, updated_by, pwd, is_active, is_admin, created_at)
            VALUES (OLD.email, changed_by::varchar, OLD.passhash != NEW.passhash, OLD.is_active != NEW.is_active, OLD.is_admin != NEW.is_admin, now());
            RETURN NEW;
        END IF;
    END IF;
    RETURN new;
END;
$$LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS user_change_audit ON users;
CREATE TRIGGER user_change_audit
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION audit_user_update();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS set_request_user;

DROP TRIGGER IF EXISTS user_delete_audit ON users;
DROP FUNCTION IF EXISTS audit_user_delete;

DROP TRIGGER IF EXISTS user_change_audit ON users;
DROP FUNCTION IF EXISTS audit_user_update;
-- +goose StatementEnd
