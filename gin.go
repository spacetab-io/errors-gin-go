package ginerrors

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spacetab-io/errors-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	validationErrorMessage = "validation error"
)

// Response makes common error response.
func Response(c *gin.Context, err interface{}) {
	errCode, data := MakeResponse(err, getLang(c))
	resp := errors.Response{Error: data}
	c.AbortWithStatusJSON(errCode, resp)
}

// MakeResponse makes ErrorObject based on error type.
func MakeResponse(err interface{}, lang langName) (int, errors.ErrorObject) {
	errObj := errors.ErrorObject{}
	errCode := http.StatusBadRequest
	errType := errors.ErrorTypeError

	switch et := err.(type) {
	case GRPCValidationError:
		errCode = http.StatusUnprocessableEntity

		errObj.Message = validationErrorMessage
		errObj.Validation = makeErrorsSliceFromViolations(et.Violations)
	case []error:
		errCode = http.StatusInternalServerError
		msgs := make([]string, 0)

		for _, e := range err.([]error) {
			msgs = append(msgs, e.Error())
		}

		errObj.Message = strings.Join(msgs, "; ")
	case validator.ValidationErrors:
		errCode = http.StatusUnprocessableEntity

		errObj.Message = validationErrorMessage
		errObj.Validation = makeErrorsSlice(et, lang)
	case error:
		st := status.Convert(et)
		if st.Code() != codes.Unknown {
			return MakeResponse(UnwrapRPCError(st), lang)
		}

		errCode, errObj.Message = getErrCode(et)
	case string:
		errObj.Message = et
	case map[string]error:
		msgs := make(map[string]string)
		for k, e := range et {
			msgs[k] = e.Error()
		}

		errObj.Message = msgs
	}

	errObj.Type = &errType

	return errCode, errObj
}

func getErrCode(et error) (errCode int, msg string) {
	msg = et.Error()
	switch msg {
	case ErrNotFound.Error():
		errCode = http.StatusNotFound
	case ErrNoMethod.Error():
		errCode = http.StatusMethodNotAllowed
	case ErrServerError.Error(), sql.ErrConnDone.Error(), sql.ErrTxDone.Error():
		errCode = http.StatusInternalServerError
	case ErrRecordNotFound.Error():
		errCode = http.StatusNotFound
	case sql.ErrNoRows.Error():
		errCode = http.StatusNotFound
		msg = ErrRecordNotFound.Error()
	default:
		errCode = http.StatusBadRequest
	}

	return
}
