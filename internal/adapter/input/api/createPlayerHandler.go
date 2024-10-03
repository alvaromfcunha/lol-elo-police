package api

import (
	usecase "github.com/alvaromfcunha/lol-elo-police/internal/domain/useCase"
	"github.com/gofiber/fiber/v2"
)

type CreatePlayerHandler struct {
	UseCase usecase.CreatePlayer
}

type CreatePlayerRequest struct {
	GameName string `json:"name"`
	TagLine  string `json:"tag"`
}

func (h CreatePlayerHandler) Handle(ctx *fiber.Ctx) error {
	request := new(CreatePlayerRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		return err
	}

	input := usecase.CreatePlayerInput{
		GameName: request.GameName,
		TagLine:  request.TagLine,
	}

	output, err := h.UseCase.Execute(input)
	// Better API error handling.
	if err != nil {
		return err
	}

	ctx.Status(201)
	ctx.JSON(output)

	return nil
}
