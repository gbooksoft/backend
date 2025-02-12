package infrastructure

import (
	"database/sql"
	"errors"
	"github.com/gbooksoft/backend/internal/domain/project2/domain"
	"github.com/gbooksoft/backend/pkg/repository"
)

type MySqlProjectRepository struct {
	*repository.BaseMysqlRepository
}

func NewMySqlProjectRepository(db *sql.DB) *MySqlProjectRepository {
	return &MySqlProjectRepository{repository.NewBaseMysqlRepository(db)}
}

func (r *MySqlProjectRepository) ExistsByUserIdAndTitle(tx *sql.Tx, authorId domain.AuthorId, title string) (bool, error) {
	query := `
SELECT 1 
FROM projects2 
WHERE author_id = ? AND title = ?
`
	var exists sql.NullBool

	err := tx.QueryRow(query, authorId.String(), title).Scan(&exists)

	if err == nil {
		return true, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	return false, err
}

func (r *MySqlProjectRepository) Add(tx *sql.Tx, project *domain.Project) error {
	query := "INSERT INTO projects2 (id, author_id, title, description) VALUES (?, ?, ?, ?)"
	desc := sql.NullString{
		String: project.Description,
		Valid:  len(project.Description) > 0,
	}

	_, err := tx.Exec(query, project.Id.String(), project.AuthorId.String(), project.Title, desc)

	if err != nil {
		return err
	}

	err = r.addGenresRelation(tx, project)

	if err != nil {
		return err
	}

	return nil
}

func (r *MySqlProjectRepository) Save(tx *sql.Tx, project *domain.Project) error {
	query := `
UPDATE projects2
SET title = ?, description = ?
WHERE id = ?
`

	desc := sql.NullString{
		String: project.Description,
		Valid:  len(project.Description) > 0,
	}

	_, err := tx.Exec(query, project.Title, desc, project.Id.String())

	if err != nil {
		return err
	}

	err = r.updateGenresRelation(tx, project.Id.String(), project.Genres)
	if err != nil {
		return err
	}

	return nil
}

func (r *MySqlProjectRepository) updateGenresRelation(tx *sql.Tx, projectId string, genreIDs []domain.GenreId) error {
	_, err := tx.Exec("DELETE FROM project2_genres WHERE project_id = ?", projectId)
	if err != nil {
		return err
	}

	if len(genreIDs) == 0 {
		return nil
	}

	query := "INSERT INTO project2_genres (project_id, genre_id) VALUES "

	var values []interface{}
	for _, genreID := range genreIDs {
		query += "(?, ?),"
		values = append(values, projectId, genreID.String())
	}

	// Убираем последнюю запятую
	query = query[:len(query)-1]

	_, err = tx.Exec(query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (r *MySqlProjectRepository) addGenresRelation(tx *sql.Tx, project *domain.Project) error {
	if len(project.Genres) == 0 {
		return nil
	}

	query := "INSERT INTO project2_genres (project_id, genre_id) VALUES (?, ?)"

	for _, genreId := range project.Genres {
		_, err := tx.Exec(query, project.Id.String(), genreId.String())
		if err != nil {
			return err
		}
	}

	return nil
}
