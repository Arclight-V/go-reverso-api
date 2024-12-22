package usecases

import (
	"errors"
	"github.com/marycka9/go-reverso-api/entities"
	"github.com/marycka9/go-reverso-api/languages"
	"github.com/marycka9/go-reverso-api/repositories"
	log "github.com/sirupsen/logrus"
)

// TranslationService manages fetching translations from various sources
type TranslationService struct {
	fetchers []repositories.TranslationFetcher
}

// NewTranslationService creates a new TranslationService
func NewTranslationService(fetchers []repositories.TranslationFetcher) *TranslationService {
	return &TranslationService{fetchers: fetchers}
}

// GetTranslations fetches translations from all available sources
func (s *TranslationService) GetTranslations(word *entities.Word, srcLang, dstLang *languages.Language) error {
	if word.Term == "" {
		return errors.New("term cannot be empty")
	}
	// If there are new sources of transfers, then change 0 to the corresponding identifier.
	translations, err := s.fetchers[0].FetchTranslations(word.Term, word.PartOfSpeech, srcLang, dstLang)
	if err != nil {
		return err
	}
	word.Translations[entities.Language(dstLang.Code)] = append(word.Translations[entities.Language(dstLang.Code)], translations...)

	return err
}

func (s *TranslationService) GetTranscriptions(word *entities.Word, srcLang, dstLang entities.Language) error {
	if word.Term == "" {
		return errors.New("term cannot be empty")
	}
	transcription, err := s.fetchers[1].FetchTranscription(word.Term, srcLang, dstLang)
	if err != nil {
		return err
	}
	log.Info(transcription)
	word.Transcription = transcription
	return nil
}
