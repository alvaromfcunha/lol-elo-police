package controller

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/alvaromfcunha/lol-elo-police/internal/adapter/input/cron"
	"github.com/alvaromfcunha/lol-elo-police/internal/adapter/output/data"
	"github.com/alvaromfcunha/lol-elo-police/internal/adapter/output/http"
	templateService "github.com/alvaromfcunha/lol-elo-police/internal/adapter/output/template"
	"github.com/alvaromfcunha/lol-elo-police/internal/adapter/output/whatsapp"
	usecase "github.com/alvaromfcunha/lol-elo-police/internal/domain/useCase"
	"go.mau.fi/whatsmeow"
)

func PatrolRoutineController(
	db *sql.DB,
	templates *template.Template,
	whatsmeowClient *whatsmeow.Client,
) func() {
	return func() {
		riotApiKey, ok := os.LookupEnv("RIOT_API_KEY")
		if !ok {
			fmt.Println("missing RIOT_API_KEY in configuration")
			return
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()

		playerData := data.NewPlayerData(
			ctx,
			db,
		)
		rankedInfoData := data.NewRankedInfoData(
			ctx,
			db,
		)

		lolApi := http.NewLolApi(riotApiKey)

		whatsappService := whatsapp.NewWhatsappService(whatsmeowClient)

		templateService := templateService.NewTemplateService(
			templates,
		)

		useCase := usecase.PolicePatrol{
			PlayerRepository:     playerData,
			RankedInfoRepository: rankedInfoData,
			LolService:           lolApi,
			WhatsappService:      whatsappService,
			TemplateService:      templateService,
		}

		handler := cron.PatrolRoutineHandler{
			UseCase: useCase,
		}

		handler.Handle()
	}
}
