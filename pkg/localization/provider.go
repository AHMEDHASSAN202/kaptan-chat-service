package localization

import (
	"bytes"
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"io/ioutil"
	"sort"
	"text/template"
)

var bundle *i18n.Bundle
var path = "pkg/localization"
var enTomlPath = path + "/en.toml"
var arTomlPath = path + "/ar.toml"

func InitLocalization() *i18n.Bundle {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.LoadMessageFile(enTomlPath)
	bundle.LoadMessageFile(arTomlPath)
	GenerateErrorCodeStruct()
	return bundle
}

func GetTranslation(c *context.Context, errorCode string, TemplateData map[string]interface{}, lang string) string {
	loc := i18n.NewLocalizer(bundle, getLangParam(c, lang))
	translation, err := loc.Localize(&i18n.LocalizeConfig{
		MessageID:    errorCode,
		TemplateData: TemplateData,
	})
	if err != nil {
		fmt.Println(err)
		translation = "error_msg_not_found"
	}
	return translation
}

func getLangParam(c *context.Context, lang string) string {
	if lang == "" {
		langCtx, ok := (*c).Value("lang").(string)
		if !ok {
			lang = "en" //default value
		} else {
			lang = langCtx
		}
	}
	return lang
}

func GetAttrByLang(ctx *context.Context, ifEn, ifAr string) string {
	lang := getLangParam(ctx, "")
	if lang == "en" {
		return ifEn
	} else if lang == "ar" {
		return ifAr
	}
	return ""
}

// ErrorCode represents a single error code with its ID and description
type ErrorCode struct {
	ID          string
	Description string
}

// Config represents the overall structure of the TOML file
type Config map[string]struct {
	Description string `toml:"description"`
	One         string `toml:"one"`
}

func GenerateErrorCodeStruct() {
	// Read and parse the TOML file
	var config Config
	if _, err := toml.DecodeFile(enTomlPath, &config); err != nil {
		fmt.Printf("Failed to parse TOML file: %v", err)
	}

	var errorCodes []ErrorCode
	for key, value := range config {
		errorCodes = append(errorCodes, ErrorCode{ID: key, Description: value.Description})
	}
	sort.Slice(errorCodes, func(i, j int) bool {
		return errorCodes[i].ID < errorCodes[j].ID
	})

	// Generate the struct definition
	structDef := generateEnumDefinition(errorCodes)

	// Write the struct definition to a Go file
	if err := ioutil.WriteFile(path+"/error_codes.go", []byte(structDef), 0644); err != nil {
		fmt.Printf("Failed to write struct definition to file: %v", err)
	}

	fmt.Println("Struct definition generated and written to error_codes.go")
}

// generateEnumDefinition generates a Go enum definition for the given error codes
func generateEnumDefinition(errorCodes []ErrorCode) string {
	const enumTemplate = `package localization

// ErrorCode is an enum representing all message IDs

const (
{{- range . }}
    {{ .ID }} = "{{ .ID }}" // {{ .Description }}
{{- end }}
)
`
	tmpl, err := template.New("enum").Parse(enumTemplate)
	if err != nil {
		fmt.Printf("Failed to parse template: %v", err)
	}

	var result bytes.Buffer
	tmpl.Execute(&result, errorCodes)

	return result.String()
}
