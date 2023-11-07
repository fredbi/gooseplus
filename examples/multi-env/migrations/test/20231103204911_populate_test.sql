-- +goose Up
-- +goose StatementBegin
INSERT INTO example(id, value) VALUES('two', 'Two for testing');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM example WHERE id = 'two';
-- +goose StatementEnd
