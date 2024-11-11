package worker

import (
	"github.com/gbooksoft/backend/internal/domain/email/queue"
	"github.com/rs/zerolog/log"
	"github.com/wneessen/go-mail"
	"os"
	"time"
)

type EmailWorker struct {
	queue      *queue.EmailQueue
	mailClient *mail.Client
}

func NewEmailWorker(queue *queue.EmailQueue, mailClient *mail.Client) *EmailWorker {
	return &EmailWorker{
		queue:      queue,
		mailClient: mailClient,
	}
}

func (worker *EmailWorker) Start() {
	for {
		worker.processEmails()
		time.Sleep(2 * time.Second) // Задержка для снижения нагрузки на базу данных
	}
}

func (worker *EmailWorker) processEmails() {
	emails, err := worker.queue.GetPendingEmailsFromQueue(5)
	if err != nil {
		log.Error().Err(err).Msgf("Ошибка при получении письма из очереди")
		return
	}

	if len(emails) == 0 {
		// Нет писем для отправки
		return
	}

	for _, email := range emails {
		err = worker.sendEmail(email)
		if err != nil {
			log.Error().Err(err).Msgf("Ошибка при отправке письма ID %v", email.Id)
			err = worker.queue.UpdateEmailStatus(email.Id, queue.StatusFailed, err.Error())

			if err != nil {
				log.Error().Err(err).Msgf("Ошибка при обновлении статуса письма ID %d: %v", email.Id, err)
			}
			return
		}

		log.Info().Msgf("Письмо ID %d успешно отправлено", email.Id)
		err = worker.queue.UpdateEmailStatus(email.Id, queue.StatusSent, "")
		if err != nil {
			log.Error().Err(err).Msgf("Ошибка при обновлении статуса письма ID %d: %v", email.Id, err)
		}
	}
}

func (worker *EmailWorker) sendEmail(email queue.Email) error {
	msg := mail.NewMsg()
	if err := msg.From(os.Getenv("EMAIL_FROM")); err != nil {
		log.Error().Err(err).Msgf("failed to set FROM address")
		return err
	}
	if err := msg.To(email.Recipient); err != nil {
		log.Error().Err(err).Msgf("failed to set TO address")
		return err
	}
	msg.Subject(email.Subject)
	msg.SetBodyString(mail.TypeTextHTML, email.Body)

	if err := worker.mailClient.DialAndSend(msg); err != nil {
		log.Error().Err(err).Msgf("Ошибка при отправке письма")
		return err
	}

	return nil
}
