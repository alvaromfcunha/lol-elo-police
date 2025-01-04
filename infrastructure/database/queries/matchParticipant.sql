-- name: CreateMatchParticipant :one
INSERT INTO match_participant (
    external_id,
    player_id,
    matches_id,
    champion,
    kills,
    deaths,
    assists,
    is_win
) VALUES (
    @external_id,
    (SELECT id FROM player WHERE player.external_id = @player_external_id),
    (SELECT id FROM matches WHERE matches.external_id = @matches_external_id),
    @champion,
    @kills,
    @deaths,
    @assists,
    @is_win
) RETURNING *;