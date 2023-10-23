package domain

type SuccessResponse struct {
	Message string `json:"message"`
}

type SuccessWithDataResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
