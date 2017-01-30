package auth

import "crypto/rand"
import "io"
import "golang.org/x/crypto/bcrypt"

import db "github.com/peaberberian/GoBanks/database"

func init() {
	_ = generateSigningKey()
}

// LoginUser logins a particular user from its credentials and returns, if
// it succeeded, the json web token for this user.
// If the autentication failed, an AuthenticationError is returned.
func LoginUser(username string, password string) (string,
	AuthenticationError) {

	if err := VerifyUser(username, password); err != nil {
		return "", err
	}

	return createToken(username)
}

// VerifyUser verifies the password for the given username and returns
// an error if the password is wrong / the user does not exists / other
// database errors
func VerifyUser(username string, password string) AuthenticationError {
	user, err := getUserFromUsername(username)
	if err != nil {
		return err
	}
	err = authenticate(user, password)
	return err
}

// RegisterUser adds a new user in the database with the given username
// and password. An administrator boolean indicates if the user is an
// administrator (more rights)
func RegisterUser(username string,
	password string, administrator bool) (db.DBUser, AuthenticationError) {

	usernameTaken, err := isUsernameTaken(username)
	if err != nil {
		return db.DBUser{}, genericAuthenticationError{}
	}
	if usernameTaken {
		return db.DBUser{}, alreadyTakenUsernameError{username}
	}

	userParams, err := newUser(username, password, administrator)
	if err != nil {
		return db.DBUser{}, genericAuthenticationError{}
	}

	user, dbErr := db.GoDB.AddUser(userParams)
	if dbErr != nil {
		return db.DBUser{}, genericAuthenticationError{}
	}

	return user, nil
}

// authenticate tries to authenticates a user based on the db.DBUser object
// and a password. Returns an error if the password is invalid.
func authenticate(user db.DBUser, password string) AuthenticationError {
	var err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash),
		[]byte(user.Salt+password))
	if err != nil {
		return wrongPasswordError{user.Name}
	}
	return nil
}

// getUserFromUsername returns the corresponding User struct for a given
// username.
// It returns an error if no or multiple users were found with that
// username or for a database error.
func getUserFromUsername(username string) (db.DBUser, AuthenticationError) {

	// setting filters
	var f db.DBUserFilters
	f.Name.SetFilter(username)

	var fields = []string{"Id", "Name", "PasswordHash", "Salt"}
	user, err := db.GoDB.GetUser(f, fields)
	if err != nil {
		return db.DBUser{}, genericAuthenticationError{}
	}
	if user.Name == "" {
		return db.DBUser{}, userNotFoundError{username}
	}

	return user, nil
}

func generateRandomKey(size int) (string, error) {
	byteSalt := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, byteSalt)
	if err != nil {
		return "", err
	}
	var salt = string(byteSalt)
	return salt, nil
}

func newUser(username string, password string,
	administrator bool) (db.DBUserParams, error) {
	salt, err := generateRandomKey(32)
	if err != nil {
		return db.DBUserParams{}, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(salt+password), 4)
	if err != nil {
		return db.DBUserParams{}, err
	}
	var user = db.DBUserParams{
		Name:         username,
		PasswordHash: string(hash),
		Salt:         salt,
	}
	return user, nil
}

func isUsernameTaken(username string) (bool, error) {
	// setting filters
	var f db.DBUserFilters
	f.Name.SetFilter(username)
	return checkExists(f)
}

func checkExists(f db.DBUserFilters) (bool, error) {
	var fields = []string{"Id", "Name", "PasswordHash", "Salt"}
	user, err := db.GoDB.GetUser(f, fields)
	if err != nil {
		return false, err
	}
	if user.Name == "" {
		return true, nil
	}
	return false, nil
}
