package response

import (
	"encoding/json"
	"github.com/gbooksoft/backend/internal/gerr"
	"github.com/gin-gonic/gin"
	"net/http"
)

type errResponse struct {
	Code int64  `json:"error_code"`
	Msg  string `json:"error_message"`
}

type errResponses struct {
	Code int               `json:"error_code"`
	Msgs map[string]string `json:"error_messages"`
}
type JsonResponse struct {
	w                  gin.ResponseWriter
	msg                string
	statusCode         int
	isErrResponse      bool
	validationMessages map[string]string
	code               int64
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

func (r *JsonResponse) ErrValidation(err error) *JsonResponse {
	messages := parseValidationErr(err)
	r.StatusCode(http.StatusBadRequest)
	r.isErrResponse = true
	r.validationMessages = messages
	return r
}

func (r *JsonResponse) BusinessErr(err gerr.BusinessLogicErr) *JsonResponse {
	if err.Code == gerr.NotFound {
		r.StatusCode(http.StatusNotFound)
		return r.ErrMsg("Not Found")
	}

	r.StatusCode(http.StatusUnprocessableEntity)
	r.code = err.Code.ToInt()
	return r.ErrMsg(err.Msg)
}

func (r *JsonResponse) Write() {
	r.w.Header().Add("Content-Type", "application/json")
	r.w.WriteHeader(r.statusCode)

	if r.isErrResponse {
		var rsp any

		if len(r.validationMessages) > 0 {
			rsp = errResponses{
				Code: r.statusCode,
				Msgs: r.validationMessages,
			}
		} else {
			rsp = errResponse{
				Code: r.code,
				Msg:  r.msg,
			}
		}

		byteRsp, _ := json.Marshal(rsp)

		r.w.Write(byteRsp)
		return
	}

	r.w.Write([]byte(r.msg))
}

func (r *JsonResponse) Unauthorized(message string) *JsonResponse {
	r.StatusCode(http.StatusUnauthorized)

	return r.ErrMsg("Unauthorized: " + message)
}

func Json(writer gin.ResponseWriter) *JsonResponse {
	return &JsonResponse{
		w:                  writer,
		statusCode:         200,
		isErrResponse:      false,
		validationMessages: make(map[string]string),
	}
}
