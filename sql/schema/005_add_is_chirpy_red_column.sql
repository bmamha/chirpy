-- +goose Up
ALTER TABLE users
ADD COLUMN is_chirpy_red BOOLEAN NOT NULL DEFAULT FALSE; 

-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
