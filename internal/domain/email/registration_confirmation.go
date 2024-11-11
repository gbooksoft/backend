package email

import (
	"database/sql"
	"github.com/gbooksoft/backend/internal/domain/email/queue"
)

const TemplateRegistrationConfirmation queue.TemplateFile = "templates/notification/email/registration_confirmation.html"
const SubjectRegistrationConfirmation = "Добро пожаловать в наш сервис!"

type RegistrationConfirmationData struct {
	ConfirmationLink string
	FirstName        string
	LastName         string
}

func (f *Facade) RegistrationConfirmationTx(tx *sql.Tx, recipient string, data RegistrationConfirmationData) error {
	body, err := f.templateGenerator.ToString(TemplateRegistrationConfirmation, f.mergeWithDefault(SubjectRegistrationConfirmation, data))

	if err != nil {
		return err
	}

	return f.queue.AddToQueueTx(tx, recipient, SubjectRegistrationConfirmation, body)
}
