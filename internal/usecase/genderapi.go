package usecase

import (
	"encoding/json"
	"fmt"
	"github.com/linqcod/user-service/internal/domain"
	"github.com/spf13/viper"
	"net/http"
)

type genderApiUsecase struct {
	uri string
}

func NewGenderApiUsecase() domain.GenderApiUsecase {
	return &genderApiUsecase{
		uri: viper.GetString("GENDER_API_URI"),
	}
}

func (u genderApiUsecase) Get(name string) (string, error) {
	resp, err := http.Get(u.uri + name)
	if err != nil {
		return "", fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}

	var gender domain.Gender
	err = json.NewDecoder(resp.Body).Decode(&gender)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return gender.Value, nil
}
