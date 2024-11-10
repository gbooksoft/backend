package email

import (
	"database/sql"
	"github.com/gbooksoft/backend/internal/domain/email/queue"
)

const TemplateEmailChangeRequestConfirmation queue.TemplateFile = "templates/notification/email/email_change_request_confirmation.html"
const SubjectEmailChangeRequestConfirmation = "Подтверждение смены почты"

type EmailChangeRequestConfirmationData struct {
	FirstName        string
	LastName         string
	ConfirmationLink string
}

func (f *Facade) EmailChangeRequestConfirmationTx(tx *sql.Tx, recipient string, data EmailChangeRequestConfirmationData) error {
	body, err := f.templateGenerator.ToString(TemplateEmailChangeRequestConfirmation, f.mergeWithDefault(SubjectEmailChangeRequestConfirmation, data))

	if err != nil {
		return err
	}
	return f.queue.AddToQueueTx(tx, recipient, SubjectEmailChangeRequestConfirmation, body)
}
