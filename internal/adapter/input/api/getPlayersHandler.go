package api

import (
	usecase "github.com/alvaromfcunha/lol-elo-police/internal/domain/useCase"
	"github.com/gofiber/fiber/v2"
)

type GetPlayersHandler struct {
	UseCase usecase.GetPlayers
}

func (h GetPlayersHandler) Handle(ctx *fiber.Ctx) error {
	output, err := h.UseCase.Execute()
	// Better API error handling.
	if err != nil {
		return err
	}

	ctx.Status(200)
	ctx.JSON(output)

	return nil
}
