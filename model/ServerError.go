package model

type BindingErrorMsg struct {
	Field   string `json:"field" validate:"required"`
	Message string `json:"message" validate:"required"`
}

type ServerError struct {
	BindingErrors []*BindingErrorMsg `json:"errors"`
	Message       string             `json:"message"`
}
