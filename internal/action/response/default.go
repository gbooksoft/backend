package response

import (
	"errors"
	"github.com/gbooksoft/backend/internal/gerr"
	"github.com/gin-gonic/gin"
)

func HandleErr(err error, w gin.ResponseWriter) {
	var bErr gerr.BusinessLogicErr

	if errors.As(err, &bErr) {
		Json(w).BusinessErr(bErr).Write()
		return
	}

	Json(w).InternalError().Write()
	return
}
