package localization

import (
	"context"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func InitLocalization() *i18n.Bundle {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.LoadMessageFile("pkg/localization/en.toml")
	bundle.LoadMessageFile("pkg/localization/ar.toml")
	return bundle
}

func GetTranslation(c *context.Context, errorCode string, TemplateData map[string]interface{}) string {
	lang := (*c).Value("lang").(string)
	loc := i18n.NewLocalizer(bundle, lang)
	translation, err := loc.Localize(&i18n.LocalizeConfig{
		MessageID:    errorCode,
		TemplateData: TemplateData,
		PluralCount:  0,
	})
	if err != nil {
		translation = "error_msg_not_found"
	}
	return translation
}
