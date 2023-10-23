package domain

type Age struct {
	Value int64 `json:"age"`
}

type AgeApiUsecase interface {
	Get(name string) (int64, error)
}
