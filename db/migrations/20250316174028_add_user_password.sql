-- +goose Up
-- +goose StatementBegin
alter table users
    add password char;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table users
    drop column password;
-- +goose StatementEnd
