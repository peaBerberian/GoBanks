package login

// Error codes for LoginError
// You can retrieve them on any LoginError{}.ErrorCode
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

	// No user found with the given token
	LoginErrorNoToken

	// Multiple users found with the given token
	LoginErrorMultipleToken
)

// Errors happening on login (not from database/crypting errors)
type LoginError struct {
	err string

	// Error code. See constants.
	ErrorCode uint32
}

// Error generate a readable error string for a LoginError
func (lerr LoginError) Error() string {
	if lerr.err != "" {
		return lerr.err
	}
	return "Login failed"
}

type noUserFoundError struct {
	username string
	token    string
}

type multipleUserFound struct {
	username string
	token    string
}

type alreadyCreatedUserError struct {
	username string
}

func (err multipleUserFound) Error() string {
	if err.username != "" {
		return "Multiple users found with that username: " + err.username
	}
	if err.token != "" {
		return "Multiple users found with that token: " + err.token
	}
	return "Multiple users found"
}

func (err noUserFoundError) Error() string {
	if err.username != "" {
		return "No user found with that username: " + err.username
	}
	if err.token != "" {
		return "No user found with that token: " + err.token
	}
	return "No user found"
}

func (err alreadyCreatedUserError) Error() string {
	return "An user with that username already exists: " + err.username
}
