package api

// Error codes for loginError
// You can retrieve them on returned errors.ErrorCode()
const (
	// Every Login error that could not be categorized
	LoginErrorUnknown uint32 = 600 + iota

	// No user found with the given username
	LoginErrorNoUsername

	// Multiple users found with the given username
	LoginErrorMultipleUsername

	// Trying to create a new user with a already taken username
	LoginErrorAlreadyTakenUsername

	// Trying to authenticate but the wrong password was given
	LoginErrorWrongPassword

	// The given token has expired, you should request a new one
	ExpiredTokenError

	// The JWT is not valid (modified by the user?)
	InvalidTokenError

	// The signing method for the token was not the one expected (modified?)
	InvalidTokenSigningMethodError

	// The JWT cannot be parsed
	UnreadableTokenError

	// The signing for the token provoked an error.
	TokenSigningError

	// The used signing key is invalid and could not be corrected. Used as a
	// security measure, this error will not appear unless the code has been
	// wrongly modified to a point where we cannot guarantee the security
	// of our tokens
	InvalidSigningKeyError
)

// loginError defines Errors happening on login (not from database/crypting
// errors)
type loginError struct {
	err string

	// Error code. See constants.
	code uint32
}

func (lerr loginError) ErrorCode() uint32 {
	return lerr.code
}

func (lerr loginError) Error() string {
	if lerr.err != "" {
		return lerr.err
	}
	return "Login failed"
}

type expiredTokenError struct{}

type invalidTokenSigningMethodError struct {
	alg string
}

type unreadableTokenError struct {
	field string
}

type invalidTokenError struct{}

type noUserFoundError struct {
	username string
}

type multipleUserFound struct {
	username string
}

type alreadyCreatedUserError struct {
	username string
}

func (err expiredTokenError) Error() string {
	return "This token has expired"
}

func (err invalidTokenSigningMethodError) Error() string {
	if err.alg != "" {
		return "Unexpected signing method: " + err.alg
	}
	return "This token has not the right signing method"
}

func (err unreadableTokenError) Error() string {
	if err.field != "" {
		return "The following field could not be parsed: " + err.field
	}
	return "This token could not be parsed"
}

func (err invalidTokenError) Error() string {
	return "This token is invalid"
}

func (err multipleUserFound) Error() string {
	if err.username != "" {
		return "Multiple users found with that username: " + err.username
	}
	return "Multiple users found"
}

func (err noUserFoundError) Error() string {
	if err.username != "" {
		return "No user found with that username: " + err.username
	}
	return "No user found"
}

func (err alreadyCreatedUserError) Error() string {
	return "An user with that username already exists: " + err.username
}
