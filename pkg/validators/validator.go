package validators

import (
	"github.com/go-playground/locales/ar"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ar_translations "github.com/go-playground/validator/v10/translations/ar"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/labstack/echo/v4"
)

type Message struct {
	En   string `json:"en"`
	Ar   string `json:"ar"`
	Code string `json:"code"`
}

type ErrorResponse struct {
	ValidationErrors   map[string][]string `json:"validation_errors"`
	IsError            bool                `json:"-"`
	ErrorMessageObject *Message            `json:"message"`
}

var (
	transEn ut.Translator
	transAr ut.Translator
)

func Init() *validator.Validate {
	validate := validator.New()
	en := en.New()
	ar := ar.New()
	uni := ut.New(en, ar)
	// this is usually know or extracted from http 'Accept-Language' header
	transEn, _ = uni.GetTranslator("en")
	transAr, _ = uni.GetTranslator("ar")
	en_translations.RegisterDefaultTranslations(validate, transEn)
	ar_translations.RegisterDefaultTranslations(validate, transAr)

	return validate
}

func ValidateStruct(c echo.Context, validate *validator.Validate, obj interface{}) ErrorResponse {
	err := validate.Struct(obj)
	lang := c.Request().Context().Value("lang")
	if err != nil {
		errs := err.(validator.ValidationErrors)
		errMap := make(map[string][]string)
		for _, e := range errs {
			// can translate each error one at a time.
			if lang == "en" {
				errMap[e.Field()] = []string{e.Translate(transEn)}
			} else {
				errMap[e.Field()] = []string{e.Translate(transAr)}
			}

		}
		return ErrorResponse{
			ValidationErrors: errMap,
			IsError:          true,
		}
	}
	return ErrorResponse{}
}

//package validators
//
//import (
//	"fmt"
//	"github.com/go-playground/validator/v10"
//	"github.com/labstack/echo/v4"
//	"github.com/nicksnyder/go-i18n/v2/i18n"
//	"golang.org/x/text/language"
//	"net/http"
//	"reflect"
//	"strings"
//)
//
//// Name of the struct tag used in examples
//const tagName = "message_key"
//
//type fieldError struct {
//	err        validator.FieldError
//	messageKey string
//}
//type validateErrorResponse struct {
//	StatusCode int    `json:"status_code"`
//	Message    string `json:"message"`
//	Key        string `json:"error_key"`
//}
//
//func (q fieldError) toI18nMessage(c echo.Context, i18nBundle *i18n.Bundle) (string, error) {
//	lang := c.Request().FormValue("lang")
//	accept := c.Request().Header["Accept-Language"]
//	localize := i18n.NewLocalizer(i18nBundle, language.English.String())
//	if len(lang) > 0 || len(accept) > 0 {
//		localize = i18n.NewLocalizer(i18nBundle, lang, accept[0])
//	}
//
//	// get from lang from request
//	localizeValidateMessage := i18n.LocalizeConfig{
//		MessageID: q.messageKey,
//		TemplateData: map[string]string{
//			"ErrorTag":             q.err.Tag(),
//			"ErrorParam":           q.err.Param(),
//			"ErrorField":           q.err.Field(),
//			"ErrorStructField":     q.err.StructField(),
//			"ErrorActualTag":       q.err.ActualTag(),
//			"ValidateError":        q.err.Error(),
//			"ErrorNamespace":       q.err.Namespace(),
//			"ErrorStructNamespace": q.err.StructNamespace(),
//		},
//	}
//
//	message, err := localize.Localize(&localizeValidateMessage)
//	return message, err
//}
//
//func (q fieldError) toString(c echo.Context, i18nBundle *i18n.Bundle) string {
//	var sb strings.Builder
//	// load message from message file
//
//	message, err := q.toI18nMessage(c, i18nBundle)
//
//	if len(message) > 0 && err == nil {
//		return message
//	}
//
//	sb.WriteString("validation failed on field '" + q.err.Field() + "'")
//	sb.WriteString(", condition: " + q.err.ActualTag())
//
//	if q.err.Param() != "" {
//		sb.WriteString(" { " + q.err.Param() + " }")
//	}
//
//	if q.err.Value() != nil && q.err.Value() != "" {
//		sb.WriteString(fmt.Sprintf(", actual: %v", q.err.Value()))
//	}
//
//	return sb.String()
//}
//
//func parseField(inputType reflect.StructField, tagName string) []string {
//	// Get the field tag value
//	tags := inputType.Tag.Get(tagName)
//	tags = strings.Replace(tags, " ", "", -1)
//	return strings.Split(tags, ",")
//}
//
//func ValidationRender[T any](c echo.Context, err error, request T) {
//	var i18nBundle = i18n.NewBundle(language.English)
//	t := reflect.TypeOf(request)
//	mapMessage := map[string]string{}
//	for i := 0; i < t.NumField(); i++ {
//		field := t.Field(i)
//		tagMessageMap := parseField(field, tagName)
//		tagBindingMap := parseField(field, "json")
//		fmt.Println(tagBindingMap, "11111")
//
//		for i := 0; i < len(tagBindingMap); i++ {
//			if i < len(tagMessageMap) {
//				tag := strings.Split(tagBindingMap[i], "=")
//				mapMessage[field.Name+"_"+tag[0]] = tagMessageMap[i]
//				fmt.Println(tag, "221111")
//				fmt.Println(mapMessage, "331111")
//			}
//		}
//	}
//	fmt.Println(err, "991111")
//	fmt.Println(err.(validator.ValidationErrors), "881111")
//	for _, fieldErr := range err.(validator.ValidationErrors) {
//		mess := fieldError{fieldErr, mapMessage[fieldErr.Field()+"_"+fieldErr.Tag()]}.toString(c, i18nBundle)
//		errResponse := validateErrorResponse{
//			StatusCode: http.StatusUnprocessableEntity,
//			Message:    mess,
//			Key:        "validate_err",
//		}
//		fmt.Println(mess, "441111")
//		c.JSON(http.StatusUnprocessableEntity, errResponse)
//		return
//	}
//}
