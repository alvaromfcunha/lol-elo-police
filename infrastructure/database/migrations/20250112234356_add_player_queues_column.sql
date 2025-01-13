-- +goose Up
-- +goose StatementBegin
ALTER TABLE player ADD COLUMN notify_queues VARCHAR(100) NOT NULL DEFAULT '420,450';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE player DROP COLUMN notify_queues;
-- +goose StatementEnd
