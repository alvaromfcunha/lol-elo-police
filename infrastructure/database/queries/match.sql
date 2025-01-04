-- name: CreateMatch :one
INSERT INTO matches (
    external_id,
    match_id,
    queue_id_type,
    game_creation_at,
    game_ended_at,
    game_duration
) VALUES (
    @external_id,
    @match_id,
    @queue_id_type,
    @game_creation_at,
    @game_ended_at,
    @game_duration
) RETURNING *;

-- name: GetLastestMatchesByPlayerExternalId :many
SELECT
    sqlc.embed(matches)
FROM
    matches
INNER JOIN
    match_participant
ON
    matches.id = match_participant.matches_id
INNER JOIN
    player
ON
    match_participant.player_id = player.id
WHERE
    player.external_id = @player_external_id
ORDER BY
    matches.game_ended_at DESC
LIMIT
    1;

-- name: GetMatchesByMatchId :one
SELECT
    sqlc.embed(matches)
FROM
    matches
WHERE
    match_id = @match_id