package errors

const AuthErrorMsg = "Auth failed"

type AuthError struct {
	cause error
}

func NewAuthError(cause error) AuthError {
	return AuthError{cause: cause}
}

func (r AuthError) Error() string {
	return AuthErrorMsg + ":" + r.cause.Error()
}
