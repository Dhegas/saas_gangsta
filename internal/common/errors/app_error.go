package errors

type AppError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Status  int         `json:"-"`
	Details interface{} `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code, message string, status int, details interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Details: details,
	}
}
