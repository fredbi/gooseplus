-- +goose Up
-- +goose StatementBegin
CREATE TABLE example(
  id text PRIMARY KEY,
  value text NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS example;
-- +goose StatementEnd
