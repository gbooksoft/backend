package infrastructure

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gbooksoft/backend/internal/domain/project2/readmodel"
	"gorm.io/gorm"
)

type MySqlProjectReadModel struct {
	db *gorm.DB
}

func NewMySqlProjectReadModel(db *gorm.DB) *MySqlProjectReadModel {
	return &MySqlProjectReadModel{
		db: db,
	}
}

func (r *MySqlProjectReadModel) GetAllLoose(authorId string) ([]readmodel.Project, error) {
	rows, err := r.db.
		Table("projects2 AS p").
		Select("p.id AS project_id", "p.title", "p.description", "pm.author", "pm.keywords", "pm.language", "pm.isbn").
		Joins("join project_metadata pm on p.id = pm.project_id").
		Where("p.author_id = ?", authorId).
		Rows()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var projectIds []string
	projects := make(map[string]*readmodel.Project)

	for rows.Next() {
		project := readmodel.Project{}

		err = rows.Scan(&project.Id, &project.Title, &project.Description, &project.Metadata.Author, &project.Metadata.Keywords, &project.Metadata.Language, &project.Metadata.Isbn)
		if err != nil {
			return nil, err
		}

		projectIds = append(projectIds, project.Id)
		projects[project.Id] = &project
	}

	genres, err := r.db.
		Table("genres AS g").
		Select("g.id AS genre_id", "g.name", "pg.project_id").
		Joins("join project2_genres pg on pg.genre_id = g.id").
		Where("pg.project_id IN ?", projectIds).
		Rows()

	if err != nil {
		return nil, err
	}
	defer genres.Close()
	for genres.Next() {
		genre := readmodel.Genre{}
		var projectId string

		err = genres.Scan(&genre.Id, &genre.Name, &projectId)
		if err != nil {
			return nil, err
		}

		projects[projectId].Genres = append(projects[projectId].Genres, genre)
	}

	var result []readmodel.Project
	for _, project := range projects {
		result = append(result, *project)
	}

	return result, nil
}

func (r *MySqlProjectReadModel) GetById(authorId string, projectId string) (readmodel.Project, error) {
	rows, err := r.db.Raw(`
        SELECT 
            p.id AS project_id, p.title, p.description,
            pm.author, pm.keywords, pm.language, pm.isbn,
            c.chapter_id AS chapter_id, c.title AS chapter_title, c.content AS chapter_content, c.order AS chapter_order,
            l.location_id AS location_id, l.name AS location_name, l.description AS location_description,
            ch.character_id AS character_id, ch.name AS character_name, ch.description AS character_description, ch.attributes AS character_attributes
        FROM projects2 p
        JOIN project_metadata pm ON p.id = pm.project_id
        LEFT JOIN chapters c ON c.project_id = p.id AND c.is_last_version = 1 AND c.deleted_at IS NULL
        LEFT JOIN locations l ON l.project_id = p.id AND l.is_last_version = 1 AND l.deleted_at IS NULL
        LEFT JOIN characters ch ON ch.project_id = p.id AND ch.is_last_version = 1 AND ch.deleted_at IS NULL
        WHERE p.id = ? AND p.author_id = ?
        ORDER BY c.order
    `, projectId, authorId).Rows()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return readmodel.Project{}, nil
		}
		return readmodel.Project{}, err
	}
	defer rows.Close()

	var project readmodel.Project
	seenChapters := make(map[string]struct{})
	seenLocations := make(map[string]struct{})
	seenCharacters := make(map[string]struct{})

	for rows.Next() {
		var projectId, title, description, author, keywords, language, isbn sql.NullString
		var chapterId, chapterTitle, chapterContent sql.NullString
		var chapterOrder sql.NullInt64
		var locationId, locationName, locationDescription sql.NullString
		var characterId, characterName, characterDescription, characterAttributes sql.NullString

		err := rows.Scan(&projectId, &title, &description, &author, &keywords, &language, &isbn,
			&chapterId, &chapterTitle, &chapterContent, &chapterOrder,
			&locationId, &locationName, &locationDescription,
			&characterId, &characterName, &characterDescription, &characterAttributes)
		if err != nil {
			return readmodel.Project{}, err
		}

		if project.Id == "" {
			project = readmodel.Project{
				Id:          projectId.String,
				Title:       title.String,
				Description: description.String,
				Metadata: readmodel.Metadata{
					Author:   author.String,
					Keywords: keywords.String,
					Language: language.String,
					Isbn:     isbn.String,
				},
			}
		}

		if chapterId.Valid {
			if _, exists := seenChapters[chapterId.String]; !exists {
				project.Chapters = append(project.Chapters, readmodel.Chapter{
					Id:      chapterId.String,
					Name:    chapterTitle.String,
					Content: chapterContent.String,
					Order:   int(chapterOrder.Int64),
				})
				seenChapters[chapterId.String] = struct{}{}
			}
		}

		if locationId.Valid {
			if _, exists := seenLocations[locationId.String]; !exists {

				var locDescription *string
				if locationDescription.Valid {
					locDescription = &locationDescription.String
				}
				project.Locations = append(project.Locations, readmodel.Location{
					Id:          locationId.String,
					Name:        locationName.String,
					Description: locDescription,
				})
				seenLocations[locationId.String] = struct{}{}
			}
		}

		if characterId.Valid {
			if _, exists := seenCharacters[characterId.String]; !exists {
				var attributes map[string]string
				if characterAttributes.Valid {
					json.Unmarshal([]byte(characterAttributes.String), &attributes)
				}
				var charDescription *string
				if characterDescription.Valid {
					charDescription = &characterDescription.String
				}
				project.Characters = append(project.Characters, readmodel.Character{
					Id:          characterId.String,
					Name:        characterName.String,
					Attributes:  attributes,
					Description: charDescription,
				})
				seenCharacters[characterId.String] = struct{}{}
			}
		}
	}

	if project.Id == "" {
		return readmodel.Project{}, nil
	}

	return project, nil
}
