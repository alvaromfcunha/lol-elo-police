package api

import (
	"github.com/alvaromfcunha/lol-elo-police/internal/adapter/output/logger"
	usecase "github.com/alvaromfcunha/lol-elo-police/internal/domain/useCase"
	"github.com/gofiber/fiber/v2"
)

type GetPlayersHandler struct {
	UseCase usecase.GetPlayersUseCase
}

func NewGetPlayersHandler(useCase usecase.GetPlayersUseCase) GetPlayersHandler {
	return GetPlayersHandler{
		UseCase: useCase,
	}
}

func (h GetPlayersHandler) Handle(ctx *fiber.Ctx) error {
	logger.Info(h, "Handling get players request")

	output, err := h.UseCase.Execute()
	if err != nil {
		logger.Error(h, "Error on get players use case", err)
		return err
	}

	ctx.Status(200)
	ctx.JSON(output)

	return nil
}
