package repositories

import (
	"github.com/marycka9/go-reverso-api/entities"
	"github.com/marycka9/go-reverso-api/languages"
	"github.com/serope/laroussefr/traduction"
)

type LarousseScarping struct{}

func NewLarousseScarping() *LarousseScarping {
	return &LarousseScarping{}
}

// TODO:: refactoring
func (p *LarousseScarping) FetchTranslations(term, partOfSpeech string, srcLang, dstLang *languages.Language) ([]string, error) {
	return nil, nil
}

func (p *LarousseScarping) FetchTranscription(term string, srcLang, dstLang entities.Language) (string, error) {
	return "", nil
}

func (p *LarousseScarping) FetchAdditionalData(word *entities.Word) error {

	result, err := traduction.New(word.Term, traduction.Fr, traduction.En)
	if err != nil {
		return err
	}
	if len(result.Words) > 0 {
		// for feminine and masculine gender
		if result.Words[0].Header.Text != word.Term {
			word.Term = result.Words[0].Header.Text
		}
		if result.Words[0].Header.TextAlt != "" {
			word.TermAlt = result.Words[0].Header.TextAlt
		}
		if result.Words[0].Header.Phonetic != "" {
			word.Transcription = result.Words[0].Header.Phonetic
		}
	}

	return nil
}

func (p *LarousseScarping) FetchConjugation(term string, lang entities.Language) (*entities.FrenchVerbConjugation, error) {
	return nil, nil
}
