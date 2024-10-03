package cron

import (
	"fmt"

	usecase "github.com/alvaromfcunha/lol-elo-police/internal/domain/useCase"
)

type PatrolRoutineHandler struct {
	UseCase usecase.PolicePatrol
}

func (h PatrolRoutineHandler) Handle() {
	err := h.UseCase.Execute()
	// Better error logging.
	if err != nil {
		fmt.Println("[CRON][PolicePatrol] - Error:", err.Error())
	}
}
