package infrastructure

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gbooksoft/backend/internal/domain/project/domain"
	"github.com/gbooksoft/backend/pkg/repository"
	"github.com/google/uuid"
	"strings"
)

type MySqlProjectRepository struct {
	*repository.BaseMysqlRepository
}

func NewMySqlProjectRepository(db *sql.DB) *MySqlProjectRepository {
	return &MySqlProjectRepository{repository.NewBaseMysqlRepository(db)}
}

func (r *MySqlProjectRepository) ExistsByUserIdAndTitle(tx *sql.Tx, userId domain.UserId, title string, ExceptProjectId *domain.ProjectId) (bool, error) {
	query := `
SELECT 1 
FROM projects 
WHERE user_id = ? AND title = ?
`
	if ExceptProjectId != nil {
		query = query + " AND id != ?"
	}

	var exists sql.NullBool

	var err error
	if ExceptProjectId != nil {
		err = tx.QueryRow(query, userId.String(), title, ExceptProjectId.String()).Scan(&exists)
	} else {
		err = tx.QueryRow(query, userId.String(), title).Scan(&exists)
	}

	if err == nil {
		return true, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	return false, err
}

func (r *MySqlProjectRepository) ExistsByUserIdAndProjectId(tx *sql.Tx, userId domain.UserId, projectId domain.ProjectId) (bool, error) {
	var exists sql.NullBool
	err := tx.QueryRow("SELECT 1 FROM projects WHERE user_id = ? AND id = ?", userId.String(), projectId.String()).Scan(&exists)

	if err == nil {
		return true, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	return false, err
}

func (r *MySqlProjectRepository) Add(tx *sql.Tx, project *domain.Project) error {
	query := "INSERT INTO projects (id, user_id, title, description, tags) VALUES (?, ?, ?, ?, ?)"
	desc := sql.NullString{
		String: project.Description,
		Valid:  len(project.Description) > 0,
	}

	var tagsJSON interface{}
	if len(project.Tags) == 0 {
		tagsJSON = nil
	} else {
		jsonData, err := json.Marshal(project.Tags)
		if err != nil {
			return err
		}
		tagsJSON = string(jsonData)
	}

	_, err := tx.Exec(query, project.Id.String(), project.UserId.String(), project.Title, desc, tagsJSON)

	if err != nil {
		return err
	}

	err = r.addGenresRelation(tx, project)

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

func (r *MySqlProjectRepository) FindMetadataByUserIdAndProjectId(tx *sql.Tx, userId domain.UserId, projectId domain.ProjectId) (*domain.Metadata, error) {
	query := `
SELECT pm.id, pm.author, pm.genre, pm.keywords, pm.language, pm.isbn
FROM project_metadata  pm
JOIN projects p ON p.id = pm.project_id
WHERE p.user_id = ? AND p.id = ?
`

	metadata := &domain.Metadata{}
	var id string
	var author sql.NullString
	var genre sql.NullString
	var keywords sql.NullString
	var language sql.NullString
	var isbn sql.NullString
	err := tx.QueryRow(query, userId.String(), projectId.String()).Scan(&id, &author, &genre, &keywords, &language, &isbn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	metadata.Id = domain.MetadataId(uuid.MustParse(id))
	metadata.ProjectId = projectId
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

func (r *MySqlProjectRepository) AddMetadata(tx *sql.Tx, metadata *domain.Metadata) error {
	query := `
INSERT INTO project_metadata (id, project_id, author, genre, keywords, language, isbn) VALUES (?, ?, ?, ?, ?, ?, ?)
`
	_, err := tx.Exec(query, metadata.Id.String(), metadata.ProjectId.String(), metadata.Author, metadata.Genre, metadata.Keywords, metadata.Language, metadata.Isbn)

	if err != nil {
		return err
	}

	return nil
}

func (r *MySqlProjectRepository) SaveMetadata(tx *sql.Tx, metadata *domain.Metadata) error {
	query := `
UPDATE project_metadata 
SET author=?, genre=?, keywords=?, language=?, isbn=?
WHERE id = ?
`
	author := sql.NullString{
		String: metadata.Author,
		Valid:  metadata.Author != "",
	}
	genre := sql.NullString{
		String: metadata.Genre,
		Valid:  metadata.Genre != "",
	}
	keywords := sql.NullString{
		String: metadata.Keywords,
		Valid:  metadata.Keywords != "",
	}
	language := sql.NullString{
		String: metadata.Language,
		Valid:  metadata.Language != "",
	}
	isbn := sql.NullString{
		String: metadata.Isbn,
		Valid:  metadata.Isbn != "",
	}

	_, err := tx.Exec(query, author, genre, keywords, language, isbn, metadata.Id.String())

	if err != nil {
		return err
	}

	return nil
}

func (r *MySqlProjectRepository) FindByUserIdAndProjectId(tx *sql.Tx, userId domain.UserId, projectId domain.ProjectId) (*domain.Project, error) {
	query := `
SELECT p.title, p.description, p.tags, GROUP_CONCAT(pg.genre_id, ',') AS genre_ids
FROM projects p 
JOIN project_genres pg ON pg.project_id = p.id
WHERE p.user_id = ? AND p.id = ?
`
	var title string
	var desc sql.NullString
	var tags sql.NullString
	var genres string

	err := r.Db.QueryRow(query, userId.String(), projectId.String()).Scan(&title, &desc, &tags, &genres)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	project := &domain.Project{
		Id:     projectId,
		UserId: userId,
		Title:  title,
	}

	if desc.Valid {
		project.Description = desc.String
	}

	if tags.Valid {
		err := json.Unmarshal([]byte(tags.String), &project.Tags)
		if err != nil {
			return nil, err
		}
	}

	for _, genre := range strings.Split(genres, ",") {
		if len(genre) == 0 {
			continue
		}
		project.GenresIds = append(project.GenresIds, domain.GenreId(uuid.MustParse(genre)))
	}

	return project, err
}

func (r *MySqlProjectRepository) Save(tx *sql.Tx, project *domain.Project) error {
	query := `
UPDATE projects
SET title = ?, description = ?, tags = ?
WHERE id = ?
`

	desc := sql.NullString{
		String: project.Description,
		Valid:  len(project.Description) > 0,
	}

	var tagsJSON interface{}
	if len(project.Tags) == 0 {
		tagsJSON = nil
	} else {
		jsonData, err := json.Marshal(project.Tags)
		if err != nil {
			return err
		}
		tagsJSON = string(jsonData)
	}

	_, err := tx.Exec(query, project.Title, desc, tagsJSON, project.Id.String())

	if err != nil {
		return err
	}

	err = r.updateGenresRelation(tx, project.Id.String(), project.GenresIds)
	if err != nil {
		return err
	}

	return nil
}

func (r *MySqlProjectRepository) updateGenresRelation(tx *sql.Tx, projectID string, genreIDs []domain.GenreId) error {
	_, err := tx.Exec("DELETE FROM project_genres WHERE project_id = ?", projectID)
	if err != nil {
		return err
	}

	if len(genreIDs) == 0 {
		return nil
	}

	query := "INSERT INTO project_genres (project_id, genre_id) VALUES "

	var values []interface{}
	for _, genreID := range genreIDs {
		query += "(?, ?),"
		values = append(values, projectID, genreID.String())
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
	if len(project.GenresIds) == 0 {
		return nil
	}

	query := "INSERT INTO project_genres (project_id, genre_id) VALUES (?, ?)"

	for _, genreID := range project.GenresIds {
		_, err := tx.Exec(query, project.Id.String(), genreID.String())
		if err != nil {
			return err
		}
	}

	return nil
}
