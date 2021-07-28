package ginerrors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	errs "github.com/spacetab-io/errors-go"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type (
	langName         string
	validationRule   string
	errorPattern     string
	validationErrors map[validationRule]errorPattern
)

func (ve errorPattern) string() string {
	return string(ve)
}

var CommonValidationErrors = map[langName]validationErrors{
	"ru": {
		"ek":       "Ошибка валидации для свойства `%s` с правилом `%s`",
		"required": "Свойство `%s` обязательно для заполнения",
		"gt":       "Свойство `%s` должно содержать более `%s` элементов",
	},
	"en": {
		"ek":       "Fail to validate field `%s` with rule `%s`",
		"required": "Field `%s` is required",
		"gt":       "Field `%s` must contain more than `%s` elements",
	},
	"es": {
		"ek":       "",
		"required": "",
		"gt":       "",
	},
}

var (
	defaultLang = "ru"

	ErrNotFound       = errors.New("route not found")
	ErrNoMethod       = errors.New("method not allowed")
	ErrServerError    = errors.New("internal server error")
	ErrRecordNotFound = errors.New("record not found")
	ErrUnknownErrVal  = errors.New("unknown error value")
)

func getLang(c *gin.Context) langName {
	lang := c.GetHeader("lang")
	if lang == "" {
		lang = c.DefaultQuery("lang", defaultLang)
	}

	return langName(lang)
}

// validationErrors Формирование массива ошибок.
func makeErrorsSlice(err error, lang langName) map[errs.FieldName][]errs.ValidationError {
	ve := make(map[errs.FieldName][]errs.ValidationError)

	//nolint: errorlint
	verrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	for _, e := range verrs {
		field := getFieldName(e.Namespace(), e.Field())
		if _, ok := ve[field]; !ok {
			ve[field] = make([]errs.ValidationError, 0)
		}

		ve[field] = append(
			ve[field],
			getErrMessage(validationRule(e.ActualTag()), field, e.Param(), lang),
		)
	}

	return ve
}

func makeErrorsSliceFromViolations(violations []*errdetails.BadRequest_FieldViolation) map[errs.FieldName][]errs.ValidationError {
	ve := make(map[errs.FieldName][]errs.ValidationError)

	for _, v := range violations {
		if v == nil {
			continue
		}

		field := errs.FieldName(v.Field)
		if _, ok := ve[field]; !ok {
			ve[field] = make([]errs.ValidationError, 0)
		}

		e := errs.ValidationError(v.Description)
		ve[field] = append(ve[field], e)
	}

	return ve
}

func getFieldName(namespace string, field string) errs.FieldName {
	namespace = strings.ReplaceAll(namespace, "]", "")
	namespace = strings.ReplaceAll(namespace, "[", ".")
	namespaceSlice := strings.Split(namespace, ".")
	fieldName := field

	if len(namespaceSlice) > 2 { //nolint: gomnd
		fieldName = strings.Join([]string{strings.Join(namespaceSlice[1:len(namespaceSlice)-1], "."), field}, ".")
	}

	return errs.FieldName(fieldName)
}

func getErrMessage(errorType validationRule, field errs.FieldName, param string, lang langName) errs.ValidationError {
	errKey := errorType

	if _, ok := CommonValidationErrors[lang][errorType]; !ok {
		errKey = "ek"
	}

	if param != "" && errKey == "ek" {
		return errs.ValidationError(fmt.Sprintf(CommonValidationErrors[lang][errKey].string(), field, errorType))
	}

	return errs.ValidationError(fmt.Sprintf(CommonValidationErrors[lang][errKey].string(), field))
}
