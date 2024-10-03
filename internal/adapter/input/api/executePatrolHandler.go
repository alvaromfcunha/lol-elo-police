package api

import (
	usecase "github.com/alvaromfcunha/lol-elo-police/internal/domain/useCase"
	"github.com/gofiber/fiber/v2"
)

type ExecutePatrolHandler struct {
	UseCase usecase.PolicePatrol
}

func (h ExecutePatrolHandler) Handle(ctx *fiber.Ctx) error {
	err := h.UseCase.Execute()
	// Better API error handling.
	if err != nil {
		return err
	}

	ctx.Status(200)

	return nil
}
