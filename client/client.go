package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/marycka9/go-reverso-api/entities"
	"github.com/marycka9/go-reverso-api/languages"
	"github.com/marycka9/go-reverso-api/voices"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Client struct {
	Client *http.Client
}

func NewClient() *Client {
	return &Client{
		Client: http.DefaultClient,
	}
}

func (c *Client) Close() {
	c.Close()
}

func (c *Client) Translate(text string, srcLang, dstLang *languages.Language) (*entities.TranslateResponse, error) {
	translateReq := entities.NewTranslateRequest(text, srcLang, dstLang)
	requestBody, err := translateReq.MarshalJson()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		translateReq.GetUrl(),
		strings.NewReader(requestBody),
	)

	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", entities.UserAgentContextBrowser)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	var translate *entities.TranslateResponse
	if err := json.NewDecoder(resp.Body).Decode(&translate); err != nil {
		return nil, err
	}

	_ = resp.Body.Close()

	return translate, nil
}

func (c *Client) Synonyms(text string, language *languages.Language) (*entities.SynonymsResponse, error) {
	synonymRequest := entities.NewSynonymRequest(text, language)

	req, err := http.NewRequest(
		http.MethodGet,
		synonymRequest.GetUrl(language.Code, text),
		nil,
	)

	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", "")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	var synonym *entities.SynonymsResponse
	if err := json.NewDecoder(resp.Body).Decode(&synonym); err != nil {
		return nil, err
	}

	_ = resp.Body.Close()

	return synonym, nil
}

func (c *Client) AutoComplete(text string, language *languages.Language) (*entities.AutoCompleteResponse, error) {
	autoCompleteRequest := entities.NewAutoCompleteRequest()

	req, err := http.NewRequest(
		http.MethodGet,
		autoCompleteRequest.GetUrl(language.Code, text),
		nil,
	)

	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", "")
	req.Header.Add("x-reverso-origin", "synonymapp")
	req.Header.Add("x-reverso-ui-lang", "en")
	req.Header.Add("authorization", fmt.Sprintf("Basic %s", entities.BearerSynonyms))

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	autocomplete := make(entities.AutoCompleteResponse, 0)
	if err := json.NewDecoder(resp.Body).Decode(&autocomplete); err != nil {
		return nil, err
	}

	_ = resp.Body.Close()

	return &autocomplete, nil
}

func (c *Client) Context(text string, srcLang, dstLang *languages.Language, page int) (*entities.ContextResponse, error) {
	queryReq := entities.NewContextRequest(text, srcLang, dstLang, page)

	req, err := http.NewRequest(
		http.MethodPost,
		queryReq.GetUrl(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("User-Agent", entities.UserAgentContextApp)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	var query *entities.ContextResponse
	if err := json.NewDecoder(resp.Body).Decode(&query); err != nil {
		return nil, err
	}

	_ = resp.Body.Close()

	return query, nil
}

func (c *Client) Suggest(text string, srcLang, dstLang *languages.Language) (*entities.SuggestResponse, error) {
	suggestReq := entities.NewSuggestRequest(text, srcLang, dstLang)

	req, err := http.NewRequest(
		http.MethodGet,
		suggestReq.GetUrl(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", entities.UserAgentContextApp)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	var query *entities.SuggestResponse
	if err := json.NewDecoder(resp.Body).Decode(&query); err != nil {
		return nil, err
	}

	_ = resp.Body.Close()

	return query, nil
}

func (c *Client) Speak(fileName, filePath, text string, mp3BitRate, voiceSpeed int) error {
	speakRequest, err := entities.NewSpeakRequest(fileName, filePath, text, voices.VoiceEnglishFemale, mp3BitRate, voiceSpeed)
	if err != nil {
		return err
	}

	if _, err := os.Stat(speakRequest.FilePath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(speakRequest.FilePath, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	fileOut, err := os.OpenFile(speakRequest.GetPath(), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodGet,
		speakRequest.GetUrl(voices.VoiceEnglishFemale),
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", entities.UserAgentContextApp)

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	var buffer *bytes.Buffer
	var body []byte

	buffer = bytes.NewBuffer(body)

	if resp.ContentLength == -1 {
		_, err = buffer.ReadFrom(resp.Body)
	} else {
		body = make([]byte, resp.ContentLength)
		_, err = io.Copy(buffer, resp.Body)
	}

	body = buffer.Bytes()
	_ = resp.Body.Close()

	if _, err = fileOut.Write(body); err != nil {
		return err
	}

	_ = fileOut.Close()

	return nil
}

func (c *Client) FetchTranslations(term, partOfSpeech string, srcLang, dstLang *languages.Language) ([]string, error) {
	res, err := c.Translate(term, srcLang, dstLang)
	if err != nil {
		return nil, err
	}
	var translations []string
	// TODO :: Transfer it to args as it can be used to determine the level of understanding of the language
	i := 0
	for _, result := range res.ContextResults.Results {
		if strings.Contains(result.PartOfSpeech, partOfSpeech) {
			if i >= 1 {
				break
			}
			translations = append(translations, result.Translation)
			i++
		}
	}
	return translations, nil
}

func (c *Client) FetchTranscription(term string, srcLang, dstLang entities.Language) (string, error) {
	return "", nil
}
func (c *Client) FetchAdditionalData(word *entities.Word) error {
	return nil
}

func (c *Client) FetchConjugation(term string, lang entities.Language) (*entities.FrenchVerbConjugation, error) {
	// Проверка языка (допустим, только французский поддерживается)
	if lang != "french" {
		return nil, fmt.Errorf("язык %s не поддерживается для спряжения", lang)
	}

	// Формируем URL, заменяя "aller" на нужный глагол.
	url := fmt.Sprintf("https://conjugator.reverso.net/conjugation-french-verb-%s.html", strings.ToLower(term))

	// Создаём новый HTTP-запрос.
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Устанавливаем необходимые заголовки.
	req.Header.Add("User-Agent", entities.UserAgentContextBrowser)
	req.Header.Add("Accept-Language", "en-US,en;q=0.9,fr;q=0.8")

	// Выполняем запрос.
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("не удалось получить спряжение: статус код %d", resp.StatusCode)
	}

	// Парсим HTML-ответ с помощью goquery.
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// Создаём структуру для хранения спряжения.
	conjugation := &entities.FrenchVerbConjugation{
		Infinitif: term,
		Indicatif: make(map[string][]string),
		Imperatif: make(map[string][]string),
	}

	// Извлекаем Infinitif
	doc.Find(".word-wrap-row").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".word-wrap-title h4").Text()
		if strings.TrimSpace(title) == "Infinitif" {
			s.Find(".blue-box-wrap.alt-tense ul.wrap-verbs-listing li").Each(func(j int, li *goquery.Selection) {
				verbForm := strings.TrimSpace(li.Find("i.verbtxt").Text())
				if verbForm != "" {
					conjugation.Infinitif = verbForm
				}
			})
		}
	})

	// Извлекаем Indicatif
	doc.Find(".result-block-api").Find(".word-wrap-row").Each(func(i int, s *goquery.Selection) {
		// Проверяем, является ли раздел "Indicatif"
		sectionTitle := s.Find(".word-wrap-title h4").Text()
		if strings.TrimSpace(sectionTitle) == "Indicatif" {
			// Для каждой колонки спряжения внутри Indicatif
			s.Find(".wrap-three-col .blue-box-wrap").Each(func(j int, box *goquery.Selection) {
				// Извлекаем время (например, Présent, Imparfait)
				tense := strings.TrimSpace(box.Find("p").First().Text())
				if tense == "" {
					return
				}

				// Извлекаем формы спряжения
				forms := []string{}
				box.Find("ul.wrap-verbs-listing li").Each(func(k int, li *goquery.Selection) {
					// Извлекаем форму глагола из элемента с классом "verbtxt"
					form := strings.TrimSpace(li.Find("i.verbtxt").Text())
					if form == "" {
						// В некоторых случаях может потребоваться объединение с дополнительными частями
						// Например, для Passé composé: "suis allé"
						// Здесь можно адаптировать парсинг при необходимости
						// Например:
						aux := strings.TrimSpace(li.Find("i.auxgraytxt").Text())
						if aux != "" {
							form = aux + " " + form
						}
					}
					if form != "" {
						forms = append(forms, form)
					}
				})

				if len(forms) > 0 {
					conjugation.Indicatif[tense] = forms
				}
			})
		}
	})

	// Поиск всех секций с классом "word-wrap-row"
	doc.Find(".word-wrap-row").Each(func(i int, s *goquery.Selection) {
		// Поиск заголовка внутри "word-wrap-title"
		title := s.Find(".word-wrap-title h4").Text()
		if strings.TrimSpace(title) == "ImpératifInfinitif" {
			// Внутри секции "Impératif" ищем все "blue-box-wrap"
			s.Find(".wrap-three-col .blue-box-wrap").Each(func(j int, box *goquery.Selection) {
				// Извлекаем время из тега <p>
				tense := strings.TrimSpace(box.Find("p").Text())
				if tense == "" {
					return // Пропускаем, если время не найдено
				}

				// Инициализируем срез для хранения форм
				var forms []string

				// Проходимся по каждому <li> внутри <ul>
				box.Find("ul.wrap-verbs-listing li").Each(func(k int, li *goquery.Selection) {
					// Для "Impératif Passé" формы содержат вспомогательный глагол
					aux := strings.TrimSpace(li.Find("i.auxgraytxt").Text())
					verb := strings.TrimSpace(li.Find("i.verbtxt").Text())

					if aux != "" {
						// Если есть вспомогательный глагол, объединяем его с основной формой
						forms = append(forms, fmt.Sprintf("%s%s", aux, verb))
					} else {
						// Иначе добавляем только основную форму
						forms = append(forms, verb)
					}
				})

				// Добавляем извлеченные данные в структуру
				if len(forms) > 0 {
					if _, ok := conjugation.Imperatif[tense]; !ok {
						conjugation.Imperatif[tense] = forms
					}

				}
			})
		}
	})

	// Проверяем, заполнено ли спряжение
	if len(conjugation.Indicatif) == 0 && conjugation.Infinitif == "" {
		return nil, fmt.Errorf("не удалось извлечь спряжение для глагола %s", term)
	}

	return conjugation, nil
}
