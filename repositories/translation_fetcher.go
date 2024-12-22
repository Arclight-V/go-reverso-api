package repositories

import (
	"github.com/marycka9/go-reverso-api/entities"
	"github.com/marycka9/go-reverso-api/languages"
)

// TranslationFetcher defines an interface for fetching translations from a source
type TranslationFetcher interface {
	FetchTranslations(term, partOfSpeech string, srcLang, dstLang *languages.Language) ([]string, error)
	FetchTranscription(term string, srcLang, dstLang entities.Language) (string, error)
	FetchAdditionalData(word *entities.Word) error
}
