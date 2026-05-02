-- +goose Up
-- +goose StatementBegin

ALTER TABLE users
ADD CONSTRAINT users_chat_id_telegram_unique UNIQUE (chat_id_telegram);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE users
DROP CONSTRAINT users_chat_id_telegram_unique;

-- +goose StatementEnd
