package infrastructure

import (
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

func (r *MySqlProjectReadModel) GetById(userId string, projectId string) (readmodel.Project, error) {
	return readmodel.Project{}, nil
}
