package auth

import "fmt"
import "time"

import jwt "github.com/dgrijalva/jwt-go"

// Duration of a token lifetime.
var jwtExpiration int = 1

// Signing key used for JSON Web Tokens
var tokenSigningKey string

// Representation of the principal attributes of our json web token
type UserToken struct {
	ExpirationDate  time.Time
	UserId          int
	IsAdministrator bool
}

// SetTokenExpiration modifies the duration of a token's lifetime.
// Applicable as soon as new user is logged in (LoginUser function)
func SetTokenExpiration(exp int) {
	jwtExpiration = exp
}

// GetTokenExpiration returns the current duration for a token lifetime.
func GetTokenExpiration() int {
	return jwtExpiration
}

// ParseToken takes in argument the token string and returns an easily
// readable UserToken struct which repeats most of the token properties.
// Returns an error if the token is invalid
func ParseToken(tokenString string) (UserToken, AuthenticationError) {
	if tokenString == "" {
		return UserToken{}, noTokenError{}
	}

	jwToken, err := jwt.Parse(tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSigningKey), nil
		})

	if err != nil || !jwToken.Valid {
		return UserToken{}, invalidTokenError{}
	}

	if _, ok := jwToken.Method.(*jwt.SigningMethodHMAC); !ok {
		var alg = fmt.Sprintf("%s", jwToken.Header["alg"])
		return UserToken{}, invalidTokenSigningMethodError{alg}
	}

	var expirationDate time.Time
	var userId int
	var isAdministrator bool
	var ok bool

	var createUnreadableError = func(field string) (UserToken,
		AuthenticationError) {

		return UserToken{}, unreadableTokenError{field}
	}

	if userId64, ok := jwToken.Claims["uid"].(float64); !ok {
		return createUnreadableError("uid")
	} else {
		userId = int(userId64)
	}

	if isAdministrator, ok = jwToken.Claims["adm"].(bool); !ok {
		return createUnreadableError("adm")
	}

	var exp64 float64
	if exp64, ok = jwToken.Claims["exp"].(float64); !ok {
		return createUnreadableError("exp")
	}
	expirationDate = time.Unix(int64(exp64), 0)

	var token = UserToken{
		UserId:          userId,
		IsAdministrator: isAdministrator,
		ExpirationDate:  expirationDate,
	}

	if token.ExpirationDate.Unix() <= time.Now().Unix() {
		return UserToken{}, expiredTokenError{}
	}

	return token, nil
}

// createToken creates a new token string for a specific user.
// Returns an error if a problem with the databases was encountered or if
// the signing key is not secure enough.
func createToken(username string) (string, AuthenticationError) {
	if tokenSigningKey == "" {
		return "", invalidSigningKeyError{}
	}

	user, err := getUserFromUsername(username)
	if err != nil {
		return "", err
	}

	var dur = time.Hour * time.Duration(GetTokenExpiration())
	var expirationDate = time.Now().Add(dur).UnixNano() / 1e6
	var userId = user.Id
	var isAdmin = user.Administrator

	jwToken := jwt.New(jwt.SigningMethodHS256)
	jwToken.Claims["exp"] = expirationDate
	jwToken.Claims["uid"] = userId
	jwToken.Claims["adm"] = isAdmin
	tokenString, serr := jwToken.SignedString([]byte(tokenSigningKey))
	if serr != nil {
		return "", tokenSigningError{}
	}

	return tokenString, nil
}

// TODO put back like it was
// This is done right now for faster tests
func generateSigningKey() error {
	tokenSigningKey = "banana"
	return nil
	// var err error
	// tokenSigningKey, err = generateRandomKey(32)
	// return err
}
