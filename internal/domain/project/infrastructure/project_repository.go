package infrastructure

import (
	"database/sql"
	"errors"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/pkg/repository"
)

type MySqlProjectRepository struct {
	*repository.BaseMysqlRepository
}

func NewMySqlProjectRepository(db *sql.DB) *MySqlProjectRepository {
	return &MySqlProjectRepository{repository.NewBaseMysqlRepository(db)}
}

func (r *MySqlProjectRepository) ExistsByUserIdAndTitle(tx *sql.Tx, userId domain.UserId, title string) (bool, error) {
	var exists sql.NullBool
	err := tx.QueryRow("SELECT 1 FROM projects WHERE user_id = ? AND title = ?", userId.String(), title).Scan(&exists)

	if err == nil {
		return true, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	return false, err
}

func (r *MySqlProjectRepository) Add(tx *sql.Tx, project *domain.Project) error {
	_, err := tx.Exec("INSERT INTO projects (id, user_id, title) VALUES (?, ?, ?)", project.Id.String(), project.UserId.String(), project.Title)

	if err != nil {
		return err
	}

	return nil
}

func (r *MySqlProjectRepository) UserHasAccess(userId domain.UserId, projectId domain.ProjectId) (bool, error) {
	var exists sql.NullBool
	err := r.Db.QueryRow("SELECT 1 FROM projects WHERE user_id = ? AND id = ?", userId.String(), projectId.String()).Scan(&exists)

	if err == nil {
		return true, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	return false, err
}
