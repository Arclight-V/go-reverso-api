package entities

type Language string

const (
	French  Language = "french"
	English Language = "english"
	Russian Language = "russian"
)

type Translations = map[Language][]string

type Word struct {
	Language      Language
	Term          string
	TermAlt       string
	PartOfSpeech  string
	Transcription string
	Translations  Translations
	// TODO:: merge with PartOfSpeech
	Type string
}

// Conjugation for verb, before adding new tenses, create a type
type FrenchVerbConjugation struct {
	Infinitif string
	Indicatif map[string][]string
	Imperatif map[string][]string
}
