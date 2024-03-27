-- +goose Up
-- +goose StatementBegin
CREATE TABLE urls (short VARCHAR(8) PRIMARY KEY, original TEXT UNIQUE);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE urls;
-- +goose StatementEnd
