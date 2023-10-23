package usecase

import (
	"encoding/json"
	"fmt"
	"github.com/linqcod/user-service/internal/domain"
	"github.com/spf13/viper"
	"net/http"
)

type nationalityApiUsecase struct {
	uri string
}

func NewNationalityApiUsecase() domain.NationalityApiUsecase {
	return &nationalityApiUsecase{
		uri: viper.GetString("NATIONALITY_API_URI"),
	}
}

func (u nationalityApiUsecase) Get(name string) (string, error) {
	resp, err := http.Get(u.uri + name)
	if err != nil {
		return "", fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}

	var nationality domain.Nationality
	err = json.NewDecoder(resp.Body).Decode(&nationality)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return nationality.Countries[0].CountryId, nil
}
