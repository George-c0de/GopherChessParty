-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN new_id UUID;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
UPDATE users
SET new_id = uuid_generate_v5('00000000-0000-0000-0000-000000000000', id::text);
ALTER TABLE users DROP COLUMN id;
ALTER TABLE users RENAME COLUMN new_id TO id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN old_id INT;
ALTER TABLE users ADD COLUMN id_new SERIAL;
ALTER TABLE users DROP COLUMN id;
ALTER TABLE users RENAME COLUMN id_new TO id;
-- +goose StatementEnd
