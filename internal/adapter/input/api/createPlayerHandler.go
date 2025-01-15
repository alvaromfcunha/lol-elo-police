package api

import (
	"github.com/alvaromfcunha/lol-elo-police/internal/adapter/output/logger"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity/enum"
	usecase "github.com/alvaromfcunha/lol-elo-police/internal/domain/useCase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CreatePlayerHandler struct {
	UseCase usecase.CreatePlayerUseCase
	Validator *validator.Validate
}

func NewCreatePlayerHandler(useCase usecase.CreatePlayerUseCase) CreatePlayerHandler {
	return CreatePlayerHandler{
		UseCase: useCase,
		Validator: validator.New(),
	}
}

type CreatePlayerRequest struct {
	GameName     string         `json:"name" validate:"required"`
	TagLine      string         `json:"tag" validate:"required"`
	NotifyQueues []enum.QueueId `json:"notifyQueues" validate:"required"`
}

func (h CreatePlayerHandler) Handle(ctx *fiber.Ctx) error {
	logger.Info(h, "Handling create player request")

	request := new(CreatePlayerRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		logger.Error(h, "Error parsing create player request", err)
		return err
	}

	if err := h.Validator.Struct(request); err != nil {
		logger.Error(h, "Error validating create player request", err)
		return err
	}

	input := usecase.CreatePlayerInput{
		GameName: request.GameName,
		TagLine:  request.TagLine,
		NotifyQueues: request.NotifyQueues,
	}

	output, err := h.UseCase.Execute(input)
	if err != nil {
		logger.Error(h, "Error on create player use case", err)
		return err
	}

	ctx.Status(201)
	ctx.JSON(output)

	return nil
}
