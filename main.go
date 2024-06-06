package main

import (
	"go.uber.org/fx"
	"samm/internal/module/example"
	"samm/internal/module/menu"
	"samm/pkg/config"
	"samm/pkg/database"
	"samm/pkg/http"
	"samm/pkg/http/echo"
	echoserver "samm/pkg/http/echo/server"
	httpclient "samm/pkg/http_client"
	"samm/pkg/logger"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

func main() {
	fx.New(
		fx.Options(
			fx.Provide(
				config.Init,
				logger.InitLogger,
				http.NewContext,
				echoserver.NewEchoServer,
				httpclient.NewHttpClient,
				validators.Init,
			),
			example.Module,
			menu.Module,
			database.Module,
			fx.Invoke(echo.RunServers),
			fx.Invoke(localization.InitLocalization),
		),
	).Run()
}

//
// You can edit this code!
// Click here and start typing.
// ####################################
// package main
//
// import (
//
//	"errors"
//	"fmt"
//	"reflect"
//	"strings"
//
//	"github.com/go-playground/validator/v10"
//
// )
//
// // Help valid struct with tag custom
// const tagCustom = "errormgs"
// const jsonTag = "json"
//
//	func errorTagFunc[T interface{}](obj interface{}, snp string, fieldname, actualTag string) error {
//		o := obj.(T)
//
//		if !strings.Contains(snp, fieldname) {
//			return nil
//		}
//
//		fieldArr := strings.Split(snp, ".")
//		rsf := reflect.TypeOf(o)
//
//		for i := 1; i < len(fieldArr); i++ {
//			field, found := rsf.FieldByName(fieldArr[i])
//			if found {
//				if fieldArr[i] == fieldname {
//					customMessage := field.Tag.Get(tagCustom)
//					fieldJsonName := field.Tag.Get(jsonTag)
//					if customMessage != "" {
//						return fmt.Errorf("%s: %s", fieldJsonName, customMessage)
//					}
//					return nil
//				} else {
//					if field.Type.Kind() == reflect.Ptr {
//						// If the field type is a pointer, dereference it
//						rsf = field.Type.Elem()
//					} else {
//						rsf = field.Type
//					}
//				}
//			}
//		}
//		return nil
//	}
//
//	func ValidateFunc[T interface{}](obj interface{}, validate *validator.Validate) (errs error) {
//		o := obj.(T)
//
//		defer func() {
//			if r := recover(); r != nil {
//				fmt.Println("Recovered in Validate:", r)
//				errs = fmt.Errorf("can't validate %+v", r)
//			}
//		}()
//
//		if err := validate.Struct(o); err != nil {
//			errorValid := err.(validator.ValidationErrors)
//			for _, e := range errorValid {
//				// snp  X.Y.Z
//				snp := e.StructNamespace()
//				errmgs := errorTagFunc[T](obj, snp, e.Field(), e.ActualTag())
//				if errmgs != nil {
//					errs = errors.Join(errs, fmt.Errorf("%w", errmgs))
//				} else {
//					errs = errors.Join(errs, fmt.Errorf("%w", e))
//				}
//			}
//		}
//
//		if errs != nil {
//			return errs
//		}
//
//		return nil
//	}
//
//	type PaymentInfo struct {
//		CreditCardNumber string   `json:"credit_card_number" validate:"required" errormgs:"Invalid credit card is required xxx"`
//		CVV              string   `json:"cvv" validate:"required,len=3" errormgs:"CVV code must be three digits long"`
//		Age              *int     `json:"age" validate:"required" errormgs:"int age custom"`
//		Message          *Message `json:"message" validate:"required" errormgs:"int age custom tttt"`
//		Text             string   `json:"text" validate:"required" errormgs:"ssss age custom tttt333"`
//	}
//
//	func (m *PaymentInfo) Validate(validate *validator.Validate) error {
//		return ValidateFunc[PaymentInfo](*m, validate)
//	}
//
//	type Message struct {
//		Text string `json:"text" validate:"required" errormgs:"int age custom tttt333"`
//	}
//
// func main() {
//
//	pm := PaymentInfo{
//		CreditCardNumber: "",
//		CVV:              "23123",
//		Age:              nil,
//		Message: &Message{
//			Text: "",
//		},
//	}
//
//	// Create a new validator instance
//	validate := validator.New()
//
//	err := pm.Validate(validate)
//	fmt.Println("##################1")
//	fmt.Println(err)
//	fmt.Println("##################2")
//	err = validate.Struct(pm)
//	fmt.Println("##################3")
//	fmt.Println(err)
//	fmt.Println("##################4")
//
// }
//package main
//
//import (
//	"fmt"
//
//	"github.com/go-playground/validator/v10"
//)
//
//// User struct to validate
//type User struct {
//	Username string `json:"username" validate:"required"`
//	Email    string `json:"email" validate:"required,email"`
//	Age      uint8  `json:"age" validate:"gte=18"` // must be 18 or older
//}
//
//// validateUser function takes a User struct and returns a map[string]string containing validation errors
//func validateUser(user User) map[string]string {
//	validate := validator.New()
//	err := validate.Struct(user)
//
//	if err != nil {
//		validationErrors := err.(validator.ValidationErrors)
//		errorMap := make(map[string]string)
//		for _, fieldError := range validationErrors {
//			errorMap[fieldError.Field()] = fieldError.Error()
//			// Placeholder for translated message (future implementation)
//			errorMap[fieldError.Field()+"_translated"] = translate(fieldError.Error())
//		}
//		return errorMap
//	}
//	return nil
//}
//
//// Placeholder function for translation (replace with actual translation logic)
//func translate(message string) string {
//	// Replace this with your translation logic using an external library or API
//	return fmt.Sprintf("Translation not available yet: %s", message)
//}
//
//func main() {
//	// Create a User instance with some invalid data
//	user := User{
//		Username: "", // missing username
//		Email:    "notanemail",
//		Age:      15, // underage
//	}
//
//	// Validate the user struct
//	validationErrors := validateUser(user)
//
//	if validationErrors != nil {
//		// Print validation errors with placeholders for translated messages
//		for field, message := range validationErrors {
//			fmt.Printf("Error: %s - %s (Translation: %s)\n", field, message, translate(message))
//		}
//	} else {
//		fmt.Println("User data is valid!")
//	}
//}

// #####################
//package main
//
//import (
//	"fmt"
//	"github.com/go-playground/locales/ar"
//	en_translations "github.com/go-playground/validator/v10/translations/en"
//
//	"github.com/go-playground/locales/en"
//	ut "github.com/go-playground/universal-translator"
//	"github.com/go-playground/validator/v10"
//)
//
//// User contains user information
//type User struct {
//	FirstName      string     `validate:"required"`
//	LastName       string     `validate:"required"`
//	Age            uint8      `validate:"gte=0,lte=130"`
//	Email          string     `validate:"required,email"`
//	FavouriteColor string     `validate:"iscolor"`                // alias for 'hexcolor|rgb|rgba|hsl|hsla'
//	Addresses      []*Address `validate:"required,dive,required"` // a person can have a home and cottage...
//}
//
//// Address houses a users address information
//type Address struct {
//	Street string `validate:"required"`
//	City   string `validate:"required"`
//	Planet string `validate:"required"`
//	Phone  string `validate:"required"`
//}
//
//// use a single instance , it caches struct info
//var (
//	uni      *ut.UniversalTranslator
//	validate *validator.Validate
//)
//
//func main() {
//
//	// NOTE: ommitting allot of error checking for brevity
//
//	en := en.New()
//	ar := ar.New()
//	uni = ut.New(en, en, ar)
//
//	// this is usually know or extracted from http 'Accept-Language' header
//	// also see uni.FindTranslator(...)
//	trans, _ := uni.GetTranslator("en")
//
//	validate = validator.New()
//	//en_translations.RegisterDefaultTranslations(validate, trans)
//	en_translations.RegisterDefaultTranslations(validate, trans)
//	translateAll(trans)
//	//translateIndividual(trans)
//	//translateOverride(trans) // yep you can specify your own in whatever locale you want!
//}
//
//func translateAll(trans ut.Translator) {
//
//	type User struct {
//		Username string     `validate:"required"`
//		Tagline  string     `validate:"required,lt=10"`
//		Tagline2 string     `validate:"required,gt=1"`
//		t        []*Address `json:"t" validate:"omitempty,dive"`
//	}
//
//	adds := make([]*Address, 0)
//	user := User{
//		Username: "",
//		Tagline:  "This tagline is way too long.",
//		Tagline2: "1",
//		t:        adds,
//	}
//
//	err := validate.Struct(user)
//	if err != nil {
//
//		// translate all error at once
//		errs := err.(validator.ValidationErrors)
//
//		// returns a map with key = namespace & value = translated error
//		// NOTICE: 2 errors are returned and you'll see something surprising
//		// translations are i18n aware!!!!
//		// eg. '10 characters' vs '1 character'
//		fmt.Println(errs.Translate(trans))
//	}
//}
//
//func translateIndividual(trans ut.Translator) {
//
//	type User struct {
//		Username string     `validate:"required"`
//		Tagline  string     `validate:"required,lt=10"`
//		Tagline2 string     `validate:"required,gt=1"`
//		t        []*Address `json:"t" validate:"omitempty,dive"`
//	}
//
//	adds := make([]*Address, 0)
//	user := User{
//		Username: "",
//		Tagline:  "This tagline is way too long.",
//		Tagline2: "1",
//		t:        adds,
//	}
//
//	err := validate.Struct(user)
//	if err != nil {
//
//		errs := err.(validator.ValidationErrors)
//
//		for _, e := range errs {
//			// can translate each error one at a time.
//			fmt.Println(e.Translate(trans))
//		}
//	}
//}
//
//func translateOverride(trans ut.Translator) {
//
//	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
//		return ut.Add("required", "{0} must have a value!", true) // see universal-translator for details
//	}, func(ut ut.Translator, fe validator.FieldError) string {
//		t, _ := ut.T("required", fe.Field())
//
//		return t
//	})
//
//	type User struct {
//		Username string `validate:"required"`
//	}
//
//	var user User
//
//	err := validate.Struct(user)
//	if err != nil {
//
//		errs := err.(validator.ValidationErrors)
//
//		for _, e := range errs {
//			// can translate each error one at a time.
//			fmt.Println(e.Translate(trans))
//		}
//	}
//}
