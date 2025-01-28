package response

import (
	"errors"
	"github.com/gbooksoft/backend/pkg/gerr"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func HandleErr(err error, w gin.ResponseWriter) {
	var bErr gerr.BusinessLogicErr

	if errors.As(err, &bErr) {
		Json(w).BusinessErr(bErr).Write()
		return
	}

	log.Error().Err(err).Msg("Internal Error")
	Json(w).InternalError().Write()
	return
}
