-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_login() RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT' AND NEW.success) THEN
        UPDATE users SET last_login=now() WHERE email = NEW.email;
    end if;
    RETURN NEW;
END;
$$LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS user_login_update on login_attempts;
CREATE TRIGGER user_login_update
    AFTER INSERT ON login_attempts
    FOR EACH ROW EXECUTE FUNCTION update_login();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS user_login_update on login_attempts;
DROP FUNCTION IF EXISTS update_login;
-- +goose StatementEnd
