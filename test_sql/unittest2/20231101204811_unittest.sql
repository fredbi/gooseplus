-- +goose Up
-- +goose StatementBegin
CREATE TABLE unittest_pre(
  id text PRIMARY KEY
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS unittest_pre;
-- +goose StatementEnd
