package usecase

import (
	"encoding/json"
	"fmt"
	"github.com/linqcod/user-service/internal/domain"
	"github.com/spf13/viper"
	"net/http"
)

type ageApiUsecase struct {
	uri string
}

func NewAgeApiUsecase() domain.AgeApiUsecase {
	return &ageApiUsecase{
		uri: viper.GetString("AGE_API_URI"),
	}
}

func (u ageApiUsecase) Get(name string) (int64, error) {
	resp, err := http.Get(u.uri + name)
	if err != nil {
		return -1, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return -1, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}

	var age domain.Age
	err = json.NewDecoder(resp.Body).Decode(&age)
	if err != nil {
		return -1, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return age.Value, nil
}
