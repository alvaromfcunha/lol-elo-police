-- +goose Up
-- +goose StatementBegin
ALTER TABLE matches ADD COLUMN NEW_game_creation_at DATETIME NOT NULL;
UPDATE matches SET NEW_game_creation_at = game_creation_at;
ALTER TABLE matches DROP COLUMN game_creation_at;
ALTER TABLE matches RENAME COLUMN NEW_game_creation_at TO game_creation_at;

ALTER TABLE matches ADD COLUMN NEW_game_ended_at DATETIME NOT NULL;
UPDATE matches SET NEW_game_ended_at = game_ended_at;
ALTER TABLE matches DROP COLUMN game_ended_at;
ALTER TABLE matches RENAME COLUMN NEW_game_ended_at TO game_ended_at;

ALTER TABLE ranked_info ADD COLUMN created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE match_participant ADD COLUMN new_ranked_info_id INTEGER REFERENCES ranked_info(id);
ALTER TABLE match_participant ADD COLUMN prev_ranked_info_id INTEGER REFERENCES ranked_info(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE matches ADD COLUMN NEW_game_creation_at TIMESTAMP NOT NULL;
UPDATE matches SET NEW_game_creation_at = game_creation_at;
ALTER TABLE matches DROP COLUMN game_creation_at;
ALTER TABLE matches RENAME COLUMN NEW_game_creation_at TO game_creation_at;

ALTER TABLE matches ADD COLUMN NEW_game_ended_at TIMESTAMP NOT NULL;
UPDATE matches SET NEW_game_ended_at = game_ended_at;
ALTER TABLE matches DROP COLUMN game_ended_at;
ALTER TABLE matches RENAME COLUMN NEW_game_ended_at TO game_ended_at;

ALTER TABLE ranked_info DROP COLUMN created_at;
ALTER TABLE match_participant DROP COLUMN new_ranked_info_id;
ALTER TABLE match_participant DROP COLUMN prev_ranked_info_id;
-- +goose StatementEnd
