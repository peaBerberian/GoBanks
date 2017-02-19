package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/peaberberian/GoBanks/auth"
	"github.com/peaberberian/GoBanks/database"
)

// DBAccount properties gettable through this handler
var gettable_account_fields = []string{
	"Id",
	"BankId",
	"Name",
	"Description",
}

// handleAccounts is the main handler for call on the /accounts api. It
// dispatches to other function based on the HTTP method used the typical
// REST CRUD naming scheme.
func handleAccounts(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	switch r.Method {
	case "GET":
		handleAccountRead(w, r, t)
	case "POST":
		handleAccountCreate(w, r, t)
	case "PUT":
		handleAccountUpdate(w, r, t)
	case "DELETE":
		handleAccountDelete(w, r, t)
	default:
		handleNotSupportedMethod(w, r.Method)
	}
}

// handleAccountRead handle GET requests on the /accounts API
func handleAccountRead(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	// look if we have an id (GET /accounts/35 => id == 35)
	var id, hasIDinURL = getApiId(r.URL.Path)

	var queryString = r.URL.Query()
	var f database.DBAccountFilters
	var limit int

	// recuperate every bank attached to this user.
	// (blocking database request here :(, TODO see what I can do, cache?)
	bankIds, err := getBankIdsForUserId(t.UserId)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}
	f.BankIds.SetFilter(bankIds)

	// if an id was set in the url, filter to the record corresponding to it
	if hasIDinURL {
		f.Ids.SetFilter([]int{id})
	} else {
		// if only some bank account names are wanted, filter
		wantedAccountNames, _ := queryStringPropertyToStringArray(queryString, "name")
		if len(wantedAccountNames) > 0 {
			f.Names.SetFilter(wantedAccountNames)
		}

		// if only some bank account ids are wanted, filter
		wantedAccountIds, _ := queryStringPropertyToIntArray(queryString, "id")
		if len(wantedAccountIds) > 0 {
			f.Ids.SetFilter(wantedAccountIds)
		}

		// if only some bank ids are wanted, filter
		wantedBankIds, _ := queryStringPropertyToIntArray(queryString, "bankId")
		if len(wantedBankIds) > 0 {
			f.BankIds.SetFilter(wantedBankIds)
		}

		// obtain limit of wanted records, if set
		limit, _ = queryStringPropertyToInt(queryString, "limit")
	}

	// perform the database request
	vals, err := database.GoDB.GetAccounts(f, gettable_account_fields, uint(limit))
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	// if an id was given, we're awaiting an object, not an array.
	if hasIDinURL {
		if len(vals) == 0 {
			fmt.Fprintf(w, "{}")
		} else {
			fmt.Fprintf(w, generateAccountResponse(vals[0]))
		}
		return
	}

	// else respond directly with the result
	if len(vals) == 0 {
		fmt.Fprintf(w, "[]")
	} else {
		fmt.Fprintf(w, generateAccountsResponse(vals))
	}
}

// handleAccountCreate handle POST requests on the /accounts API
func handleAccountCreate(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	// you cannot post on a specific id, reject if you want to do that
	if _, hasID := getApiId(r.URL.Path); hasID {
		handleNotSupportedMethod(w, r.Method)
		return
	}

	// recuperate every bank attached to this user.
	// (blocking database request here :(, TODO see what I can do, cache?)
	bankIds, err := getBankIdsForUserId(t.UserId)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	// convert body to map[string]interface{}
	bodyMap, err := readBodyAsStringMap(r.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	// translate data into a DBAccountParams element
	// (also check mandatory fields)
	accountElem, err := inputToAccountParams(bodyMap)
	if err != nil {
		handleError(w, err)
		return
	}

	// check if the user has the bank he tries to add to
	if !intInArray(accountElem.BankId, bankIds) {
		handleError(w, notPermittedOperationError{})
		return
	}

	// perform database add request
	account, err := database.GoDB.AddAccount(accountElem)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	fmt.Fprintf(w, generateAccountResponse(account))
}

// handleAccountUpdate handle PUT requests on the /accounts API
func handleAccountUpdate(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	// look if we have an id (GET /accounts/35 => id == 35)
	var id, hasID = getApiId(r.URL.Path)

	// if an id was found, it means that we want to replace an element
	// redirect to the right function
	if !hasID {
		handleAccountReplace(w, r, t)
		return
	}

	// recuperate every bank ids associated to this user
	// (blocking database request here :(, TODO see what I can do, cache?)
	bankIds, err := getBankIdsForUserId(t.UserId)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	// recuperate every account ids associated to this user
	// (blocking database request here :(, TODO see what I can do, cache?)
	accountIds, err := getAccountIdsForBankIds(bankIds)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	// if the wanted account does not belong to the user, reject
	if !intInArray(id, accountIds) {
		handleError(w, notPermittedOperationError{})
		return
	}

	// convert body to map[string]interface{}
	bodyMap, err := readBodyAsStringMap(r.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	// -- check fields and update only the ones there --

	var fields []string
	var accountElem database.DBAccountParams

	if val, ok := bodyMap["name"]; ok {
		if str, ok := val.(string); !ok {
			handleError(w, bodyParsingError{})
			return
		} else {
			accountElem.Name = str
			fields = append(fields, "Name")
		}
	}
	if val, ok := bodyMap["description"]; ok {
		if str, ok := val.(string); !ok {
			handleError(w, bodyParsingError{})
			return
		} else {
			accountElem.Description = str
			fields = append(fields, "Description")
		}
	}

	// Filter the account id
	var f database.DBAccountFilters
	f.Ids.SetFilter([]int{id})

	// perform the database request
	if err = database.GoDB.UpdateAccounts(f, fields,
		accountElem); err != nil {
		handleError(w, queryOperationError{})
		return
	}

	handleSuccess(w, r)
}

// handleAccountDelete handle DELETE requests on the /accounts API
func handleAccountDelete(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	// look if we have an id (GET /banks/35 => id == 35)
	var id, hasID = getApiId(r.URL.Path)

	var f database.DBAccountFilters

	// recuperate every bank ids associated to this user
	// (blocking database request here :(, TODO see what I can do, cache?)
	bankIds, err := getBankIdsForUserId(t.UserId)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	// if we have an id, check permission and set filter
	if hasID {
		// recuperate every account ids associated to this user
		// (blocking database request here :(, TODO see what I can do, cache?)
		accountIds, err := getAccountIdsForBankIds(bankIds)
		if err != nil {
			handleError(w, queryOperationError{})
			return
		}

		// (blocking database request here :(, TODO see what I can do, cache?)
		if !intInArray(id, accountIds) {
			handleError(w, notPermittedOperationError{})
			return
		}
		f.Ids.SetFilter([]int{id})
	} else {
		// filter by bankId
		f.BankIds.SetFilter(bankIds)
	}

	// perform the database request
	if err := database.GoDB.RemoveAccounts(f); err != nil {
		handleError(w, queryOperationError{})
		return
	}
	handleSuccess(w, r)
}

// handleAccountReplace handle specifically PUT requests on the main /accounts
// API.
// (not restricted to a certain id).
func handleAccountReplace(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	var bodyMaps, err = readBodyAsArrayOfStringMap(r.Body)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	// recuperate every account ids associated to this user
	// (blocking database request here :(, TODO see what I can do, cache?)
	bankIds, err := getBankIdsForUserId(t.UserId)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	var accs []database.DBAccountParams

	// translate data into DBBankParams elements
	// (also check mandatory fields)
	for _, bodyMap := range bodyMaps {
		accElem, err := inputToAccountParams(bodyMap)
		if err != nil {
			handleError(w, err)
			return
		}

		// check that the bankId indicated is attached to this user
		if !intInArray(accElem.BankId, bankIds) {
			handleError(w, notPermittedOperationError{})
			return
		}
		accs = append(accs, accElem)
	}

	// Remove old accounts linked to this user
	var f database.DBAccountFilters
	f.BankIds.SetFilter(bankIds)

	if err := database.GoDB.RemoveAccounts(f); err != nil {
		handleError(w, queryOperationError{})
		return
	}

	// add each account indicated to the database
	for _, acc := range accs {
		if _, err := database.GoDB.AddAccount(acc); err != nil {
			handleError(w, queryOperationError{})
			return
		}
	}
	handleSuccess(w, r)
}

// generateAccountResponse generates a JSON string representing the DBAccount
// struct provided for the API user. If the marshalling fails or if the
// result is nil, an empty JSON object is returned ('{}')
func generateAccountResponse(acc database.DBAccount) string {
	var resJSON = dbAccountToAccountJSON(acc)

	resBytes, err := json.Marshal(resJSON)
	if err != nil || resBytes == nil {
		return "{}"
	}
	return string(resBytes)
}

// generateAccountResponse generates a JSON string representing a collection
// of DBAccount structs provided for the API user. If the marshalling fails or
// if the result is nil, an empty JSON array is returned ('[]')
func generateAccountsResponse(acc []database.DBAccount) string {
	var resJSON []AccountJSON
	for _, t := range acc {
		resJSON = append(resJSON, dbAccountToAccountJSON(t))
	}
	resBytes, err := json.Marshal(resJSON)
	if err != nil || resBytes == nil {
		return "[]"
	}
	return string(resBytes)
}

// dbAccountToAccountJSON takes a DBAccount and convert it to its corresponding
// AccountJSON response.
func dbAccountToAccountJSON(acc database.DBAccount) AccountJSON {
	return AccountJSON{
		Id:          acc.Id,
		Name:        acc.Name,
		Description: acc.Description,
		BankId:      acc.BankId,
	}
}

// process map[string]interface{} input to create a DBAccountParams object.
// if mandatory fields are not found, this function returns an error.
func inputToAccountParams(input map[string]interface{}) (database.DBAccountParams, error) {
	var res database.DBAccountParams
	var valid bool

	// The "name" field is mandatory
	res.Name, valid = input["name"].(string)
	if !valid {
		return res, missingParameterError{"name"}
	}

	// The "name" field is mandatory
	// TODO check why can't cast directly to int
	// (don't ask me why. It just werks...)
	bankIDStr, valid := input["bankId"].(float64)
	if !valid {
		return res, missingParameterError{"bankId"}
	}

	res.BankId = int(bankIDStr)
	res.Description, _ = input["description"].(string)

	return res, nil
}

// // readBodyAsAccountParams generates a DBBankParams structs from the
// // body passed with the given http request.
// func readBodyAsAccountParamsArray(r io.reader) (database.DBAccountParams,
// 	error) {
// 	bodyBytes, err := ioutil.ReadAll(r)
// 	if err != nil {
// 		return database.DBAccountParams{}, err
// 	}
// 	return parseAccountJson(string(bodyBytes))
// }

// // readBodyAsAccountParamsArray generates an array of DBBankParams structs from the
// // body passed with the given http request.
// func readBodyAsAccountParamsArray(r io.reader) ([]database.DBAccountParams,
// 	error) {
// 	bodyBytes, err := ioutil.ReadAll(r)
// 	if err != nil {
// 		return []database.DBAccountParams{}, err
// 	}
// 	return parseAccountsJson(string(bodyBytes))
// }

// func accountJSONToDBAccountParams(accj AccountJSON) database.DBAccountParams {
// 	return database.DBAccountParams{
// 		Name:        accj.Name,
// 		Description: accj.Description,
// 		BankId:      accj.BankId,
// 	}
// }

// func parseAccountJson(input string) (database.DBAccountParams, error) {
// 	var inputJson AccountJSON
// 	if err := json.Unmarshal([]byte(input), &inputJson); err != nil {
// 		return database.DBAccountParams{}, err
// 	}
// 	return accountJSONToDBAccountParams(inputJson), nil
// }

// func parseAccountsJson(input string) ([]database.DBAccountParams, error) {
// 	var res []database.DBAccountParams
// 	var inputJson []AccountJSON
// 	if err := json.Unmarshal([]byte(input), &inputJson); err != nil {
// 		return []database.DBAccountParams{}, err
// 	}
// 	for _, val := range inputJson {
// 		res = append(res, accountJSONToDBAccountParams(val))
// 	}
// 	return res, nil
// }
