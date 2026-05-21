-- +goose Up
-- +goose StatementBegin

ALTER TABLE users ADD COLUMN profile_document_url TEXT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE users DROP COLUMN IF EXISTS profile_document_url;

-- +goose StatementEnd
