package controller

import (
	"database/sql"
	"os"

	"github.com/alvaromfcunha/lol-elo-police/internal/adapter/input/api"
	"github.com/alvaromfcunha/lol-elo-police/internal/adapter/output/data"
	"github.com/alvaromfcunha/lol-elo-police/internal/adapter/output/http"
	usecase "github.com/alvaromfcunha/lol-elo-police/internal/domain/useCase"
	"github.com/gofiber/fiber/v2"
)

func CreatePlayerController(db *sql.DB) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		riotApiKey, ok := os.LookupEnv("RIOT_API_KEY")
		if !ok {
			panic("missing RIOT_API_KEY in configuration")
		}

		playerData := data.NewPlayerData(
			ctx.Context(),
			db,
		)
		rankedInfoData := data.NewRankedInfoData(
			ctx.Context(),
			db,
		)
		lolApi := http.NewLolApi(riotApiKey)

		useCase := usecase.CreatePlayer{
			PlayerRepository:     playerData,
			RankedInfoRepository: rankedInfoData,
			LolService:           lolApi,
		}

		handler := api.CreatePlayerHandler{
			UseCase: useCase,
		}

		return handler.Handle(ctx)
	}
}
