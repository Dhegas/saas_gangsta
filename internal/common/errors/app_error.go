package errors

// AppError adalah error terstruktur yang digunakan di seluruh aplikasi.
// Field Details sengaja TIDAK di-export ke JSON untuk mencegah kebocoran data internal.
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

// New membuat AppError baru.
// Parameter details TIDAK diteruskan ke client; gunakan hanya untuk logging internal jika diperlukan.
func New(code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}
