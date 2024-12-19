package email

import (
	"database/sql"
	"github.com/gbooksoft/backend/internal/domain/email/queue"
)

const TemplateForgotPassword queue.TemplateFile = "templates/notification/email/forgot_password.html"
const ForgotPassword = "GBook - сброс пароля"

type ForgotPasswordData struct {
	ResetPasswordLink string
}

func (f *Facade) ForgotPasswordTx(tx *sql.Tx, recipient string, data ForgotPasswordData) error {
	body, err := f.templateGenerator.ToString(TemplateForgotPassword, f.mergeWithDefault(ForgotPassword, data))

	if err != nil {
		return err
	}

	return f.queue.AddToQueueTx(tx, recipient, ForgotPassword, body)
}
