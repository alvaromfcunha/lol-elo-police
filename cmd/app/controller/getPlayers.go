package controller

import (
	"database/sql"

	"github.com/alvaromfcunha/lol-elo-police/internal/adapter/input/api"
	"github.com/alvaromfcunha/lol-elo-police/internal/adapter/output/data"
	usecase "github.com/alvaromfcunha/lol-elo-police/internal/domain/useCase"
	"github.com/gofiber/fiber/v2"
)

func GetPlayersController(db *sql.DB) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		playerData := data.NewPlayerData(
			ctx.Context(),
			db,
		)
		rankedInfoData := data.NewRankedInfoData(
			ctx.Context(),
			db,
		)

		useCase := usecase.GetPlayers{
			PlayerRepository:     playerData,
			RankedInfoRepository: rankedInfoData,
		}

		handler := api.GetPlayersHandler{
			UseCase: useCase,
		}

		return handler.Handle(ctx)
	}
}
