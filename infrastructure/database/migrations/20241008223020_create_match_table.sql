-- +goose Up
-- +goose StatementBegin
CREATE TABLE matches (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    external_id VARCHAR(36) UNIQUE NOT NULL,
    match_id VARCHAR(14) UNIQUE NOT NULL,
    queue_id_type INTEGER NOT NULL,
    game_creation_at VARCHAR(33) NOT NULL,
    game_ended_at VARCHAR(33) NOT NULL,
    game_duration INTEGER NOT NULL
);
CREATE TABLE match_participant (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    external_id VARCHAR(36) UNIQUE NOT NULL,
    player_id INTEGER NOT NULL REFERENCES player(id),
    champion VARCHAR(25) NOT NULL,
    kills INTEGER NOT NULL,
    deaths INTEGER NOT NULL,
    assists INTEGER NOT NULL,
    is_win BOOLEAN NOT NULL,
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE matches;
DROP TABLE match_participant;
-- +goose StatementEnd
