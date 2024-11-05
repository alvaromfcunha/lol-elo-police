package template

import (
	"bytes"
	"text/template"

	"github.com/alvaromfcunha/lol-elo-police/internal/domain/service"
)

type TemplateService struct {
	templates *template.Template
}

func NewTemplateService(templates *template.Template) TemplateService {
	return TemplateService{templates}
}

func (s TemplateService) ExecuteNewRankedMatchMessageTemplate(data service.NewRankedMatchData) (text string, err error) {
	var textBuf bytes.Buffer

	err = s.templates.ExecuteTemplate(&textBuf, "NewRankedMatch", data)
	if err != nil {
		return
	}

	text = textBuf.String()

	return
}

func (s TemplateService) ExecuteNewUnrankedMatchMessageTemplate(data service.NewUnrankedMatchData) (text string, err error) {
	var textBuf bytes.Buffer

	err = s.templates.ExecuteTemplate(&textBuf, "NewUnrankedMatch", data)
	if err != nil {
		return
	}

	text = textBuf.String()

	return
}
