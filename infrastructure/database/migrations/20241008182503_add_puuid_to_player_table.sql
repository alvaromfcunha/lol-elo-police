-- +goose Up
-- +goose StatementBegin
ALTER TABLE player ADD COLUMN puuid VARCHAR(78) NOT NULL DEFAULT '-';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE player DROP COLUMN puuid;
-- +goose StatementEnd
