package data

import (
	"strconv"
	"strings"
	"time"

	"github.com/alvaromfcunha/lol-elo-police/internal/adapter/output/logger"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity/enum"
	"github.com/alvaromfcunha/lol-elo-police/internal/generated/database"
	"github.com/google/uuid"
)

func AssemblePlayer(player database.Player) entity.Player {
	var nqs []enum.QueueId

	if player.NotifyQueues != "" {
		nqss := strings.Split(player.NotifyQueues, ",")
		nqs = make([]enum.QueueId, len(nqss))
		for idx, q := range nqss {
			if qi, err := strconv.Atoi(q); err == nil {
				nqs[idx] = enum.QueueId(qi)
			} else {
				logger.Warn(AssemblePlayer, "Player: '" + player.ExternalID + "' NotifyQueues item: '" + q + "' not atoi'able")
			}
		}
	}

	return entity.Player{
		Id:           uuid.MustParse(player.ExternalID),
		SummonerId:   player.SummonerID,
		Puuid:        player.Puuid,
		GameName:     player.GameName,
		TagLine:      player.TagLine,
		NotifyQueues: nqs,
	}
}

func AssembleRankedInfo(rankedInfo database.RankedInfo, player database.Player) entity.RankedInfo {
	return entity.RankedInfo{
		Id:           uuid.MustParse(rankedInfo.ExternalID),
		Player:       AssemblePlayer(player),
		QueueType:    enum.QueueType(rankedInfo.QueueType),
		Tier:         rankedInfo.Tier,
		Rank:         rankedInfo.Rank,
		LeaguePoints: int(rankedInfo.LeaguePoints),
		Wins:         int(rankedInfo.Wins),
		Losses:       int(rankedInfo.Losses),
	}
}

func AssembleMatch(match database.Match) entity.Match {
	return entity.Match{
		Id:             uuid.MustParse(match.ExternalID),
		MatchId:        match.MatchID,
		QueueIdType:    int(match.QueueIDType),
		GameCreationAt: match.GameCreationAt,
		GameEndedAt:    match.GameEndedAt,
		GameDuration:   time.Duration(match.GameDuration),
	}
}
