package application

import (
	"context"
	"database/sql"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/rs/zerolog/log"
)

type LogoutCommand struct {
	RefreshToken string
}

type LogoutHandlerResponse struct{}

type LogoutHandler struct {
	db *sql.DB
}

func NewLogoutHandler(db *sql.DB) *LogoutHandler {
	return &LogoutHandler{db}
}

func (h *LogoutHandler) Handle(ctx context.Context, cmd *LogoutCommand) (*LogoutHandlerResponse, error) {
	rt := cmd.RefreshToken

	_, rtClaims, err := ParseToken(rt, false)

	if err != nil {
		log.Error().Err(err).Msgf("can not parse refresh token")
		return nil, err
	}

	if IsTokenExpired(rtClaims) {
		return nil, gerr.NewBusinessLogicErr("token expired", gerr.AuthorizationRefreshTokenExpired)
	}

	err = deleteByRefreshToken(h.db, rt)
	if err != nil {
		log.Error().Err(err).Msgf("can not delete refresh token")
		return nil, err
	}

	return nil, nil
}
