package domain

import (
	"errors"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gbooksoft/backend/pkg/model"
	"time"

	"github.com/google/uuid"
)

type GenreId = model.Uuid
type AuthorId = model.Uuid
type ChapterId = model.Uuid
type CharacterId = model.Uuid
type LocationId = model.Uuid
type ProjectId = model.Uuid
type MetadataId = model.Uuid

type Chapter struct {
	Id        ChapterId
	Title     string
	Content   string
	DeletedAt *time.Time
}

type Character struct {
	Id          CharacterId
	Name        string
	Attributes  map[string]string
	Description string
	DeletedAt   *time.Time
}

type Location struct {
	Id          LocationId
	Name        string
	Description string
	DeletedAt   *time.Time
}

type Version[T any] struct {
	Entity    T
	Version   int
	CreatedAt time.Time
}

type Metadata struct {
	Id       MetadataId
	Author   string
	Genre    string
	Keywords string
	Language string
	Isbn     string
}

func nextVersion[T any](prev Version[T], entity T) Version[T] {
	return Version[T]{
		Entity:    entity,
		Version:   prev.Version + 1, // Увеличиваем версию
		CreatedAt: time.Now(),
	}
}

func newVersion[T any](entity T) Version[T] {
	return Version[T]{
		Entity:    entity,
		Version:   1,
		CreatedAt: time.Now(),
	}
}

type Project struct {
	Id          ProjectId
	Title       string
	Description string
	Genres      []GenreId
	AuthorId    AuthorId
	Chapters    []Version[Chapter]
	Characters  []Version[Character]
	Locations   []Version[Location]
	Metadata    *Metadata
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewProject(projectId ProjectId, title string, description string, genres []GenreId, authorId AuthorId) (*Project, error) {
	if len(genres) > 5 {
		return nil, gerr.NewBusinessLogicErr("max 5 genres", gerr.ProjectShouldHaveMaxFiveGenres)
	}

	return &Project{
		Id:          projectId,
		Title:       title,
		Description: description,
		Genres:      genres,
		AuthorId:    authorId,
		Chapters:    []Version[Chapter]{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    &Metadata{},
	}, nil
}

func (p *Project) AddChapter(title, content string) {
	chapter := Chapter{
		Id:      ChapterId(uuid.New()),
		Title:   title,
		Content: content,
	}
	p.Chapters = append(p.Chapters, newVersion(chapter))
	p.UpdatedAt = time.Now()
}

func (p *Project) AddCharacter(name, description string, attributes map[string]string) {
	character := Character{
		Id:          CharacterId(uuid.New()),
		Name:        name,
		Description: description,
		Attributes:  attributes,
	}
	p.Characters = append(p.Characters, newVersion(character))
	p.UpdatedAt = time.Now()
}

func (p *Project) AddLocation(name, description string) {
	location := Location{
		Id:          LocationId(uuid.New()),
		Name:        name,
		Description: description,
	}
	p.Locations = append(p.Locations, newVersion(location))
	p.UpdatedAt = time.Now()
}

func (p *Project) UpdateChapter(chapterId ChapterId, newTitle, newContent string) error {
	for i, v := range p.Chapters {
		if v.Entity.Id.Equal(chapterId) && v.Entity.DeletedAt == nil {
			updatedChapter := v.Entity
			updatedChapter.Title = newTitle
			updatedChapter.Content = newContent

			// Создаём новую версию с увеличенным номером
			p.Chapters[i] = nextVersion(v, updatedChapter)
			p.UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("глава не найдена или удалена")
}

func (p *Project) UpdateCharacter(characterId CharacterId, newName, newDescription string, newAttributes map[string]string) error {
	for i, v := range p.Characters {
		if v.Entity.Id.Equal(characterId) && v.Entity.DeletedAt == nil {
			updatedCharacter := v.Entity
			updatedCharacter.Name = newName
			updatedCharacter.Description = newDescription
			updatedCharacter.Attributes = newAttributes

			// Создаём новую версию с увеличенным номером
			p.Characters[i] = nextVersion(v, updatedCharacter)
			p.UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("персонаж не найден или удалён")
}

func (p *Project) UpdateLocation(locationId LocationId, newName, newDescription string) error {
	for i, v := range p.Locations {
		if v.Entity.Id.Equal(locationId) && v.Entity.DeletedAt == nil {
			updatedLocation := v.Entity
			updatedLocation.Name = newName
			updatedLocation.Description = newDescription

			// Создаём новую версию с увеличенным номером
			p.Locations[i] = nextVersion(v, updatedLocation)
			p.UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("локация не найдена или удалена")
}

func (p *Project) DeleteChapter(chapterId ChapterId) error {
	for i, v := range p.Chapters {
		if v.Entity.Id.Equal(chapterId) && v.Entity.DeletedAt == nil {
			now := time.Now()
			updatedChapter := v.Entity
			updatedChapter.DeletedAt = &now

			// Создаём новую версию с увеличенным номером
			p.Chapters[i] = nextVersion(v, updatedChapter)
			p.UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("глава не найдена или уже удалена")
}

func (p *Project) DeleteCharacter(characterId CharacterId) error {
	for i, v := range p.Characters {
		if v.Entity.Id.Equal(characterId) && v.Entity.DeletedAt == nil {
			now := time.Now()
			updatedCharacter := v.Entity
			updatedCharacter.DeletedAt = &now

			// Создаём новую версию с увеличенным номером
			p.Characters[i] = nextVersion(v, updatedCharacter)
			p.UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("персонаж не найден или уже удалён")
}

func (p *Project) DeleteLocation(locationId LocationId) error {
	for i, v := range p.Locations {
		if v.Entity.Id.Equal(locationId) && v.Entity.DeletedAt == nil {
			now := time.Now()
			updatedLocation := v.Entity
			updatedLocation.DeletedAt = &now

			// Создаём новую версию с увеличенным номером
			p.Locations[i] = nextVersion(v, updatedLocation)
			p.UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("локация не найдена или уже удалена")
}
