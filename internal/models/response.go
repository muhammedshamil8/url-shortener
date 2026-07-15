package models

type SuccessResponse struct {
	Status    string `json:"status" example:"success"`
	Data      any    `json:"data"`
	RequestID string `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type ErrorResponse struct {
	Status    string `json:"status" example:"error"`
	Message   string `json:"message" example:"Invalid URL"`
	RequestID string `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}
