package queue

import (
	"database/sql"
	"time"
)

type TemplateData interface{}

type Status string

const (
	StatusPending    Status = "pending"
	StatusSent       Status = "sent"
	StatusFailed     Status = "failed"
	StatusProcessing Status = "processing"
)

type Email struct {
	Id         int64
	Recipient  string
	Subject    string
	Body       string
	RetryCount int
}

type EmailQueue struct {
	db *sql.DB
}

func NewEmailQueue(db *sql.DB) *EmailQueue {
	return &EmailQueue{
		db: db,
	}
}

func (s *EmailQueue) AddToQueue(recipient, subject, body string) error {
	query := `INSERT INTO email_queue (recipient, subject, body, status) VALUES (?, ?, ?, ?)`
	_, err := s.db.Exec(query, recipient, subject, body, StatusPending)
	if err != nil {
		return err
	}

	return nil
}

func (s *EmailQueue) AddToQueueTx(tx *sql.Tx, recipient, subject, body string) error {
	query := `INSERT INTO email_queue (recipient, subject, body, status) VALUES (?, ?, ?, ?)`
	_, err := tx.Exec(query, recipient, subject, body, StatusPending)
	if err != nil {
		return err
	}

	return nil
}

func (s *EmailQueue) GetPendingEmailsFromQueue(limit int) ([]Email, error) {
	query := `
SELECT id, recipient, subject, body, retry_count
FROM email_queue 
WHERE status = ? OR (status = ? AND next_attempt_at <= NOW())
LIMIT ? FOR UPDATE SKIP LOCKED
`
	rows, err := s.db.Query(query, StatusPending, StatusFailed, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []Email
	for rows.Next() {
		var email Email
		if err := rows.Scan(&email.Id, &email.Recipient, &email.Subject, &email.Body, &email.RetryCount); err != nil {
			return nil, err
		}
		// Обновляем статус на "Processing" для текущего письма
		_, err = s.db.Exec("UPDATE email_queue SET status = ? WHERE id = ?", StatusProcessing, email.Id)
		if err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return emails, nil
}

func (s *EmailQueue) HandleFailed(email Email, errorMessage string) error {
	maxRetries := 3

	if email.RetryCount >= maxRetries {
		_, err := s.db.Exec("UPDATE email_queue SET status = ?, error_message = ?, next_attempt_at = NULL WHERE id = ?", StatusFailed, errorMessage, email.Id)
		if err != nil {
			return err
		}

		return nil
	}

	retryDelay := time.Duration(email.RetryCount+1) * time.Minute // Увеличиваем задержку с каждой попыткой

	nextAttempt := time.Now().Add(retryDelay)

	query := `
UPDATE email_queue 
SET status = ?, error_message = ?, retry_count = retry_count + 1, next_attempt_at = ? 
WHERE id = ?
`
	_, err := s.db.Exec(query, StatusFailed, errorMessage, nextAttempt, email.Id)

	if err != nil {
		return err
	}

	return nil
}

func (s *EmailQueue) UpdateEmailStatus(queueId int64, status Status, errorMessage string) error {

	_, err := s.db.Exec("UPDATE email_queue SET status = ?, error_message = ? WHERE id = ?", string(status), errorMessage, queueId)
	if err != nil {
		return err
	}

	return nil
}
