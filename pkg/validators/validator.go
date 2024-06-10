package validators

import (
	"context"
	"fmt"
	"github.com/go-playground/locales/ar"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ar_translations "github.com/go-playground/validator/v10/translations/ar"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"samm/pkg/validators/localization"
	"strings"
)

type Message struct {
	Text string `json:"text"`
	Code string `json:"code"`
}

type ErrorResponse struct {
	ValidationErrors   map[string][]string `json:"validation_errors"`
	IsError            bool                `json:"-"`
	ErrorMessageObject *Message            `json:"message"`
}
type Response struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  bool        `json:"status"`
}

var (
	transEn ut.Translator
	transAr ut.Translator
	langEn  = "en"
	langAr  = "ar"
)

func Init() *validator.Validate {
	validate := validator.New()
	en := en.New()
	ar := ar.New()
	uni := ut.New(en, ar)
	// this is usually know or extracted from http 'Accept-Language' header
	transEn, _ = uni.GetTranslator(langEn)
	transAr, _ = uni.GetTranslator(langAr)
	en_translations.RegisterDefaultTranslations(validate, transEn)
	ar_translations.RegisterDefaultTranslations(validate, transAr)

	//NewRegisterCustomValidator(validate)

	return validate
}
func GetTrans(c context.Context) ut.Translator {
	lang := c.Value("lang")
	switch lang {
	case langEn:
		return transEn
	case langAr:
		return transAr
	default:
		return transEn
	}
}

type CustomErrorTags struct {
	ValidationTag          string
	RegisterValidationFunc func(fl validator.FieldLevel) bool
}

func GetFiledName(e validator.FieldError) string {
	filedName := ""
	for i, s := range strings.Split(e.Namespace(), ".") {
		if i == 0 {
			continue
		}
		filedName = filedName + s

		if i != len(strings.Split(e.Namespace(), "."))-1 {
			filedName = filedName + "."
		}
	}
	return strings.ToLower(filedName)
}

func ValidateStruct(c context.Context, validate *validator.Validate, obj interface{}, customErrorTags ...CustomErrorTags) ErrorResponse {
	registerCustomValidation(c, validate, customErrorTags...)
	NewRegisterCustomValidator(c, validate)

	err := validate.Struct(obj)
	lang := c.Value("lang")
	fmt.Println(lang)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		errMap := make(map[string][]string)
		for _, e := range errs {

			filedName := GetFiledName(e)
			// can translate each error one at a time.
			if lang == langEn {
				errMap[filedName] = []string{e.Translate(transEn)}
			} else {
				errMap[filedName] = []string{e.Translate(transAr)}
			}

		}
		return ErrorResponse{
			ValidationErrors: errMap,
			IsError:          true,
		}
	}
	return ErrorResponse{}
}

func registerCustomValidation(c context.Context, validate *validator.Validate, customErrorTags ...CustomErrorTags) {
	for _, tag := range customErrorTags {
		validate.RegisterTranslation(tag.ValidationTag, GetTrans(c), func(ut ut.Translator) error {
			return nil
		}, func(ut ut.Translator, fe validator.FieldError) string {
			return localization.GetTranslation(&c, tag.ValidationTag, nil, ut.Locale())
		})
		validate.RegisterValidation(tag.ValidationTag, tag.RegisterValidationFunc)
		fmt.Println("registerCustomValidation", tag.ValidationTag)
	}
}

func GetErrorResponseFromErr(e error) ErrorResponse {
	return ErrorResponse{
		ValidationErrors: nil,
		IsError:          true,
		ErrorMessageObject: &Message{
			Text: e.Error(),
			Code: "",
		},
	}
}

func GetErrorResponse(ctx *context.Context, code string, data map[string]interface{}) ErrorResponse {
	message := localization.GetTranslation(ctx, code, data, "")
	return ErrorResponse{
		ValidationErrors: nil,
		IsError:          true,
		ErrorMessageObject: &Message{
			Text: message,
			Code: code,
		},
	}
}

func Success(c echo.Context, data any) error {
	c.JSON(http.StatusOK, data)
	return nil
}
func SuccessResponse(c echo.Context, data any) error {
	if data == nil {
		data = make(map[string]interface{})
	}

	res := Response{Status: true, Message: "Success", Data: data}
	c.JSON(http.StatusOK, res)
	return nil
}

func ErrorStatusUnprocessableEntity(c echo.Context, validationErr ErrorResponse) error {
	c.JSON(http.StatusUnprocessableEntity, validationErr)
	return errors.New("ErrorStatusUnprocessableEntity")
}
func ErrorStatusBadRequest(c echo.Context, validationErr ErrorResponse) error {
	c.JSON(http.StatusBadRequest, validationErr)
	return errors.New("ErrorStatusBadRequest")
}
func ErrorStatusInternalServerError(c echo.Context, validationErr ErrorResponse) error {
	c.JSON(http.StatusInternalServerError, validationErr)
	return errors.New("ErrorStatusInternalServerError")
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
