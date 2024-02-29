package playerHandler

import (
	"net/http"

	"github.com/alvaromfcunha/lol-elo-police/internal/model"
	"github.com/alvaromfcunha/lol-elo-police/pkg/lol"
	"github.com/gofiber/fiber/v2"
)

type CreatePlayerRequest struct {
	GameName string `json:"name"`
	TagLine  string `json:"tag"`
}

func (h PlayerHandler) CreatePlayer(c *fiber.Ctx) error {
	request := new(CreatePlayerRequest)
	err := c.BodyParser(request)
	if err != nil {
		return err
	}

	acc, err := h.LolApi.GetAccountByRiotId(lol.RiotId{
		GameName: request.GameName,
		TagLine:  request.TagLine,
	})
	if err != nil {
		return &fiber.Error{
			Message: "error getting account",
			Code:    http.StatusBadRequest,
		}
	}

	smm, err := h.LolApi.GetSummonerByPuuid(acc.Puuid)
	if err != nil {
		return &fiber.Error{
			Message: "error getting summoner",
			Code:    http.StatusBadRequest,
		}
	}

	leagues, err := h.LolApi.GetLeaguesBySummonerId(smm.Id)
	if err != nil {
		return &fiber.Error{
			Message: "error getting leagues",
			Code:    http.StatusBadRequest,
		}
	}

	var solo *lol.LeagueEntry
	for _, league := range leagues {
		if league.QueueType == lol.QueueType(lol.Solo) {
			solo = &league
		}
	}

	if solo == nil {
		return &fiber.Error{
			Message: "no solo ranked league for player specified",
			Code:    http.StatusBadRequest,
		}
	}

	player := &model.Player{
		SummonerId:   smm.Id,
		GameName:     acc.GameName,
		TagLine:      acc.TagLine,
		Tier:         solo.Tier,
		Rank:         solo.Rank,
		LeaguePoints: solo.LeaguePoints,
		Wins:         solo.Wins,
		Losses:       solo.Losses,
	}

	h.Db.Create(player)

	c.Status(201)
	return c.JSON(player)
}
