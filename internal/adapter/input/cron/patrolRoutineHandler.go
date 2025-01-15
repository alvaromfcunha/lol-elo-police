package cron

import (
	"github.com/alvaromfcunha/lol-elo-police/internal/adapter/output/logger"
	usecase "github.com/alvaromfcunha/lol-elo-police/internal/domain/useCase"
)

type PatrolRoutineHandler struct {
	UseCase usecase.PolicePatrolUseCase
}

func NewPatrolRoutineHandler(useCase usecase.PolicePatrolUseCase) PatrolRoutineHandler {
	return PatrolRoutineHandler{
		UseCase: useCase,
	}
}

func (h PatrolRoutineHandler) Handle() {
	logger.Info(h, "Handling patrol routine")

	err := h.UseCase.Execute()
	if err != nil {
		logger.Error(h, "Error on patrol routine use case", err)
	}
}
