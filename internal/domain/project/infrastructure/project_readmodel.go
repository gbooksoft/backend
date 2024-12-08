package infrastructure

import (
	"database/sql"
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
	if err := row.Scan(&project.Id, &project.Title, &project.Description); err != nil {
		return readmodel.Project{}, err
	}

	return project, nil
}
