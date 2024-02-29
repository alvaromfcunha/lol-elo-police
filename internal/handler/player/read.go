package playerHandler

import (
	"github.com/alvaromfcunha/lol-elo-police/internal/model"
	"github.com/gofiber/fiber/v2"
)

func (h PlayerHandler) ReadAllPlayers(c *fiber.Ctx) error {
	var players []model.Player
	h.Db.Find(&players)

	return c.JSON(players)
}
