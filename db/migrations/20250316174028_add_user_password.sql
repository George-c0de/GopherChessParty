-- +goose Up
-- +goose StatementBegin
alter table users
    add password VARCHAR(255) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table users
    drop column password;
-- +goose StatementEnd
