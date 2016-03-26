package login

import "crypto/rand"
import "io"
import "github.com/peaberberian/GoBanks/database/types"
import "golang.org/x/crypto/bcrypt"

type NoUserFoundError struct {
	username string
	token    string
}

type MultipleUserFound struct {
	username string
	token    string
}

type AlreadyCreatedUserError struct {
	username string
}

func (err MultipleUserFound) Error() string {
	if err.username != "" {
		return "Multiple users found with that username: " + err.username
	}
	if err.token != "" {
		return "Multiple users found with that token: " + err.token
	}
	return "Multiple users found"
}

func (err NoUserFoundError) Error() string {
	if err.username != "" {
		return "No user found with that username: " + err.username
	}
	if err.token != "" {
		return "No user found with that token: " + err.token
	}
	return "No user found"
}

func (err AlreadyCreatedUserError) Error() string {
	return "An user with that username already exists: " + err.username
}

func NewUser(username string, password string) (user types.User,
	err error) {
	byteSalt := make([]byte, 32)
	_, err = io.ReadFull(rand.Reader, byteSalt)
	var salt = string(byteSalt)
	if err != nil {
		return user, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(salt+password), 4)
	if err != nil {
		return user, err
	}
	user = types.User{
		Name:         username,
		PasswordHash: string(hash),
		Salt:         salt,
		Permanent:    true,
	}
	return user, nil
}

func AuthenticateUser(user types.User, password string) (err error) {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash),
		[]byte(user.Salt+password))
}

func RegisterUser(db types.GoBanksDataBase, username string,
	password string) (user types.User, err error) {

	usernameTaken, err := isUsernameTaken(db, username)
	if err != nil {
		return
	}
	if usernameTaken {
		return types.User{}, AlreadyCreatedUserError{username}
	}

	user, err = NewUser(username, password)
	if err != nil {
		return
	}

	_, err = db.AddUser(user)
	if err != nil {
		return
	}

	return
}

// TODO look at JWT
func generateToken() string {
	return "aaaa"
}

func LoginUser(db types.GoBanksDataBase, username string,
	password string) (string, error) {

	user, err := GetUserFromUsername(db, username)
	if err != nil {
		return "", err
	}
	err = AuthenticateUser(user, password)
	if err != nil {
		return "", err
	}

	// Add token
	user.Token = generateToken()
	err = db.UpdateUser(user)
	if err != nil {
		return "", err
	}

	return user.Token, nil
}

func GetUserFromUsername(db types.GoBanksDataBase, username string) (types.User,
	error) {

	// setting filters
	var f types.UserFilters
	f.Filters.Names = true
	f.Values.Names = []string{username}

	users, err := db.GetUsers(f)
	if len(users) < 1 {
		return types.User{}, NoUserFoundError{username: username}
	}
	if len(users) >= 2 {
		return users[0], MultipleUserFound{username: username}
	}

	return users[0], err
}

func GetUserFromToken(db types.GoBanksDataBase, token string) (types.User,
	error) {
	// setting filters
	var f types.UserFilters
	f.Filters.Tokens = true
	f.Values.Tokens = []string{token}

	users, err := db.GetUsers(f)
	if len(users) < 1 {
		return users[0], NoUserFoundError{token: token}
	}
	if len(users) >= 2 {
		return users[0], MultipleUserFound{token: token}
	}
	return users[0], err
}

func isUsernameTaken(db types.GoBanksDataBase, username string) (bool, error) {
	// setting filters
	var f types.UserFilters
	f.Filters.Names = true
	f.Values.Names = []string{username}
	return checkExists(db, f)
}

func isTokenTaken(db types.GoBanksDataBase, token string,
) (bool, error) {
	// setting filters
	var f types.UserFilters
	f.Filters.Tokens = true
	f.Values.Tokens = []string{token}
	return checkExists(db, f)
}

func checkExists(db types.GoBanksDataBase, f types.UserFilters,
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
