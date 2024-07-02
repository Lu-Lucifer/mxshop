package code

//go:generate codegen -type=int
const (
	// ErrUserNotFound - 404: User not found.
	ErrUserNotFound int = iota + 100401
	// ErrUserAlreadyExists - 400: User already exists.
	ErrUserAlreadyExists
	// ErrPasswordIncorrect - 400: Password is incorrect.
	ErrPasswordIncorrect
	// ErrSmsSend - 400: Send sms error.
	ErrSmsSend
	// ErrCodeNotExist - 404: Sms code incorrect or expired.
	ErrCodeNotExist
	// ErrCodeIncorrect - 400: Sms code incorrect.
	ErrCodeIncorrect
)
