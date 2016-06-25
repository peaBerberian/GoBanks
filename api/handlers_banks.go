package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/peaberberian/GoBanks/auth"
	"github.com/peaberberian/GoBanks/database"
)

// DBBank properties gettable through this handler
var gettable_bank_fields = []string{
	"Id",
	"UserId",
	"Name",
	"Description",
}

// handleBanks is the main handler for call on the /bank api. It dispatches
// to other function based on the HTTP method used the typical REST CRUD
// naming scheme.
func handleBanks(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	switch r.Method {
	case "GET":
		handleBankRead(w, r, t)
	case "POST":
		handleBankCreate(w, r, t)
	case "PUT":
		handleBankUpdate(w, r, t)
	case "DELETE":
		handleBankDelete(w, r, t)
	default:
		handleNotSupportedMethod(w, r.Method)
	}
}

// handleBankRead handle GET requests on the /bank API
func handleBankRead(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	var id, hasIdInUrl = getApiId(r.URL.Path)
	var queryString = r.URL.Query()
	var f database.DBBankFilters

	// always filter on the current user
	addIntFilter(t.UserId, &f.UserId)

	if hasIdInUrl {
		addIntArrayFilter([]int{id}, &f.Ids)
	}

	// if only some banks are wanted, construct filter
	queryStringToStringArrayFilter(queryString, "name", &f.Names)
	limit := getQueryStringLimit(queryString)

	vals, err := database.GoDB.GetBanks(f, gettable_bank_fields, limit)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	// if an id was given, we're awaiting a single element.
	if hasIdInUrl {
		if len(vals) == 0 {
			fmt.Fprintf(w, "{}")
		} else {
			fmt.Fprintf(w, generateBankResponse(vals[0]))
		}
		return
	}

	if len(vals) == 0 {
		fmt.Fprintf(w, "[]")
	} else {
		fmt.Fprintf(w, generateBanksResponse(vals))
	}
}

// handleBankCreate handle POST requests on the /bank API
func handleBankCreate(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	// you cannot post on an id
	if _, hasIdInUrl := getApiId(r.URL.Path); hasIdInUrl {
		handleNotSupportedMethod(w, r.Method)
		return
	}

	if err := checkMandatoryFields(r, []string{"name"}); err != nil {
		handleError(w, err)
		return
	}

	// translate data to a DBBankParams
	bankElem, err := httpRequestToBank(r)
	if err != nil {
		handleError(w, genericOperationError{})
		return
	}

	// attach elem to current user
	bankElem.UserId = t.UserId
	bank, err := database.GoDB.AddBank(bankElem)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	fmt.Fprintf(w, generateBankResponse(bank))
}

// handleBankUpdate handle PUT requests on the /bank API
func handleBankUpdate(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	var id, hasIdInUrl = getApiId(r.URL.Path)

	if !hasIdInUrl {
		handleBankReplace(w, r, t)
		return
	}

	if err := checkPermissionForBank(t, id); err != nil {
		handleError(w, err)
		return
	}

	jsonElem, err := httpRequestToMap(r)
	if err != nil {
		handleError(w, genericOperationError{})
		return
	}

	var fields []string
	var bankElem database.DBBankParams

	// TODO reflection
	if val, ok := jsonElem["name"]; ok {
		if str, ok := val.(string); !ok {
			handleError(w, bodyParsingError{})
			return
		} else {
			bankElem.Name = str
			fields = append(fields, "Name")
		}
	}

	if val, ok := jsonElem["description"]; ok {
		if str, ok := val.(string); !ok {
			handleError(w, bodyParsingError{})
			return
		} else {
			bankElem.Description = str
			fields = append(fields, "Description")
		}
	}

	var f database.DBBankFilters
	f.Ids.SetFilter([]int{id})

	if err = database.GoDB.UpdateBanks(f, fields, bankElem); err != nil {
		handleError(w, queryOperationError{})
		return
	}

	handleSuccess(w, r)
}

// handleBankDelete handle DELETE requests on the /bank API
func handleBankDelete(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	var id, hasIdInUrl = getApiId(r.URL.Path)
	var f database.DBBankFilters

	if hasIdInUrl {
		if err := checkPermissionForBank(t, id); err != nil {
			handleError(w, err)
			return
		}

		f.Ids.SetFilter([]int{id})
	}

	f.UserId.SetFilter(t.UserId)

	if err := database.GoDB.RemoveBanks(f); err != nil {
		handleError(w, queryOperationError{})
		return
	}
	handleSuccess(w, r)
}

// handleBankReplace handle specifically PUT requests on the main /bank API
// (not restricted to a certain id).
func handleBankReplace(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	var bnks, err = httpRequestToBanks(r)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	var f database.DBBankFilters
	f.UserId.SetFilter(t.UserId)

	if err := database.GoDB.RemoveBanks(f); err != nil {
		handleError(w, queryOperationError{})
		return
	}

	for _, bnk := range bnks {
		bnk.UserId = t.UserId
		if _, err := database.GoDB.AddBank(bnk); err != nil {
			handleError(w, queryOperationError{})
			return
		}
	}
	handleSuccess(w, r)
}

// checkPermissionForBank checks if an user related to the given token has
// the given bankId. It returns an error if the user doesn't have this bank
// or if the database query failed.
func checkPermissionForBank(t *auth.UserToken, bankId int) OperationError {
	if val, err := userHasBank(t.UserId, bankId); !val {
		return notPermittedOperationError{}
	} else if err != nil {
		return queryOperationError{}
	}
	return nil
}

// userHasBank checks if the userId given possess the bankId also given in
// argument. It can return an error if the database query failed.
func userHasBank(userId int, bankId int) (bool, error) {
	var f database.DBBankFilters

	f.UserId.SetFilter(userId)

	f.Ids.SetFilter([]int{bankId})

	var fields = gettable_bank_fields
	val, err := database.GoDB.GetBanks(f, fields, 0)
	if len(val) == 0 || err != nil {
		return false, err
	}

	return true, nil
}

// generateBankResponse generates a JSON string representing the DBBank
// struct provided for the API user. If the marshalling fails or if the
// result is nil, an empty JSON object is returned ('{}')
func generateBankResponse(bnk database.DBBank) string {
	var resJson = dbBankToBankJSON(bnk)

	resBytes, err := json.Marshal(resJson)
	if err != nil || resBytes == nil {
		return "{}"
	}
	return string(resBytes)
}

// generateBankResponse generates a JSON string representing a collection
// of DBBank structs provided for the API user. If the marshalling fails or
// if the result is nil, an empty JSON array is returned ('[]')
func generateBanksResponse(bnk []database.DBBank) string {
	var resJson []BankJSON
	for _, t := range bnk {
		resJson = append(resJson, dbBankToBankJSON(t))
	}
	resBytes, err := json.Marshal(resJson)
	if err != nil || resBytes == nil {
		return "[]"
	}
	return string(resBytes)
}

// httpRequestToBank generates a DBBankParams struct from the body passed
// with the given http request.
func httpRequestToBank(r *http.Request) (database.DBBankParams, error) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return database.DBBankParams{}, err
	}
	return unmarshallToDBBankParams(bodyBytes)
}

// httpRequestsToBank generates an array of DBBankParams structs from the
// body passed with the given http request.
func httpRequestToBanks(r *http.Request) ([]database.DBBankParams, error) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []database.DBBankParams{}, err
	}
	return unmarshallToDBBankParamsArray(bodyBytes)
}

// dbBankToBankJSON takes a DBBank and convert it to its corresponding
// BankJSON response.
func dbBankToBankJSON(bnk database.DBBank) BankJSON {
	return BankJSON{
		Id:          bnk.Id,
		Name:        bnk.Name,
		Description: bnk.Description,
	}
}

// bankJSONToDBBankParams takes a BankJSON and convert it to its
// corresponding DBBankParams struct (used for the database).
func bankJSONToDBBankParams(bnkj BankJSON) database.DBBankParams {
	return database.DBBankParams{
		Name:        bnkj.Name,
		Description: bnkj.Description,
	}
}

// unmarshallToDBBankParams unmarshall an array of bytes into a
// DBBankParams struct
func unmarshallToDBBankParams(input []byte) (database.DBBankParams, error) {
	var inputJson BankJSON
	if err := json.Unmarshal([]byte(input), &inputJson); err != nil {
		return database.DBBankParams{}, err
	}
	return bankJSONToDBBankParams(inputJson), nil
}

// unmarshallToDBBankParamsArray unmarshall an array of bytes into an array
// of DBBankParams structs.
func unmarshallToDBBankParamsArray(input []byte) ([]database.DBBankParams,
	error) {
	var res []database.DBBankParams
	var inputJson []BankJSON
	if err := json.Unmarshal([]byte(input), &inputJson); err != nil {
		return []database.DBBankParams{}, err
	}
	for _, val := range inputJson {
		res = append(res, bankJSONToDBBankParams(val))
	}
	return res, nil
}
