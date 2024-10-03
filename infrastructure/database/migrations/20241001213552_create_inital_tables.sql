-- +goose Up
-- +goose StatementBegin
CREATE TABLE player (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    external_id VARCHAR(36) UNIQUE NOT NULL,
    summoner_id VARCHAR(50) UNIQUE NOT NULL,
    game_name VARCHAR(16) NOT NULL,
    tag_line VARCHAR(5) NOT NULL
);
CREATE TABLE ranked_info (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    external_id VARCHAR(36) UNIQUE NOT NULL,
    player_id INTEGER REFERENCES player(id),
    queue_type VARCHAR(20) NOT NULL,
    tier VARCHAR(15) NOT NULL,
    rank VARCHAR(15) NOT NULL,
    league_points INTEGER NOT NULL,
    wins INTEGER NOT NULL,
    losses INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE player;
DROP TABLE ranked_info;
-- +goose StatementEnd
