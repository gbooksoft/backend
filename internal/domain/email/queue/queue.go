package queue

import (
	"database/sql"
)

type TemplateData interface{}

type Status string

const (
	StatusPending    Status = "pending"
	StatusSent       Status = "sent"
	StatusFailed     Status = "failed"
	StatusProcessing Status = "processing"
)

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
