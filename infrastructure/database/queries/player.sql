-- name: CreatePlayer :one
INSERT INTO player (
    external_id,
    summoner_id,
    game_name,
    tag_line
) VALUES (
    @external_id,
    @summoner_id,
    @game_name,
    @tag_line
) RETURNING *;

-- name: GetPlayers :many
SELECT
    sqlc.embed(player)
FROM
    player;

-- name: UpdatePlayer :exec
UPDATE 
    player
SET
    solo_queue_id = (SELECT id FROM ranked_info WHERE ranked_info.external_id = @solo_queue_id),
    flex_queue_id = (SELECT id FROM ranked_info WHERE ranked_info.external_id = @flex_queue_id)
WHERE
    player.external_id = @external_id;