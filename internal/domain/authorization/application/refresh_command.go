package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/rs/zerolog/log"
)

type RefreshCommand struct {
	RefreshToken string
}

type RefreshHandlerResponse struct {
	AccessToken  string
	RefreshToken string
}

type RefreshHandler struct {
	db *sql.DB
}

func NewRefreshHandler(db *sql.DB) *RefreshHandler {
	return &RefreshHandler{db}
}

func (h *RefreshHandler) Handle(ctx context.Context, cmd *RefreshCommand) (*RefreshHandlerResponse, error) {
	rt := cmd.RefreshToken

	_, rtClaims, err := ParseToken(rt, true)

	if err != nil {
		log.Error().Err(err).Msgf("can not parse refresh token")
		return nil, err
	}

	if IsTokenExpired(rtClaims) {
		return nil, gerr.NewBusinessLogicErr("token expired", gerr.AuthorizationRefreshTokenExpired)
	}

	var at string
	err = h.db.QueryRow("SELECT access FROM auth_token WHERE refresh = ?", rt).Scan(&at)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, gerr.NewBusinessLogicErr("refresh not found", gerr.AuthorizationRefreshTokenNotFound)
		}

		log.Error().Err(err).Msgf("can not refresh token")
		return nil, err
	}

	_, atClaims, err := ParseToken(at, true)

	if err != nil {
		log.Error().Err(err).Msgf("can not parse refresh token")
		return nil, err
	}

	err = deleteByRefreshToken(h.db, rt)
	if err != nil {
		log.Error().Err(err).Msgf("can not delete refresh token")
		return nil, err
	}

	at, rt, err = createTokens(h.db, atClaims.Subject, atClaims.Email, atClaims.FirstName, atClaims.LastName)

	if err != nil {
		log.Error().Err(err).Msgf("can not create access token when refreshing in %s", atClaims.Subject)
		return nil, fmt.Errorf("can not create access token when refreshing in %w", err)
	}
	return &RefreshHandlerResponse{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}
