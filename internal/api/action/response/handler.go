package response

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func HandleErr(err error, w gin.ResponseWriter) {
	log.Error().Err(err).Msg("Internal Error")
	Json(w).InternalError().Write()
	return
}
