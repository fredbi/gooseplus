-- +goose Up
-- +goose StatementBegin
CREATE TABLE unittest(
  id text PRIMARY KEY
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS unittest;
-- +goose StatementEnd
