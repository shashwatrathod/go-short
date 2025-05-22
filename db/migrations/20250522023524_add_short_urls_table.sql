-- +goose Up
-- +goose StatementBegin
CREATE TABLE short_urls (
    short_url varchar(8) PRIMARY KEY,
    original_url VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ language 'plpgsql';
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER update_short_urls_updated_at
BEFORE UPDATE ON short_urls
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_short_urls_updated_at ON short_urls;
-- +goose StatementEnd
-- +goose StatementBegin
DROP FUNCTION IF EXISTS update_updated_at_column();
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS short_urls;
-- +goose StatementEnd
