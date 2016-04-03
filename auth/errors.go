package auth

// Error codes for authentication errors
// You can retrieve them on returned errors.ErrorCode()
const (
	// Every Login error that could not be categorized
	UnknownAuthenticationErrorCode uint32 = 300 + iota

	// No user found with the given username
	UserNotFoundErrorCode

	// Multiple users found with the given username
	MultipleUserFoundErrorCode

	// Trying to create a new user with a already taken username
	AlreadyTakenUsernameErrorCode

	// Trying to authenticate but the wrong password was given
	WrongPasswordErrorCode

	// The given token has expired, you should request a new one
	ExpiredTokenErrorCode

	// The JWT is not valid (modified by the user?)
	InvalidTokenErrorCode

	// No jwt was provided
	NoTokenErrorCode

	// The signing method for the token was not the one expected (modified?)
	InvalidTokenSigningMethodErrorCode

	// The JWT cannot be parsed
	UnreadableTokenErrorCode

	// The token could not be signed with our signing key.
	TokenSigningErrorCode

	// The used signing key is invalid and could not be corrected. Used as a
	// security measure, this error will not appear unless the code has been
	// wrongly modified to a point where we cannot guarantee the security
	// of our tokens
	InvalidSigningKeyErrorCode
)

type AuthenticationError interface {
	error
	ErrorCode() uint32
}

type genericAuthenticationError struct {
	// Error message
	err string

	// Error code. See constants.
	code uint32
}

func (lerr genericAuthenticationError) Error() string {
	if lerr.err != "" {
		return lerr.err
	}
	return "Authentication failed."
}

func (lerr genericAuthenticationError) ErrorCode() uint32 {
	if lerr.code == 0 {
		return UnknownAuthenticationErrorCode
	}
	return lerr.code
}

type userNotFoundError struct{ username string }
type multipleUserFoundError struct{ username string }
type alreadyTakenUsernameError struct{ username string }
type wrongPasswordError struct{ username string }
type expiredTokenError struct{}
type invalidTokenError struct{}
type noTokenError struct{}
type invalidTokenSigningMethodError struct{ alg string }
type unreadableTokenError struct{ field string }
type tokenSigningError struct{}
type invalidSigningKeyError struct{ field string }

func (err userNotFoundError) Error() string {
	if err.username != "" {
		return "No user found with that username: " + err.username + "."
	}
	return "No user found."
}

func (err userNotFoundError) ErrorCode() uint32 {
	return UserNotFoundErrorCode
}

func (err multipleUserFoundError) Error() string {
	if err.username != "" {
		return "Multiple users found with that username: " + err.username + "."
	}
	return "Multiple users found."
}

func (err multipleUserFoundError) ErrorCode() uint32 {
	return MultipleUserFoundErrorCode
}

func (err alreadyTakenUsernameError) Error() string {
	return "An user with that username already exists: " + err.username + "."
}

func (err alreadyTakenUsernameError) ErrorCode() uint32 {
	return AlreadyTakenUsernameErrorCode
}

func (err wrongPasswordError) Error() string {
	if err.username != "" {
		return "Wrong password for that username: " + err.username + "."
	}
	return "Wrong password."
}

func (err wrongPasswordError) ErrorCode() uint32 {
	return WrongPasswordErrorCode
}

func (err expiredTokenError) Error() string {
	return "This token has expired."
}

func (err expiredTokenError) ErrorCode() uint32 {
	return ExpiredTokenErrorCode
}

func (err invalidTokenError) Error() string {
	return "This token is invalid."
}

func (err invalidTokenError) ErrorCode() uint32 {
	return InvalidTokenErrorCode
}

func (err noTokenError) Error() string {
	return "No token provided."
}

func (err noTokenError) ErrorCode() uint32 {
	return NoTokenErrorCode
}

func (err invalidTokenSigningMethodError) Error() string {
	if err.alg != "" {
		return "Unexpected signing method: " + err.alg + "."
	}
	return "This token has not the right signing method."
}

func (err invalidTokenSigningMethodError) ErrorCode() uint32 {
	return InvalidTokenSigningMethodErrorCode
}

func (err unreadableTokenError) Error() string {
	if err.field != "" {
		return "The following field could not be parsed: " + err.field + "."
	}
	return "This token could not be parsed."
}

func (err unreadableTokenError) ErrorCode() uint32 {
	return UnreadableTokenErrorCode
}

func (err tokenSigningError) Error() string {
	return "This token could not be signed."
}

func (err tokenSigningError) ErrorCode() uint32 {
	return TokenSigningErrorCode
}

func (err invalidSigningKeyError) Error() string {
	return "The signing key for JWT is not secure enough."
}

func (err invalidSigningKeyError) ErrorCode() uint32 {
	return InvalidSigningKeyErrorCode
}
