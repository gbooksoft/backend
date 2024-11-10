package email

import (
	"database/sql"
	"github.com/gbooksoft/backend/internal/domain/email/queue"
)

const TemplateEmailChangeRequestInformation queue.TemplateFile = "templates/notification/email/email_change_request_information.html"
const SubjectEmailChangeRequestInformation = "Вы сменили почту"

type EmailChangeRequestInformationData struct {
	FirstName string
	LastName  string
	NewEmail  string
}

func (f *Facade) EmailChangeRequestInformationTx(tx *sql.Tx, recipient string, data EmailChangeRequestInformationData) error {
	body, err := f.templateGenerator.ToString(TemplateEmailChangeRequestInformation, f.mergeWithDefault(SubjectEmailChangeRequestInformation, data))

	if err != nil {
		return err
	}
	return f.queue.AddToQueueTx(tx, recipient, SubjectEmailChangeRequestInformation, body)
}
