-- name: CreateRankedInfo :one
INSERT INTO ranked_info (
    external_id,
    player_id,
    queue_type,
    tier,
    rank,
    league_points,
    wins,
    losses
) VALUES (
    @external_id,
    (SELECT player.id FROM player WHERE player.external_id = @player_external_id),
    @queue_type,
    @tier,
    @rank,
    @league_points,
    @wins,
    @losses
) RETURNING *;

-- name: UpdateRankedInfo :exec
UPDATE
    ranked_info
SET
    tier = @tier,
    rank = @rank,
    league_points = @league_points,
    wins = @wins,
    losses = @losses
WHERE
    external_id = @external_id;

-- name: GetByPlayerExternalIdAndQueueType :one
SELECT
    sqlc.embed(ranked_info),
    sqlc.embed(player)
FROM
    ranked_info
INNER JOIN
    player
ON
    ranked_info.player_id = player.id
WHERE
    player_id = (SELECT player.id FROM player WHERE player.external_id = @player_external_id)
    AND queue_type = @queue_type;