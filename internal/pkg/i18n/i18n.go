package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

// Init initializes the i18n bundle and loads language files
func Init() error {
	// Create a new bundle with English as the default language
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Get the project root directory
	rootDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	localesDir := filepath.Join(rootDir, "locales")

	// Load language files
	languages := []string{"en", "zh"}
	for _, lang := range languages {
		filename := filepath.Join(localesDir, fmt.Sprintf("%s.json", lang))
		if _, err := os.Stat(filename); err == nil {
			if _, err := bundle.LoadMessageFile(filename); err != nil {
				return fmt.Errorf("failed to load %s language file: %w", lang, err)
			}
		}
	}

	return nil
}

// GetLocalizer retrieves the localizer from the Gin context
// If not found, returns a localizer with English as default
func GetLocalizer(c *gin.Context) *i18n.Localizer {
	if localizer, exists := c.Get("localizer"); exists {
		if loc, ok := localizer.(*i18n.Localizer); ok {
			return loc
		}
	}
	// Fallback to English if no localizer in context
	return i18n.NewLocalizer(bundle, language.English.String())
}

// T translates a message key with optional template data
func T(c *gin.Context, messageID string, templateData map[string]interface{}) string {
	localizer := GetLocalizer(c)

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})

	if err != nil {
		// If translation fails, return the message ID as fallback
		return messageID
	}

	return msg
}

// TSimple translates a message key without template data
func TSimple(c *gin.Context, messageID string) string {
	return T(c, messageID, nil)
}

// GetBundle returns the i18n bundle (useful for creating custom localizers)
func GetBundle() *i18n.Bundle {
	return bundle
}
