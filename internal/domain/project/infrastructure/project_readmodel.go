package infrastructure

import (
	"database/sql"
	"errors"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
)

type MySqlProjectReadModel struct {
	db *sql.DB
}

func NewMySqlProjectReadModel(db *sql.DB) *MySqlProjectReadModel {
	return &MySqlProjectReadModel{db: db}
}

func (r *MySqlProjectReadModel) GetAll(userId string) ([]readmodel.Project, error) {
	query := `
SELECT p.id, p.title, p.description 
FROM projects p
WHERE p.user_id = ?
ORDER BY p.created_at DESC
`
	rows, err := r.db.Query(query, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var projects []readmodel.Project

	for rows.Next() {
		var project readmodel.Project

		var pDesc sql.NullString
		if err := rows.Scan(&project.Id, &project.Title, &pDesc); err != nil {
			return nil, err
		}

		if pDesc.Valid {
			project.Description = pDesc.String
		}

		projects = append(projects, project)
	}

	return projects, nil
}

func (r *MySqlProjectReadModel) GetRecent(userId string, limit int) ([]readmodel.Project, error) {
	query := `
SELECT p.id, p.title, p.description, MAX(pe.updated_at) AS last_updated
FROM projects p
JOIN project_menu_items pmi ON pmi.project_id = p.id
JOIN project_editors pe ON pe.project_menu_id = pmi.id
WHERE p.user_id = ?
GROUP BY p.id, p.title, p.description
ORDER BY last_updated DESC
LIMIT ?
`
	rows, err := r.db.Query(query, userId, limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var projects []readmodel.Project

	for rows.Next() {
		var project readmodel.Project

		var pDesc sql.NullString
		var pUpdateAt sql.NullString
		if err := rows.Scan(&project.Id, &project.Title, &pDesc, &pUpdateAt); err != nil {
			return nil, err
		}

		if pDesc.Valid {
			project.Description = pDesc.String
		}

		projects = append(projects, project)
	}

	return projects, nil
}

func (r *MySqlProjectReadModel) GetById(userId string, projectId string) (readmodel.Project, error) {
	query := `
SELECT p.id, p.title, p.description
FROM projects p
WHERE p.user_id = ? AND p.id = ?
`
	row := r.db.QueryRow(query, userId, projectId)

	var project readmodel.Project
	var desc sql.NullString
	if err := row.Scan(&project.Id, &project.Title, &desc); err != nil {
		return readmodel.Project{}, err
	}

	if desc.Valid {
		project.Description = desc.String
	}

	return project, nil
}

func (r *MySqlProjectReadModel) GetMetadata(userId string, projectId string) (*readmodel.Metadata, error) {
	query := `
SELECT pm.author, pm.genre, pm.keywords, pm.language, pm.isbn
FROM project_metadata  pm
JOIN projects p ON p.id = pm.project_id
WHERE p.user_id = ? AND p.id = ?
`
	row := r.db.QueryRow(query, userId, projectId)

	metadata := &readmodel.Metadata{}
	var author sql.NullString
	var genre sql.NullString
	var keywords sql.NullString
	var language sql.NullString
	var isbn sql.NullString

	if err := row.Scan(&author, &genre, &keywords, &language, &isbn); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	if author.Valid {
		metadata.Author = author.String
	}
	if genre.Valid {
		metadata.Genre = genre.String
	}
	if keywords.Valid {
		metadata.Keywords = keywords.String
	}
	if language.Valid {
		metadata.Language = language.String
	}
	if isbn.Valid {
		metadata.Isbn = isbn.String
	}

	return metadata, nil
}
