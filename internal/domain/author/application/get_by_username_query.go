package application

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gbooksoft/backend/internal/domain/author/readmodel"
	"github.com/gbooksoft/backend/pkg/gerr"
)

type GetAuthorByUsernameQuery struct {
	Username string
}

type GetAuthorByUsernameHandlerResponse struct {
	Author readmodel.Author
}

type GetAuthorByUsernameHandler struct {
	authorReadModel readmodel.AuthorReadModel
}

func NewGetAuthorByUsernameHandler(authorReadModel readmodel.AuthorReadModel) *GetAuthorByUsernameHandler {
	return &GetAuthorByUsernameHandler{authorReadModel: authorReadModel}
}

func (h *GetAuthorByUsernameHandler) Handle(ctx context.Context, query *GetAuthorByUsernameQuery) (*GetAuthorByUsernameHandlerResponse, error) {
	a, err := h.authorReadModel.GetByUsername(query.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, gerr.NewNotFoundErr("author not found")
		}
		return nil, err
	}

	return &GetAuthorByUsernameHandlerResponse{Author: a}, nil
}
