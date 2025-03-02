package infrastructure

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gbooksoft/backend/internal/domain/project/readmodel"
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
		Table("projects AS p").
		Select("p.id AS project_id", "p.title", "p.description", "pm.author", "pm.keywords", "pm.language", "pm.isbn").
		Joins("JOIN project_metadata pm ON p.id = pm.project_id").
		Where("p.author_id = ?", authorId).
		Order("p.updated_at DESC").
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []readmodel.Project
	var projectIds []string

	for rows.Next() {
		var project readmodel.Project
		err = rows.Scan(&project.Id, &project.Title, &project.Description, &project.Metadata.Author, &project.Metadata.Keywords, &project.Metadata.Language, &project.Metadata.Isbn)
		if err != nil {
			return nil, err
		}

		projectIds = append(projectIds, project.Id)
		projects = append(projects, project)
	}

	if len(projectIds) == 0 {
		return projects, nil
	}

	genres, err := r.db.
		Table("genres AS g").
		Select("g.id AS genre_id", "g.name", "pg.project_id").
		Joins("JOIN project_genres pg ON pg.genre_id = g.id").
		Where("pg.project_id IN ?", projectIds).
		Rows()

	if err != nil {
		return nil, err
	}
	defer genres.Close()

	for genres.Next() {
		var genre readmodel.Genre
		var projectId string

		err = genres.Scan(&genre.Id, &genre.Name, &projectId)
		if err != nil {
			return nil, err
		}

		for i := range projects {
			if projects[i].Id == projectId {
				projects[i].Genres = append(projects[i].Genres, genre)
				break
			}
		}
	}

	return projects, nil
}

func (r *MySqlProjectReadModel) GetById(authorId string, projectId string) (readmodel.Project, error) {
	rows, err := r.db.Raw(`
        SELECT 
            p.id AS project_id, p.title, p.description,
            pm.author, pm.keywords, pm.language, pm.isbn,
            c.chapter_id AS chapter_id, c.title AS chapter_title, c.content AS chapter_content, c.order AS chapter_order, c.version AS chapter_version,
            l.location_id AS location_id, l.name AS location_name, l.description AS location_description,
            ch.character_id AS character_id, ch.name AS character_name, ch.description AS character_description, ch.attributes AS character_attributes,
            g.id AS genre_id, g.name AS genre_name
        FROM projects p
        JOIN project_metadata pm ON p.id = pm.project_id
		JOIN project_genres pg ON pg.project_id = p.id
		JOIN genres g ON g.id = pg.genre_id
        LEFT JOIN project_chapters c ON c.project_id = p.id AND c.is_last_version = 1 AND c.deleted_at IS NULL
        LEFT JOIN project_locations l ON l.project_id = p.id AND l.is_last_version = 1 AND l.deleted_at IS NULL
        LEFT JOIN project_characters ch ON ch.project_id = p.id AND ch.is_last_version = 1 AND ch.deleted_at IS NULL
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
	seenGenres := make(map[string]struct{})

	for rows.Next() {
		var projectId, title, description, author, keywords, language, isbn sql.NullString
		var chapterId, chapterTitle, chapterContent sql.NullString
		var chapterOrder sql.NullInt64
		var chapterVersion sql.NullInt64
		var locationId, locationName, locationDescription sql.NullString
		var characterId, characterName, characterDescription, characterAttributes sql.NullString
		var genreId, genreName string

		err := rows.Scan(&projectId, &title, &description, &author, &keywords, &language, &isbn,
			&chapterId, &chapterTitle, &chapterContent, &chapterOrder, &chapterVersion,
			&locationId, &locationName, &locationDescription,
			&characterId, &characterName, &characterDescription, &characterAttributes,
			&genreId, &genreName)
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

		if _, exists := seenGenres[genreId]; !exists {
			project.Genres = append(project.Genres, readmodel.Genre{
				Id:   genreId,
				Name: genreName,
			})
		}

		if chapterId.Valid {
			if _, exists := seenChapters[chapterId.String]; !exists {
				project.Chapters = append(project.Chapters, readmodel.Chapter{
					Id:      chapterId.String,
					Name:    chapterTitle.String,
					Content: chapterContent.String,
					Order:   int(chapterOrder.Int64),
					Version: int(chapterVersion.Int64),
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

func (r *MySqlProjectReadModel) GetPublication(authorId string, projectId string) (readmodel.Publication, error) {
	/*	query := `
		SELECT pc.chapter_id as project_chapter_id, bc.title, bc.version, bc.order
		FROM projects p
		JOIN book_chapters bc ON p.book_id = bc.book_id
		join project_chapters pc ON pc.book_chapter_id = bc.id AND pc.is_last_version = 1
		WHERE p.author_id = ? AND p.id = ? AND p.book_id IS NOT NULL
		ORDER BY bc.order
		`*/
	query := `
SELECT pc.chapter_id as project_chapter_id, bc.title, bc.version, bc.order 
FROM book_chapters bc
JOIN projects p ON p.book_id = bc.book_id
join project_chapters pc ON pc.book_chapter_id = bc.id AND pc.is_last_version = 1
WHERE p.author_id = ? AND p.id = ? AND p.book_id IS NOT NULL
ORDER BY bc.order 
`
	rows, err := r.db.Raw(query, authorId, projectId).Rows()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return readmodel.Publication{}, nil
		}
		return readmodel.Publication{}, err
	}
	defer rows.Close()

	var publication readmodel.Publication

	for rows.Next() {
		var chapter readmodel.PublicationChapter

		err := rows.Scan(&chapter.ProjectChapterId, &chapter.Title, &chapter.Version, &chapter.Order)
		if err != nil {
			return readmodel.Publication{}, err
		}

		publication.Chapters = append(publication.Chapters, chapter)
	}

	return publication, nil
}
