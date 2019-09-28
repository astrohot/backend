package zodiac

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var apiSign = map[Sign]string{
	Capricorn:   "capricorn",
	Aquarius:    "aquarius",
	Pisces:      "pisces",
	Taurus:      "taurus",
	Gemini:      "gemini",
	Cancer:      "cancer",
	Leo:         "leo",
	Virgo:       "virgo",
	Libra:       "libra",
	Scorpio:     "scoprio",
	Sagittarius: "sagittarius",
}

// API constants.
const (
	apiURL        = "https://aztro.sameerkumar.website/?sign=%s&day=%s"
	apiDay        = "today"
	responseField = "description"
)

// GetHoroscope ...
func GetHoroscope(sign Sign) (horoscope string, err error) {
	url := fmt.Sprintf(apiURL, apiSign[sign], apiDay)
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return
	}

	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return
	}

	horoscope = result[responseField].(string)
	return
}
