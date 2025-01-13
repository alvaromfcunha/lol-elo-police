-- name: CreateMatchParticipant :one
INSERT INTO match_participant (
    external_id,
    player_id,
    matches_id,
    new_ranked_info_id,
    prev_ranked_info_id,
    champion,
    kills,
    deaths,
    assists,
    is_win
) VALUES (
    @external_id,
    (SELECT id FROM player WHERE player.external_id = @player_external_id),
    (SELECT id FROM matches WHERE matches.external_id = @matches_external_id),
    (SELECT id FROM ranked_info WHERE ranked_info.external_id = @new_ranked_info_external_id),
    (SELECT id FROM ranked_info WHERE ranked_info.external_id = @prev_ranked_info_external_id),
    @champion,
    @kills,
    @deaths,
    @assists,
    @is_win
) RETURNING *;
