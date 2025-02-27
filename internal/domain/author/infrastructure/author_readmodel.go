package infrastructure

import (
	"database/sql"
	"errors"
	"github.com/gbooksoft/backend/internal/domain/author/readmodel"
	"gorm.io/gorm"
)

type MySqlAuthorReadModel struct {
	db *gorm.DB
}

func NewMySqlAuthorReadModel(db *gorm.DB) *MySqlAuthorReadModel {
	return &MySqlAuthorReadModel{
		db: db,
	}
}
func (r *MySqlAuthorReadModel) GetByUsername(username string) (readmodel.Author, error) {
	row := r.db.
		Table("users AS u").
		Select("u.id", "u.email", "u.username", "u.first_name", "u.last_name", "u.avatar_path").
		Where("u.username = ? AND confirmed = 1", username).
		Row()

	var author readmodel.Author

	var uAvatarPath sql.NullString

	err := row.Scan(&author.Id, &author.Email, &author.Username, &author.FirstName, &author.LastName, &uAvatarPath)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return readmodel.Author{}, nil
		}
		return readmodel.Author{}, err
	}

	author.AvatarPath = uAvatarPath.String

	return author, nil
}
