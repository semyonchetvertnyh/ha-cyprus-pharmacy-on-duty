package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func Translate(apiKey, text string) (string, error) {
	url := "https://translation.googleapis.com/language/translate/v2?key=" + apiKey

	reqBody := map[string]interface{}{
		"q":      text,
		"source": "el",
		"target": "en",
		"format": "text",
	}

	body, _ := json.Marshal(reqBody)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Translations []struct {
				TranslatedText string `json:"translatedText"`
			} `json:"translations"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Data.Translations[0].TranslatedText, nil
}
