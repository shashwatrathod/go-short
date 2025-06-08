-- +goose Up
-- +goose StatementBegin
ALTER TABLE short_urls
RENAME TO url_aliases;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE url_aliases
RENAME COLUMN short_url
TO alias;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE url_aliases
RENAME COLUMN alias
TO short_url;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE url_aliases
RENAME TO short_urls;
-- +goose StatementEnd