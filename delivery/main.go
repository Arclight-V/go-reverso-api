package main

import (
	"flag"
	"fmt"
	"github.com/atselvan/ankiconnect"
	"github.com/marycka9/go-reverso-api/client"
	"github.com/marycka9/go-reverso-api/entities"
	"github.com/marycka9/go-reverso-api/languages"
	"github.com/marycka9/go-reverso-api/repositories"
	"github.com/marycka9/go-reverso-api/usecases"
	log "github.com/sirupsen/logrus"
	"strings"
)

//func InsertSpaceBeforeSymbol(s string, target rune) string {
//	runes := []rune(s)
//	var result []rune
//
//	for i, r := range runes {
//		if r == target {
//			// Проверяем, не является ли предыдущий символ пробелом
//			if i > 0 && runes[i-1] != ' ' {
//				result = append(result, ' ')
//			}
//		}
//		result = append(result, r)
//	}
//
//	return string(result)
//}

func main() {
	logger := log.New()

	// Flags for CSV file paths
	frenchFilePath := flag.String("french", "", "Path to the French CSV file")
	englishFilePath := flag.String("english", "", "Path to the English CSV file")
	russianFilePath := flag.String("russian", "", "Path to the Russian CSV file")
	flag.Parse()

	// Checking for mandatory flags
	if *frenchFilePath == "" || *englishFilePath == "" || *russianFilePath == "" {
		logger.Error("Error: all file paths must be provided")
		flag.Usage()
		return
	}
	// Repositories
	csvRepo := repositories.NewCSVRepository()

	// Read data from CSV files
	frenchWords, err := csvRepo.ReadWordsFromFile(*frenchFilePath, entities.French)
	if err != nil {
		logger.Fatal("Error reading French words:", err)
		return
	}

	englishWords, err := csvRepo.ReadWordsFromFile(*englishFilePath, entities.English)
	if err != nil {
		logger.Fatal("Error reading English words:", err)
		return
	}

	russianWords, err := csvRepo.ReadWordsFromFile(*russianFilePath, entities.Russian)
	if err != nil {
		logger.Fatal("Error reading Russian words:", err)
		return
	}

	// UseCases
	wordTranslator := usecases.NewWordTranslator()

	// Aggregate all words by language
	wordsByLanguage := map[entities.Language][]entities.Word{
		entities.French:  frenchWords,
		entities.English: englishWords,
		entities.Russian: russianWords,
	}

	// Translate words between languages
	translatedWords := wordTranslator.TranslateWords(wordsByLanguage)

	// Initialize clients
	reversoContextClient := client.NewClient()
	dictionaryCambridgeParser := repositories.NewDictionaryCambridgeParser()
	larousseScarper := repositories.NewLarousseScarping()

	// Register parsers in the service
	translationService := usecases.NewTranslationService([]repositories.TranslationFetcher{
		reversoContextClient,
		dictionaryCambridgeParser,
		larousseScarper,
	})

	langs := languages.GetLanguages()
	// Display the translated words
	for _, word := range translatedWords {
		if word.Language == entities.French {
			// For three-way cards
			//err = translationService.GetTranslations(&word, langs[string(entities.French)], langs[string(entities.English)])
			//err = translationService.GetTranscriptions(&word, word.Language, entities.English)
			//if err != nil {
			//	log.Errorf("Error GetTranscriptions", err)
			//	continue
			//}
			err = larousseScarper.FetchAdditionalData(&word)
			if err != nil {
				log.Errorf("Error FetchAdditionalData", err)
				continue
			}
			err = translationService.GetTranslations(&word, langs[string(entities.French)], langs[string(entities.Russian)])
			if err != nil {
				log.Errorf("Error GetTranslations", err)
				continue
			}
			ankiClient := ankiconnect.NewClient()
			if word.PartOfSpeech == "v" {
				verb, err := reversoContextClient.FetchConjugation(word.Term, word.Language)
				if err != nil {
					log.Errorf("Error FetchConjugation", err)
					continue
				}
				log.Infof("Conjugation: %s", verb)
				note := ankiconnect.Note{
					DeckName:  "Francais_conjugation",
					ModelName: "Basic (de conjugaison A1)",
					Fields: ankiconnect.Fields{
						"Infinitif": verb.Infinitif,
						"Présent":   strings.Join(verb.Indicatif["Présent"], "<br>"),
						"Impératif": strings.Join(verb.Imperatif["Présent"], "<br>"),
					},
				}
				restErr := ankiClient.Notes.Add(note)
				if restErr != nil {
					log.Error(restErr)
				}
			}
			if strings.IndexRune(word.Transcription, ',') != -1 && word.TermAlt == "" {
				note := ankiconnect.Note{
					DeckName:  "Francais_mots_corriger",
					ModelName: "Basic (and reversed card french)",
					Fields: ankiconnect.Fields{
						"Front": strings.Join([]string{fmt.Sprintf("%s %s", word.Term, "ERROR"), word.Transcription, word.Type}, "<br>"),
						"Back":  strings.Join(word.Translations["ru"], "<br>"),
					},
				}
				restErr := ankiClient.Notes.Add(note)
				if restErr != nil {
					log.Error(restErr)
				}
			}
			note := ankiconnect.Note{
				DeckName:  "Francais_mots",
				ModelName: "Basic (and reversed card french)",
				Fields: ankiconnect.Fields{
					"Front": strings.Join([]string{fmt.Sprintf("%s %s", word.Term, word.TermAlt), word.Transcription, word.Type}, "<br>"),
					"Back":  strings.Join(word.Translations["ru"], "<br>"),
				},
			}
			restErr := ankiClient.Notes.Add(note)
			if restErr != nil {
				log.Error(restErr)
			}

		} else {
			// For three-way cards
			//err = translationService.GetTranslations(&word, langs[string(entities.English)], langs[string(entities.French)])
			//if err != nil {
			//	log.Errorf("Error GetTranslations", err)
			//	continue
			//}
			err = translationService.GetTranscriptions(&word, word.Language, entities.Russian)
			if err != nil {
				log.Errorf("Error GetTranscriptions", err)
			}
			err = translationService.GetTranslations(&word, langs[string(entities.English)], langs[string(entities.Russian)])
			if err != nil {
				log.Errorf("Error GetTranslations", err)
				continue
			}
		}
		logger.Infof("[%s] %s %s (%s) %s\n", word.Language, word.Term, word.TermAlt, word.PartOfSpeech, word.Transcription)
		for k, v := range word.Translations {
			log.Infof(" [%s] (%s)\n", k, v)
		}

	}
}
