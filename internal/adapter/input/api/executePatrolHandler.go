package api

import (
	"github.com/alvaromfcunha/lol-elo-police/internal/adapter/output/logger"
	usecase "github.com/alvaromfcunha/lol-elo-police/internal/domain/useCase"
	"github.com/gofiber/fiber/v2"
)

type ExecutePatrolHandler struct {
	UseCase usecase.PolicePatrolUseCase
}

func NewExecutePatrolHandler(useCase usecase.PolicePatrolUseCase) ExecutePatrolHandler {
	return ExecutePatrolHandler{
		UseCase: useCase,
	}
}

func (h ExecutePatrolHandler) Handle(ctx *fiber.Ctx) error {
	logger.Info(h, "Handling execute patrol request")

	err := h.UseCase.Execute()
	if err != nil {
		logger.Error(h, "Error on execute patrol use case", err)
		return err
	}

	ctx.Status(200)

	return nil
}
