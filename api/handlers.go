package api

import "log"
import "net/http"

import "github.com/peaberberian/GoBanks/auth"

var api_calls = map[string]string{
	"authentication": "auth",
	"transactions":   "transactions",
	"banks":          "banks",
	"accounts":       "accounts",
	"categories":     "categories",
	"users":          "users",
}

// handlerV1 is the handler for all calls concerning the API version 1
func handlerV1(w http.ResponseWriter, r *http.Request) {
	var route = getApiRoute(r.URL.Path)
	log.Println("Request received for API:", route)

	w.Header().Set("content-type", "application/json")

	var token auth.UserToken

	// only route where the token shouldn't be needed
	if route != api_calls["authentication"] {
		var tokenString = getTokenFromRequest(r)
		var err error
		token, err = auth.ParseToken(tokenString)
		if err != nil {
			handleError(w, err)
			return
		}
	}

	switch route {
	case api_calls["authentication"]:
		handleAuthentication(w, r, &token)
	case api_calls["transactions"]:
		handleTransactions(w, r, &token)
	case api_calls["banks"]:
		handleBanks(w, r, &token)
	case api_calls["accounts"]:
		handleAccounts(w, r, &token)
	default:
		http.NotFound(w, r)
	}
}

func routeIsInApi(route string) bool {
	for _, val := range api_calls {
		if val == route {
			return true
		}
	}
	return false
}
