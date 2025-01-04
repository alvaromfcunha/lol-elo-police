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

func (s TemplateService) ExecuteNewMatchMessageTemplate(data service.NewMatchData) (text string, err error) {
	var textBuf bytes.Buffer

	err = s.templates.ExecuteTemplate(&textBuf, "NewMatch", data)
	if err != nil {
		return
	}

	text = textBuf.String()

	return
}
