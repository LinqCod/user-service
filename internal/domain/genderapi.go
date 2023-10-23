package domain

type Gender struct {
	Value string `json:"gender"`
}

type GenderApiUsecase interface {
	Get(name string) (string, error)
}
