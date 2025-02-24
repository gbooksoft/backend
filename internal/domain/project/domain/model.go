package domain

import (
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
type LocationInternalId = model.Uuid
type ProjectId = model.Uuid
type MetadataId = model.Uuid
type BookId = model.Uuid
type BookChapterId = model.Uuid

type versional interface {
	getInternalId() model.Uuid
	isLastVersion() bool
	isDeleted() bool
}

type Chapter struct {
	Id                ChapterId         `gorm:"type:char(36);primaryKey"`
	ChapterInternalId ChapterInternalId `gorm:"type:char(36);index;column:chapter_id"`
	ProjectId         ProjectId         `gorm:"type:char(36);index"`
	BookChapterId     *BookChapterId    `gorm:"type:char(36);index;column:book_chapter_id"`
	Title             string
	Content           string
	ContentHtml       string
	Version           int
	IsLastVersion     bool
	Order             int `gorm:"not null,default:1"`
	DeletedAt         *time.Time
}

func (Chapter) TableName() string {
	return "project_chapters"
}

func (c Chapter) getInternalId() ChapterInternalId {
	return c.ChapterInternalId
}

func (c Chapter) isLastVersion() bool {
	return c.IsLastVersion
}

func (c Chapter) isDeleted() bool {
	return c.DeletedAt != nil
}

type Character struct {
	Id                  CharacterId         `gorm:"type:char(36);primaryKey"`
	CharacterInternalId CharacterInternalId `gorm:"type:char(36);index;column:character_id"`
	ProjectId           ProjectId           `gorm:"type:char(36);index"`
	Name                string
	Attributes          model.JSONMap `gorm:"type:json"`
	Description         *string
	Version             int
	IsLastVersion       bool
	DeletedAt           *time.Time
}

func (Character) TableName() string {
	return "project_characters"
}

func (c Character) getInternalId() CharacterInternalId {
	return c.CharacterInternalId
}

func (c Character) isLastVersion() bool {
	return c.IsLastVersion
}

func (c Character) isDeleted() bool {
	return c.DeletedAt != nil
}

type Location struct {
	Id                 LocationId         `gorm:"type:char(36);primaryKey"`
	LocationInternalId LocationInternalId `gorm:"type:char(36);index;column:location_id"`
	ProjectId          ProjectId          `gorm:"type:char(36);index"`
	Name               string
	Description        *string
	Version            int
	IsLastVersion      bool
	DeletedAt          *time.Time
}

func (Location) TableName() string {
	return "project_locations"
}

func (l Location) getInternalId() LocationInternalId {
	return l.LocationInternalId
}

func (l Location) isLastVersion() bool {
	return l.IsLastVersion
}

func (l Location) isDeleted() bool {
	return l.DeletedAt != nil
}

type Metadata struct {
	Id        MetadataId `gorm:"type:char(36);primaryKey"`
	ProjectId ProjectId  `gorm:"type:char(36);index"`
	Author    string
	Keywords  string
	Language  string
	Isbn      string
}

func (*Metadata) TableName() string {
	return "project_metadata"
}

func (m *Metadata) fulfill(author, keywords, language, isbn string) {
	m.Author = author
	m.Keywords = keywords
	m.Language = language
	m.Isbn = isbn
}

type Genre struct {
	Id   GenreId `gorm:"type:char(36);primaryKey"`
	Name string
}

type Project struct {
	Id          ProjectId `gorm:"type:char(36);primaryKey"`
	Title       string    `gorm:"not null"`
	Description string
	Genres      []Genre     `gorm:"many2many:project_genres;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AuthorId    AuthorId    `gorm:"type:char(36);not null"`
	Chapters    []Chapter   `gorm:"foreignKey:project_id"`
	Characters  []Character `gorm:"foreignKey:project_id"`
	Locations   []Location  `gorm:"foreignKey:project_id"`
	Metadata    *Metadata   `gorm:"foreignKey:project_id"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	BookId      *BookId `gorm:"type:char(36);index"`
}

func (p *Project) TableName() string {
	return "projects"
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

func (p *Project) AddChapter(newInternalId ChapterInternalId, title, content, contentHtml string) error {
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
		ContentHtml:       contentHtml,
		Version:           1,
		IsLastVersion:     true,
		Order:             p.findLastOrderChapter() + 1,
	}
	p.Chapters = append(p.Chapters, chapter)
	p.UpdatedAt = time.Now()

	return nil
}

func (p *Project) AddCharacter(newInternalId CharacterInternalId, name string, description *string, attributes map[string]string) error {
	for _, c := range p.Characters {
		if c.CharacterInternalId.Equal(newInternalId) {
			return gerr.NewBusinessLogicErr("can not create character with this id", gerr.ProjectCharacterAlreadyExists)
		}
	}

	character := Character{
		Id:                  CharacterId(uuid.New()),
		CharacterInternalId: newInternalId,
		Name:                name,
		Description:         description,
		Attributes:          attributes,
		IsLastVersion:       true,
		Version:             1,
	}
	p.Characters = append(p.Characters, character)
	p.UpdatedAt = time.Now()

	return nil
}

func (p *Project) AddLocation(newInternalId LocationInternalId, name string, description *string) error {
	for _, l := range p.Locations {
		if l.LocationInternalId.Equal(newInternalId) {
			return gerr.NewBusinessLogicErr("can not create location with this id", gerr.ProjectLocationAlreadyExists)
		}
	}

	location := Location{
		Id:                 LocationId(uuid.New()),
		LocationInternalId: newInternalId,
		Name:               name,
		Description:        description,
		IsLastVersion:      true,
		Version:            1,
	}
	p.Locations = append(p.Locations, location)
	p.UpdatedAt = time.Now()

	return nil
}

func (p *Project) UpdateChapter(chapterInternalId ChapterInternalId, newTitle, newContent, newContentHtml string) error {
	lastVersion, ok := findLastVersion(p.Chapters, chapterInternalId)

	if !ok {
		return gerr.NewBusinessLogicErr("chapter not found", gerr.ProjectChapterNotFound)
	}

	if lastVersion.Title == newTitle && lastVersion.Content == newContent && lastVersion.ContentHtml == newContentHtml {
		return nil
	}

	for i, c := range p.Chapters {
		if c.Id.Equal(lastVersion.Id) {
			p.Chapters[i].IsLastVersion = false
		}
	}

	newChapter := Chapter{
		Id:                ChapterId(uuid.New()),
		ChapterInternalId: chapterInternalId,
		ProjectId:         lastVersion.ProjectId,
		Title:             newTitle,
		Content:           newContent,
		ContentHtml:       newContentHtml,
		Version:           lastVersion.Version + 1,
		Order:             lastVersion.Order,
		IsLastVersion:     true,
	}

	p.Chapters = append(p.Chapters, newChapter)
	p.UpdatedAt = time.Now()

	return nil
}

func (p *Project) UpdateCharacter(characterInternalId CharacterInternalId, newName string, newDescription *string, newAttributes map[string]string) error {
	lastVersion, ok := findLastVersion(p.Characters, characterInternalId)

	if !ok {
		return gerr.NewBusinessLogicErr("character not found", gerr.ProjectCharacterNotFound)
	}

	for i, c := range p.Characters {
		if c.Id.Equal(lastVersion.Id) {
			p.Characters[i].IsLastVersion = false
		}
	}

	newCharacter := Character{
		Id:                  CharacterId(uuid.New()),
		CharacterInternalId: characterInternalId,
		ProjectId:           lastVersion.ProjectId,
		Name:                newName,
		Description:         newDescription,
		Attributes:          newAttributes,
		Version:             lastVersion.Version + 1,
		IsLastVersion:       true,
	}

	p.Characters = append(p.Characters, newCharacter)
	p.UpdatedAt = time.Now()

	return nil
}

func (p *Project) UpdateLocation(locationInternalId LocationInternalId, newName string, newDescription *string) error {
	lastVersion, ok := findLastVersion(p.Locations, locationInternalId)

	if !ok {
		return gerr.NewBusinessLogicErr("character not found", gerr.ProjectCharacterNotFound)
	}

	for i, l := range p.Locations {
		if l.Id.Equal(lastVersion.Id) {
			p.Locations[i].IsLastVersion = false
		}
	}

	newLocation := Location{
		Id:                 LocationId(uuid.New()),
		LocationInternalId: locationInternalId,
		ProjectId:          lastVersion.ProjectId,
		Name:               newName,
		Description:        newDescription,
		Version:            lastVersion.Version + 1,
		IsLastVersion:      true,
	}

	p.Locations = append(p.Locations, newLocation)
	p.UpdatedAt = time.Now()

	return nil
}

func (p *Project) DeleteChapter(chapterInternalId ChapterInternalId) error {
	lastVersion, ok := findLastVersion(p.Chapters, chapterInternalId)

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
	p.UpdatedAt = time.Now()

	return nil
}

func (p *Project) DeleteCharacter(characterInternalId CharacterInternalId) error {
	lastVersion, ok := findLastVersion(p.Characters, characterInternalId)

	if !ok {
		return gerr.NewBusinessLogicErr("character not found", gerr.ProjectCharacterNotFound)
	}

	for i, c := range p.Characters {
		if c.Id.Equal(lastVersion.Id) && c.DeletedAt == nil {
			now := time.Now()
			updatedCharacter := c
			updatedCharacter.DeletedAt = &now

			p.Characters[i] = updatedCharacter
			break
		}
	}
	p.UpdatedAt = time.Now()

	return nil
}

func (p *Project) DeleteLocation(locationInternalId LocationInternalId) error {
	lastVersion, ok := findLastVersion(p.Locations, locationInternalId)

	if !ok {
		return gerr.NewBusinessLogicErr("character not found", gerr.ProjectCharacterNotFound)
	}

	for i, l := range p.Locations {
		if l.Id.Equal(lastVersion.Id) && l.DeletedAt == nil {
			now := time.Now()
			updatedLocation := l
			updatedLocation.DeletedAt = &now

			p.Locations[i] = updatedLocation
			break
		}
	}
	p.UpdatedAt = time.Now()

	return nil
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

func findLastVersion[T versional](slice []T, internalId model.Uuid) (T, bool) {
	var lastVersion T
	var found bool

	for _, c := range slice {
		if c.getInternalId().Equal(internalId) {
			if c.isLastVersion() {
				lastVersion = c
				found = true
				break
			}
		}
	}

	if !found {
		return lastVersion, false
	}

	if lastVersion.isDeleted() {
		return lastVersion, false
	}

	return lastVersion, true
}

func (p *Project) UpdateMetadata(author, keywords, language, isbn string) {
	if p.Metadata == nil {
		p.Metadata = &Metadata{}
	}

	p.Metadata.fulfill(
		author,
		keywords,
		language,
		isbn,
	)
	p.UpdatedAt = time.Now()
}

func (p *Project) ReorderChapter(chapterInternalId ChapterInternalId, newOrder int) {
	lastVersion, ok := findLastVersion(p.Chapters, chapterInternalId)
	if !ok {
		return
	}

	oldOrder := lastVersion.Order

	lastVersion.Order = newOrder
	p.UpdatedAt = time.Now()

	if newOrder < oldOrder {
		for i, c := range p.Chapters {
			if !c.IsLastVersion || c.DeletedAt != nil {
				continue
			}

			if c.Order >= newOrder && c.Order < oldOrder {
				n := c
				n.Order = n.Order + 1
				p.Chapters[i] = n
			}

			if c.Id.Equal(lastVersion.Id) {
				n := c
				n.Order = newOrder
				p.Chapters[i] = n
			}
		}

		return
	}

	if newOrder > oldOrder {
		for i, c := range p.Chapters {
			if !c.IsLastVersion || c.DeletedAt != nil {
				continue
			}

			if c.Order > oldOrder && c.Order <= newOrder {
				n := c
				n.Order = n.Order - 1
				p.Chapters[i] = n
			}

			if c.Id.Equal(lastVersion.Id) {
				n := c
				n.Order = newOrder
				p.Chapters[i] = n
			}
		}

		return
	}
}

func (p *Project) findLastOrderChapter() int {
	var order int
	for _, c := range p.Chapters {
		if !c.IsLastVersion || c.DeletedAt != nil {
			continue
		}

		order = max(c.Order, order)
	}

	return order
}
