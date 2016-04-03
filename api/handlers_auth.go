package api

import "net/http"
import "log"
import "fmt"
import "time"
import "encoding/json"

import "github.com/peaberberian/GoBanks/auth"

// handleAuthenticateAPI is called each time an user called the authenticate
// route.
func handleAuthentication(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	if r.Method != "POST" {
		handleNotSupportedMethod(w, r.Method)
		return
	}

	// 1 - parse request
	decoder := json.NewDecoder(r.Body)
	var authJson AuthenticationJSON
	err := decoder.Decode(&authJson)
	if err != nil {
		handleError(w, bodyParsingError{})
		return
	}

	var user string = authJson.User
	var password string = authJson.Password

	// 2 - login user
	token, err := auth.LoginUser(user, password)
	if err != nil {
		handleError(w, err)
		return
	}

	// 3 - send back token
	var expiration int = auth.GetTokenExpiration() * int(time.Hour)
	fmt.Fprintf(w, generateTokenResponse(token, expiration))

	log.Println(user, "just logged in")
}

// getTokenFromRequest recuperates the token string from an http request.
func getTokenFromRequest(r *http.Request) string {
	// Token should be in the Authorization header
	return r.Header.Get("Authorization")
}

// generateTokenResponse generates a JSON ready to be sent describing
// the jwt token given in argument. See TokenJSON for more informations.
func generateTokenResponse(token string, expires int) string {
	var resJSON = TokenJSON{
		Token:     token,
		TokenType: "bearer",
		Expires:   expires / 1e6,
	}
	resBytes, err := json.Marshal(resJSON)
	if err != nil {
		return ""
	} else {
		return string(resBytes)
	}
}
