package email

import (
	"database/sql"
	"github.com/gbooksoft/backend/internal/domain/email/queue"
)

const TemplatePasswordChanged queue.TemplateFile = "templates/notification/email/password_changed.html"
const SubjectPasswordChanged = "GBook | Вы сменили пароли"

type PasswordChangedData struct {
	FirstName string
	LastName  string
}

func (f *Facade) PasswordChangedTx(tx *sql.Tx, recipient string, data PasswordChangedData) error {
	body, err := f.templateGenerator.ToString(TemplatePasswordChanged, f.mergeWithDefault(SubjectPasswordChanged, data))

	if err != nil {
		return err
	}

	return f.queue.AddToQueueTx(tx, recipient, SubjectPasswordChanged, body)
}
