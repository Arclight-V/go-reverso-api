package repositories

import (
	"errors"
	"github.com/gocolly/colly"
	"github.com/marycka9/go-reverso-api/common"
	"github.com/marycka9/go-reverso-api/entities"
	"github.com/marycka9/go-reverso-api/languages"
	log "github.com/sirupsen/logrus"
	"strings"
)

const (
	baseUrlCambridge = "https://dictionary.cambridge.org/dictionary/"
)

type DictionaryCambridgeParser struct{}

func NewDictionaryCambridgeParser() *DictionaryCambridgeParser {
	return &DictionaryCambridgeParser{}
}

// Extract the translation from French to English
func fetchTranslationFromFrenchToEnglis(term, partOfSpeech string) ([]string, error) {
	c := colly.NewCollector()
	c.UserAgent = DefaultUserAgent

	var firstMatchFound bool
	parser := common.GetPartOfSpeechParserInstance()
	var translations []string
	var transcription string
	c.OnHTML("div.pr.dictionary", func(e *colly.HTMLElement) {
		if firstMatchFound {
			return
		}
		findPartOfSpeech := e.DOM.Find("span.pos.dpos").First().Text()
		if parser.Parse(findPartOfSpeech) != partOfSpeech {
			return
		}

		//transcription = e.DOM.Find("span.pron.dpron").First().Text()

		// TODO: When you need several translations of one word, then change here
		translation := e.DOM.Find("span.trans.dtrans").First().Text()
		translations = append(translations, translation)
		log.Infof("partOfSpeach: %s, translation: %s , transcription: %s", findPartOfSpeech, translation, transcription)

		// Set flag, for stopping process
		firstMatchFound = true
	})

	var builder strings.Builder
	builder.WriteString(baseUrlCambridge)
	builder.WriteString(string(entities.French))
	builder.WriteRune('-')
	builder.WriteString(string(entities.English))
	builder.WriteRune('/')
	builder.WriteString(term)

	err := c.Visit(builder.String())

	return translations, err

}

func fetchTranscription(term string, srcLang, dstLang entities.Language) (string, error) {
	c := colly.NewCollector()
	c.UserAgent = DefaultUserAgent

	var firstMatchFound bool
	var transcription string

	c.OnHTML("div.pr.dictionary", func(e *colly.HTMLElement) {
		if firstMatchFound {
			return
		}
		transcription = e.DOM.Find("span.pron.dpron").First().Text()
		// Set flag, for stopping process
		firstMatchFound = true
	})

	var builder strings.Builder
	builder.WriteString(baseUrlCambridge)
	builder.WriteString(string(srcLang))
	builder.WriteRune('-')
	builder.WriteString(string(dstLang))
	builder.WriteRune('/')
	builder.WriteString(term)

	err := c.Visit(builder.String())

	return transcription, err

}

// Don't use this metod TODO:: refactoring
func (p *DictionaryCambridgeParser) FetchTranslations(term, partOfSpeech string, srcLang, dstLang *languages.Language) ([]string, error) {
	if term == "" {
		return nil, errors.New("term cannot be empty")
	}
	if partOfSpeech == "" {
		return nil, errors.New("part_of_speech cannot be empty")
	}

	var translations []string
	var err error
	//if srcLang == entities.French {
	translations, err = fetchTranslationFromFrenchToEnglis(term, partOfSpeech)
	if err != nil {
		log.Error("Failed to visit URL:", err)
	}
	//}

	return translations, nil
}

func (p *DictionaryCambridgeParser) FetchTranscription(term string, srcLang, dstLang entities.Language) (string, error) {
	if term == "" {
		return "", errors.New("term cannot be empty")
	}
	if srcLang == "" {
		return "", errors.New("src_lang cannot be empty")
	}
	if dstLang == "" {
		return "", errors.New("dst_lang cannot be empty")
	}
	transcription, err := fetchTranscription(term, srcLang, dstLang)
	return transcription, err
}

func (p *DictionaryCambridgeParser) FetchAdditionalData(word *entities.Word) error {
	return nil
}
