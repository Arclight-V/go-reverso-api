package usecases

import (
	"errors"
	"fmt"
	"github.com/marycka9/go-reverso-api/entities"
	"github.com/marycka9/go-reverso-api/languages"
	"github.com/marycka9/go-reverso-api/repositories"
)

// TranslationServiceType represents the type of translation service
type TranslationServiceType int

// Translation Service Type lists the available translation services. These constants are used to select the appropriate
// fetcher in the Translation Service structure. When adding a new translation service, define a new constant in this
// block and add the corresponding fetcher to the fetchers field of the Translation Service structure.
const (
	REVERSO   TranslationServiceType = iota // Reverso Context
	CAMBRIDGE                               // Dictionary Cambridge
	LAROUSSE                                // Larousse
)

// String returns a string representation of the service type
func (t TranslationServiceType) String() string {
	switch t {
	case REVERSO:
		return "Reverso Context"
	case CAMBRIDGE:
		return "Dictionary Cambridge"
	case LAROUSSE:
		return "Larousse"
	default:
		return "Unknown"
	}
}

// TranslationService manages fetching translations from various sources
type TranslationService struct {
	fetchers map[TranslationServiceType]repositories.TranslationFetcher
}

// NewTranslationService creates a new TranslationService
func NewTranslationService(fetchers map[TranslationServiceType]repositories.TranslationFetcher) *TranslationService {
	return &TranslationService{fetchers: fetchers}
}

// GetTranslations fetches translations from all available sources
func (s *TranslationService) GetTranslations(service TranslationServiceType, word *entities.Word, srcLang, dstLang *languages.Language) error {
	if word.Term == "" {
		return errors.New("term cannot be empty")
	}

	if fetcher, ok := s.fetchers[service]; ok {
		// If there are new sources of transfers, then change 0 to the corresponding identifier.
		translations, err := fetcher.FetchTranslations(word.Term, word.PartOfSpeech, srcLang, dstLang)
		if err != nil {
			return err
		}
		word.Translations[entities.Language(dstLang.Code)] = append(word.Translations[entities.Language(dstLang.Code)], translations...)
		return nil
	}

	return errors.New("translation service not found")
}

func (s *TranslationService) GetTranscriptions(service TranslationServiceType, word *entities.Word, srcLang, dstLang entities.Language) error {
	if word.Term == "" {
		return errors.New("term cannot be empty")
	}
	if fetcher, ok := s.fetchers[service]; ok {
		transcription, err := fetcher.FetchTranscription(word.Term, srcLang, dstLang)
		if err != nil {
			return err
		}
		word.Transcription = transcription
		return nil
	}
	return errors.New("transcriptions service not found")
}

func (s *TranslationService) GetAdditionalData(service TranslationServiceType, word *entities.Word) error {
	if fetcher, ok := s.fetchers[service]; ok {
		if err := fetcher.FetchAdditionalData(word); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("AdditionalData service: %s not found", service.String())
}

func (s *TranslationService) GetConjugation(service TranslationServiceType, term string, lang entities.Language) (*entities.FrenchVerbConjugation, error) {
	if fetcher, ok := s.fetchers[service]; ok {
		verbConj, err := fetcher.FetchConjugation(term, lang)
		if err != nil {
			return nil, err
		}
		return verbConj, nil
	}
	return nil, fmt.Errorf("Conjugation service: %s not found", service.String())
}
