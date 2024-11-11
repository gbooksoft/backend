package worker

import (
	"database/sql"
	"github.com/rs/zerolog/log"
	"time"
)

type ChangeEmailRequestWorker struct {
	db *sql.DB
}

func NewChangeEmailRequestWorker(db *sql.DB) *ChangeEmailRequestWorker {
	return &ChangeEmailRequestWorker{
		db: db,
	}
}

func (worker *ChangeEmailRequestWorker) Start() {
	for {

		query := `
DELETE email_change_requests
FROM email_change_requests
LEFT JOIN users ON email_change_requests.user_id = users.id
WHERE email_change_requests.created_at < NOW() - INTERVAL 15 MINUTE
   OR users.id IS NULL;
`
		result, err := worker.db.Exec(query)

		if err != nil {
			log.Error().Err(err).Msg("Failed to delete old email change requests")
		} else {
			n, _ := result.RowsAffected()

			if n > 0 {
				log.Info().Msgf("Email change requests deleted: %d", n)
			}
		}

		time.Sleep(15 * time.Minute) // Задержка для снижения нагрузки на базу данных
	}
}
