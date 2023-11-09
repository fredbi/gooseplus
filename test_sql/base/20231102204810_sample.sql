-- +goose Up
-- +goose StatementBegin
CREATE TABLE items(
  id integer NOT NULL PRIMARY KEY,
  name text NOT NULL,
  warehouse_location text NOT NULL,
  weight double NOT NULL,
  
  description text,

  last_updated timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items;
-- +goose StatementEnd
