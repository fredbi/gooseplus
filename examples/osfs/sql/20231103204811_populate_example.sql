-- +goose Up
-- +goose StatementBegin
INSERT INTO example(id, value) VALUES('one', 'One');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM example WHERE id = 'one';
-- +goose StatementEnd
