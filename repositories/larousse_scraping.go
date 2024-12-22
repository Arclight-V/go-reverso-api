package repositories

import (
	"github.com/marycka9/go-reverso-api/entities"
	"github.com/marycka9/go-reverso-api/languages"
	"github.com/serope/laroussefr/traduction"
	log "github.com/sirupsen/logrus"
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
		log.Errorf("Error traduction: %v", err)
	}
	log.Info(result.Words)
	return nil
}
