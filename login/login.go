package login

import "crypto/rand"
import "io"
import "golang.org/x/crypto/bcrypt"

import def "github.com/peaberberian/GoBanks/database/definitions"

func init() {
	_ = generateSigningKey()
}

// LoginUser verifies the password for the given username and returns
// an error if the password is wrong / the user does not exists / other
// database errors
func LoginUser(db def.GoBanksDataBase, username string,
	password string) error {

	user, err := GetUserFromUsername(db, username)
	if err != nil {
		return err
	}
	err = authenticateUser(user, password)
	return err
}

// RegisterUser adds a new user in the database with the given username
// and password. An administrator boolean indicates if the user is an
// administrator (more rights)
func RegisterUser(db def.GoBanksDataBase, username string,
	password string, administrator bool) (user def.User, err error) {

	usernameTaken, err := isUsernameTaken(db, username)
	if err != nil {
		return
	}
	if usernameTaken {
		err = alreadyCreatedUserError{username: username}
		return def.User{}, LoginError{err: err.Error(),
			code: LoginErrorAlreadyTakenUsername}
	}

	user, err = newUser(username, password, administrator)
	if err != nil {
		return
	}

	_, err = db.AddUser(user)
	if err != nil {
		return
	}

	return
}

// GetUserFromUsername returns the corresponding User struct for a given
// username.
// It returns an error if no or multiple users were found with that
// username or for a database error.
func GetUserFromUsername(db def.GoBanksDataBase, username string,
) (def.User, error) {

	// setting filters
	var f def.UserFilters
	f.Filters.Names = true
	f.Values.Names = []string{username}

	users, err := db.GetUsers(f)
	if len(users) < 1 {
		err = noUserFoundError{username: username}
		return def.User{}, LoginError{err: err.Error(),
			code: LoginErrorNoUsername,
		}
	}
	if len(users) >= 2 {
		err = multipleUserFound{username: username}
		return def.User{}, LoginError{err: err.Error(),
			code: LoginErrorMultipleUsername,
		}
	}

	return users[0], err
}

func generateRandomKey(size int) (salt string, err error) {
	byteSalt := make([]byte, size)
	_, err = io.ReadFull(rand.Reader, byteSalt)
	if err != nil {
		return
	}
	salt = string(byteSalt)
	return salt, nil
}

func newUser(username string, password string, administrator bool,
) (user def.User, err error) {
	salt, err := generateRandomKey(32)
	if err != nil {
		return user, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(salt+password), 4)
	if err != nil {
		return user, err
	}
	user = def.User{
		Name:          username,
		PasswordHash:  string(hash),
		Salt:          salt,
		Administrator: administrator,
	}
	return user, nil
}

func authenticateUser(user def.User, password string) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash),
		[]byte(user.Salt+password))
	if err != nil {
		return LoginError{err: err.Error(),
			code: LoginErrorWrongPassword}
	}
	return nil
}

func isUsernameTaken(db def.GoBanksDataBase, username string) (bool,
	error) {
	// setting filters
	var f def.UserFilters
	f.Filters.Names = true
	f.Values.Names = []string{username}
	return checkExists(db, f)
}

func checkExists(db def.GoBanksDataBase, f def.UserFilters,
) (bool, error) {
	users, err := db.GetUsers(f)
	if err != nil {
		return false, err
	}
	if len(users) > 0 {
		return true, nil
	}
	return false, nil
}
