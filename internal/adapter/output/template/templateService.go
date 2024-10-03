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

func (s TemplateService) ExecuteQueueUpdateMessageTemplate(data service.QueueUpdateData) (text string, err error) {
	var textBuf bytes.Buffer

	err = s.templates.ExecuteTemplate(&textBuf, "QueueUpdate", data)
	if err != nil {
		return
	}

	text = textBuf.String()

	return
}

func (s TemplateService) ExecuteQueueNewEntryMessageTemplate(data service.QueueNewEntryData) (text string, err error) {
	var textBuf bytes.Buffer

	err = s.templates.ExecuteTemplate(&textBuf, "QueueNewEntry", data)
	if err != nil {
		return
	}

	text = textBuf.String()

	return
}
