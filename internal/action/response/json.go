package response

import (
	"encoding/json"
	"github.com/gbooksoft/backend/internal/gerr"
	"github.com/gin-gonic/gin"
	"net/http"
)

type errResponse struct {
	Code int    `json:"error_code"`
	Msg  string `json:"error_message"`
}
type JsonResponse struct {
	w             gin.ResponseWriter
	msg           string
	statusCode    int
	isErrResponse bool
}

func (r *JsonResponse) Struct(strct any) *JsonResponse {
	respByte, _ := json.Marshal(strct)

	return r.Msg(string(respByte))
}

func (r *JsonResponse) Msg(message string) *JsonResponse {
	r.msg = message
	return r
}

func (r *JsonResponse) ErrMsg(message string) *JsonResponse {
	r.isErrResponse = true
	return r.Msg(message)
}

func (r *JsonResponse) StatusCode(code int) *JsonResponse {
	r.statusCode = code
	return r
}

func (r *JsonResponse) InternalError() *JsonResponse {
	r.StatusCode(500)
	return r.ErrMsg("Internal Error")
}

func (r *JsonResponse) Write() {
	r.w.Header().Add("Content-Type", "application/json")
	r.w.WriteHeader(r.statusCode)

	if r.isErrResponse {
		rsp := errResponse{
			Code: r.statusCode,
			Msg:  r.msg,
		}

		byteRsp, _ := json.Marshal(rsp)

		r.w.Write(byteRsp)
		return
	}

	r.w.Write([]byte(r.msg))
}

func (r *JsonResponse) ErrValidation(err error) *JsonResponse {
	r.StatusCode(http.StatusBadRequest)
	return r.ErrMsg("Input validation error: " + err.Error())
}

func (r *JsonResponse) BusinessErr(err gerr.BusinessLogicErr) *JsonResponse {
	r.StatusCode(http.StatusUnprocessableEntity)
	return r.ErrMsg(err.Msg)
}

func Json(writer gin.ResponseWriter) *JsonResponse {
	return &JsonResponse{
		w:             writer,
		statusCode:    200,
		isErrResponse: false,
	}
}
