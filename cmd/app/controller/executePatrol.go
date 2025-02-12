package controller

import (
	"database/sql"
	"errors"
	"os"
	"text/template"

	"github.com/alvaromfcunha/lol-elo-police/internal/adapter/input/api"
	"github.com/alvaromfcunha/lol-elo-police/internal/adapter/output/data"
	"github.com/alvaromfcunha/lol-elo-police/internal/adapter/output/http"
	templateService "github.com/alvaromfcunha/lol-elo-police/internal/adapter/output/template"
	"github.com/alvaromfcunha/lol-elo-police/internal/adapter/output/whatsapp"
	usecase "github.com/alvaromfcunha/lol-elo-police/internal/domain/useCase"
	"github.com/gofiber/fiber/v2"
	"go.mau.fi/whatsmeow"
)

func ExecutePatrolController(
	db *sql.DB,
	riotHttpClient *http.RateLimitedClient,
	templates *template.Template,
	whatsmeowClient *whatsmeow.Client,
) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		riotApiKey, ok := os.LookupEnv("RIOT_API_KEY")
		if !ok {
			return errors.New("missing RIOT_API_KEY in configuration")
		}

		context := ctx.Context()

		playerData := data.NewPlayerData(
			context,
			db,
		)
		rankedInfoData := data.NewRankedInfoData(
			context,
			db,
		)
		matchData := data.NewMatchData(
			context,
			db,
		)
		matchParticipantData := data.NewMatchParticipantData(
			context,
			db,
		)

		lolApi := http.NewLolApi(
			riotHttpClient,
			riotApiKey,
		)

		whatsappService := whatsapp.NewWhatsappService(whatsmeowClient)

		templateService := templateService.NewTemplateService(
			templates,
		)

		useCase := usecase.NewPolicePatrolUseCase(
			playerData,
			rankedInfoData,
			matchData,
			matchParticipantData,
			lolApi,
			whatsappService,
			templateService,
		)

		handler := api.NewExecutePatrolHandler(useCase)

		return handler.Handle(ctx)
	}
}
