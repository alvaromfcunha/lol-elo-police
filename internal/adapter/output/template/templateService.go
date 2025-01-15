package template

import (
	"bytes"
	"text/template"

	"github.com/alvaromfcunha/lol-elo-police/internal/adapter/output/logger"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"
)

type TemplateService struct {
	templates *template.Template
}

func NewTemplateService(templates *template.Template) TemplateService {
	return TemplateService{templates}
}

func (s TemplateService) ExecuteNewMatchMessageTemplate(match entity.Match, participants []entity.MatchParticipant) (text string, err error) {
	logger.Debug(s, "Executing new match message template")

	var textBuf bytes.Buffer

	type newMatchTemplateData struct {
		Match        entity.Match
		Participants []entity.MatchParticipant
	}
	err = s.templates.ExecuteTemplate(&textBuf, "NewMatch", newMatchTemplateData{match, participants})
	if err != nil {
		logger.Error(s, "Cannot execute new match message template", err)
		return
	}

	text = textBuf.String()

	return
}
