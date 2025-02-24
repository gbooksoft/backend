package domain

import (
	"github.com/gbooksoft/backend/pkg/model"
	"time"
)

type BookId = model.Uuid
type AuthorId = model.Uuid
type ChapterId = model.Uuid
type GenreId = model.Uuid

type Book struct {
	Id          BookId    `gorm:"type:char(36);primaryKey"`
	Title       string    `gorm:"not null"`
	AuthorId    AuthorId  `gorm:"type:char(36);not null"`
	Genres      []Genre   `gorm:"many2many:book_genres;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Chapters    []Chapter `gorm:"foreignKey:book_id"`
	PublishedAt *time.Time
	IsPublished bool
}

func NewBook(id BookId, title string, authorId AuthorId, genres []Genre, chapters []Chapter, isPublished bool) *Book {
	var publishedAt *time.Time
	if isPublished {
		now := time.Now()
		publishedAt = &now
	}

	return &Book{
		Id:          id,
		Title:       title,
		AuthorId:    authorId,
		Genres:      genres,
		Chapters:    chapters,
		PublishedAt: publishedAt,
		IsPublished: isPublished,
	}
}

type Chapter struct {
	Id      ChapterId `gorm:"type:char(36);primaryKey"`
	BookId  BookId    `gorm:"type:char(36);index"`
	Title   string
	Content string
	Order   int
}

func (*Chapter) TableName() string {
	return "book_chapters"
}

type Genre struct {
	Id   GenreId `gorm:"type:char(36);primaryKey"`
	Name string
}
