package email

import (
	"github.com/gbooksoft/backend/internal/domain/email/queue"
	"os"
)

type baseData struct {
	Title        string
	AppName      string
	WebAppDomain string
}

type Facade struct {
	queue             *queue.EmailQueue
	templateGenerator *queue.TemplateGenerator
}

func NewFacade(queue *queue.EmailQueue, templateGenerator *queue.TemplateGenerator) *Facade {
	return &Facade{
		queue:             queue,
		templateGenerator: templateGenerator,
	}
}

func (f *Facade) mergeWithDefault(title string, data any) any {
	return map[string]any{
		"Base": f.getBaseData(title),
		"Data": data,
	}
}

func (f *Facade) getBaseData(pageTitle string) baseData {
	return baseData{
		Title:        pageTitle,
		AppName:      os.Getenv("APP_NANE"),
		WebAppDomain: os.Getenv("WEB_FRONT_FULL_DOMAIN"),
	}
}
