package login

import "crypto/rand"
import "io"
import "golang.org/x/crypto/bcrypt"

import def "github.com/peaberberian/GoBanks/database/definitions"

func LoginUser(db def.GoBanksDataBase, username string,
	password string) (string, error) {

	user, err := GetUserFromUsername(db, username)
	if err != nil {
		return "", err
	}
	err = authenticateUser(user, password)
	if err != nil {
		return "", err
	}

	// Add token TODO look at JWT
	user.Token = generateToken()
	err = db.UpdateUser(user)
	if err != nil {
		return "", err
	}

	return user.Token, nil
}

func RegisterUser(db def.GoBanksDataBase, username string,
	password string) (user def.User, err error) {

	usernameTaken, err := isUsernameTaken(db, username)
	if err != nil {
		return
	}
	if usernameTaken {
		err = alreadyCreatedUserError{username: username}
		return def.User{}, LoginError{err: err.Error(),
			ErrorCode: LoginErrorAlreadyTakenUsername}
	}

	user, err = newUser(username, password)
	if err != nil {
		return
	}

	_, err = db.AddUser(user)
	if err != nil {
		return
	}

	return
}

func GetUserFromUsername(db def.GoBanksDataBase, username string,
) (def.User,
	error) {

	// setting filters
	var f def.UserFilters
	f.Filters.Names = true
	f.Values.Names = []string{username}

	users, err := db.GetUsers(f)
	if len(users) < 1 {
		err = noUserFoundError{username: username}
		return def.User{}, LoginError{err: err.Error(),
			ErrorCode: LoginErrorNoUsername,
		}
	}
	if len(users) >= 2 {
		err = multipleUserFound{username: username}
		return def.User{}, LoginError{err: err.Error(),
			ErrorCode: LoginErrorMultipleUsername,
		}
	}

	return users[0], err
}

func GetUserFromToken(db def.GoBanksDataBase, token string) (def.User,
	error) {
	// setting filters
	var f def.UserFilters
	f.Filters.Tokens = true
	f.Values.Tokens = []string{token}

	users, err := db.GetUsers(f)
	if len(users) < 1 {
		err = noUserFoundError{token: token}
		return def.User{}, LoginError{err: err.Error(),
			ErrorCode: LoginErrorNoToken,
		}
	}
	if len(users) >= 2 {
		err = multipleUserFound{token: token}
		return def.User{}, LoginError{err: err.Error(),
			ErrorCode: LoginErrorMultipleToken,
		}
	}
	return users[0], err
}

func newUser(username string, password string) (user def.User,
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
	user = def.User{
		Name:         username,
		PasswordHash: string(hash),
		Salt:         salt,
		Permanent:    true,
	}
	return user, nil
}

func authenticateUser(user def.User, password string) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash),
		[]byte(user.Salt+password))
	if err != nil {
		return LoginError{err: err.Error(),
			ErrorCode: LoginErrorWrongPassword}
	}
	return nil
}

// TODO look at JWT
func generateToken() string {
	return "aaaa"
}

func isUsernameTaken(db def.GoBanksDataBase, username string) (bool,
	error) {
	// setting filters
	var f def.UserFilters
	f.Filters.Names = true
	f.Values.Names = []string{username}
	return checkExists(db, f)
}

func isTokenTaken(db def.GoBanksDataBase, token string,
) (bool, error) {
	// setting filters
	var f def.UserFilters
	f.Filters.Tokens = true
	f.Values.Tokens = []string{token}
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
