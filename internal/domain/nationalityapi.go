package domain

type Nationality struct {
	Countries []Country `json:"country"`
}

type Country struct {
	CountryId string `json:"country_id"`
}

type NationalityApiUsecase interface {
	Get(name string) (string, error)
}
