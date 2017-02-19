package api

import "log"
import "net/http"

import "github.com/peaberberian/GoBanks/auth"

var apiCalls = map[string]string{
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
	if route != apiCalls["authentication"] {
		var tokenString = getTokenFromRequest(r)
		var err error
		token, err = auth.ParseToken(tokenString)
		if err != nil {
			handleError(w, err)
			return
		}
	}

	switch route {
	case apiCalls["authentication"]:
		handleAuthentication(w, r, &token)
	case apiCalls["transactions"]:
		handleTransactions(w, r, &token)
	case apiCalls["banks"]:
		handleBanks(w, r, &token)
	case apiCalls["accounts"]:
		handleAccounts(w, r, &token)
	case apiCalls["categories"]:
		handleCategories(w, r, &token)
	default:
		http.NotFound(w, r)
	}
}

// routeIsInApi simply checks if the given route is in the
// apiCalls map values
func routeIsInApi(route string) bool {
	for _, val := range apiCalls {
		if val == route {
			return true
		}
	}
	return false
}
