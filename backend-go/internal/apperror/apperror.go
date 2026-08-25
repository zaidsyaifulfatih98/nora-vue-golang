package apperror

// AppError mirrors backend/src/utils/app-error.utils.ts: an error carrying an
// HTTP status code whose message is safe to expose to the client.
type AppError struct {
	Message    string
	StatusCode int
	Expose     bool
}

func (e *AppError) Error() string {
	return e.Message
}

// New creates an exposable AppError (status defaults to 500 if 0 is passed).
func New(message string, statusCode int) *AppError {
	if statusCode == 0 {
		statusCode = 500
	}
	return &AppError{Message: message, StatusCode: statusCode, Expose: true}
}
