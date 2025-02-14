package domain

import (
	"errors"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gbooksoft/backend/pkg/model"
	"time"

	"github.com/google/uuid"
)

const (
	maxGenres = 5
)

type GenreId = model.Uuid
type AuthorId = model.Uuid
type ChapterId = model.Uuid
type ChapterInternalId = model.Uuid
type CharacterId = model.Uuid
type CharacterInternalId = model.Uuid
type LocationId = model.Uuid
type ProjectId = model.Uuid
type MetadataId = model.Uuid

type Chapter struct {
	Id                ChapterId         `gorm:"type:char(36);primaryKey"`
	ChapterInternalId ChapterInternalId `gorm:"type:char(36);index;column:chapter_id"`
	ProjectId         ProjectId         `gorm:"type:char(36);index"`
	Title             string
	Content           string
	Version           int
	DeletedAt         *time.Time
}

type Character struct {
	Id                  CharacterId         `gorm:"type:char(36);primaryKey"`
	CharacterInternalId CharacterInternalId `gorm:"type:char(36);index;column:character_id"`
	ProjectId           ProjectId           `gorm:"type:char(36);index"`
	Name                string
	Attributes          model.JSONMap `gorm:"type:json"`
	Description         *string
	Version             int
	DeletedAt           *time.Time
}

type Location struct {
	Id          LocationId `gorm:"type:char(36);primaryKey"`
	ProjectId   ProjectId  `gorm:"type:char(36);index"`
	Name        string
	Description string
	Version     int
	DeletedAt   *time.Time
}

type Metadata struct {
	Id        MetadataId `gorm:"type:char(36);primaryKey"`
	ProjectId ProjectId  `gorm:"type:char(36);index"`
	Author    string
	Genre     string
	Keywords  string
	Language  string
	Isbn      string
}

func (Metadata) TableName() string {
	return "project_metadata"
}

type Genre struct {
	Id   GenreId `gorm:"type:char(36);primaryKey"`
	Name string
}

type Project struct {
	Id          ProjectId `gorm:"type:char(36);primaryKey"`
	Title       string    `gorm:"not null"`
	Description string
	Genres      []Genre     `gorm:"many2many:project2_genres;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AuthorId    AuthorId    `gorm:"type:char(36);not null"`
	Chapters    []Chapter   `gorm:"foreignKey:project_id"`
	Characters  []Character `gorm:"foreignKey:project_id"`
	Locations   []Location  `gorm:"foreignKey:project_id"`
	Metadata    *Metadata   `gorm:"foreignKey:project_id"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *Project) TableName() string {
	return "projects2"
}

func NewProject(projectId ProjectId, title string, description string, genres []Genre, authorId AuthorId) (*Project, error) {
	if len(genres) > maxGenres {
		return nil, gerr.NewBusinessLogicErr("max 5 genres", gerr.ProjectShouldHaveMaxFiveGenres)
	}

	return &Project{
		Id:          projectId,
		Title:       title,
		Description: description,
		Genres:      genres,
		AuthorId:    authorId,
		Metadata: &Metadata{
			Id: MetadataId(uuid.New()),
		},
	}, nil
}

func (p *Project) AddChapter(newInternalId ChapterInternalId, title, content string) error {
	for _, c := range p.Chapters {
		if c.ChapterInternalId.Equal(newInternalId) {
			return gerr.NewBusinessLogicErr("can not create chapter with this id", gerr.ProjectChapterAlreadyExists)
		}
	}

	chapter := Chapter{
		Id:                ChapterId(uuid.New()),
		ChapterInternalId: newInternalId,
		Title:             title,
		Content:           content,
		Version:           1,
	}
	p.Chapters = append(p.Chapters, chapter)
	p.UpdatedAt = time.Now()

	return nil
}

func (p *Project) AddCharacter(newInternalId CharacterInternalId, name string, description *string, attributes map[string]string) error {
	for _, c := range p.Characters {
		if c.CharacterInternalId.Equal(newInternalId) {
			return gerr.NewBusinessLogicErr("can not create character with this id", gerr.ProjectChapterAlreadyExists)
		}
	}

	character := Character{
		Id:                  CharacterId(uuid.New()),
		CharacterInternalId: newInternalId,
		Name:                name,
		Description:         description,
		Attributes:          attributes,
		Version:             1,
	}
	p.Characters = append(p.Characters, character)
	p.UpdatedAt = time.Now()

	return nil
}

func (p *Project) AddLocation(name, description string) {
	location := Location{
		Id:          LocationId(uuid.New()),
		Name:        name,
		Description: description,
		Version:     1,
	}
	p.Locations = append(p.Locations, location)
	p.UpdatedAt = time.Now()
}

func (p *Project) UpdateChapter(chapterInternalId ChapterInternalId, newTitle, newContent string) error {
	lastVersion, ok := p.findLastVersionChapter(chapterInternalId)

	if !ok {
		return gerr.NewBusinessLogicErr("chapter not found", gerr.ProjectChapterNotFound)
	}

	if lastVersion.Title == newTitle && lastVersion.Content == newContent {
		return nil
	}

	newChapter := Chapter{
		Id:                ChapterId(uuid.New()),
		ChapterInternalId: chapterInternalId,
		ProjectId:         lastVersion.ProjectId,
		Title:             newTitle,
		Content:           newContent,
		Version:           lastVersion.Version + 1,
	}

	p.Chapters = append(p.Chapters, newChapter)

	return nil
}

func (p *Project) UpdateCharacter(characterId CharacterId, newName string, newDescription *string, newAttributes map[string]string) error {
	for i, c := range p.Characters {
		if c.Id.Equal(characterId) && c.DeletedAt == nil {
			updatedCharacter := c
			updatedCharacter.Name = newName
			updatedCharacter.Description = newDescription
			updatedCharacter.Attributes = newAttributes
			updatedCharacter.Version = updatedCharacter.Version + 1

			// Создаём новую версию с увеличенным номером
			p.Characters[i] = updatedCharacter
			p.UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("персонаж не найден или удалён")
}

func (p *Project) UpdateLocation(locationId LocationId, newName, newDescription string) error {
	for i, l := range p.Locations {
		if l.Id.Equal(locationId) && l.DeletedAt == nil {
			updatedLocation := l
			updatedLocation.Name = newName
			updatedLocation.Description = newDescription
			updatedLocation.Version = updatedLocation.Version + 1

			// Создаём новую версию с увеличенным номером
			p.Locations[i] = updatedLocation
			p.UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("локация не найдена или удалена")
}

func (p *Project) DeleteChapter(chapterInternalId ChapterInternalId) error {
	lastVersion, ok := p.findLastVersionChapter(chapterInternalId)

	if !ok {
		return gerr.NewBusinessLogicErr("chapter not found", gerr.ProjectChapterNotFound)
	}

	for i, c := range p.Chapters {
		if c.Id.Equal(lastVersion.Id) && c.DeletedAt == nil {
			now := time.Now()
			updatedChapter := c
			updatedChapter.DeletedAt = &now

			p.Chapters[i] = updatedChapter
			break
		}
	}

	return nil
}

func (p *Project) DeleteCharacter(characterId CharacterId) error {
	for i, c := range p.Characters {
		if c.Id.Equal(characterId) && c.DeletedAt == nil {
			now := time.Now()
			updatedCharacter := c
			updatedCharacter.DeletedAt = &now

			// Создаём новую версию с увеличенным номером
			p.Characters[i] = updatedCharacter
			p.UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("персонаж не найден или уже удалён")
}

func (p *Project) DeleteLocation(locationId LocationId) error {
	for i, l := range p.Locations {
		if l.Id.Equal(locationId) && l.DeletedAt == nil {
			now := time.Now()
			updatedLocation := l
			updatedLocation.DeletedAt = &now

			// Создаём новую версию с увеличенным номером
			p.Locations[i] = updatedLocation
			p.UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("локация не найдена или уже удалена")
}

func (p *Project) Update(title string, description *string, genres []Genre) error {
	if len(genres) > maxGenres {
		return gerr.NewBusinessLogicErr("max 5 genres", gerr.ProjectShouldHaveMaxFiveGenres)
	}
	if description != nil {
		p.Description = *description
	} else {
		p.Description = ""
	}

	p.Title = title
	p.Genres = genres
	p.UpdatedAt = time.Now()

	return nil
}

func (p *Project) findLastVersionChapter(chapterInternalId ChapterInternalId) (Chapter, bool) {
	var lastVersion Chapter
	var found bool
	maxVersion := 0

	for _, c := range p.Chapters {
		if c.ChapterInternalId.Equal(chapterInternalId) {
			if c.Version > maxVersion {
				lastVersion = c
				found = true
			}

			maxVersion = max(c.Version, maxVersion)
		}
	}

	if !found {
		return Chapter{}, false
	}

	if lastVersion.DeletedAt != nil {
		return Chapter{}, false

	}

	return lastVersion, true
}
