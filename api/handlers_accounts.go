package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/peaberberian/GoBanks/auth"
	"github.com/peaberberian/GoBanks/database"
)

var gettable_account_infos = []string{"Id", "BankId", "Name", "Description"}

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

func handleAccountRead(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	var id, hasIdInUrl = getApiId(r.URL.Path)
	var queryString = r.URL.Query()
	var f database.DBAccountFilters

	// always filter on the current user
	f.UserId.Activated = true
	f.UserId.Value = t.UserId

	if hasIdInUrl {
		f.Ids.Activated = true
		f.Ids.Value = []int{id}
	}

	// if only some accounts are wanted, construct filter
	queryStringToStringArrayFilter(queryString, "name", &f.Names)
	queryStringToIntArrayFilter(queryString, "id", &f.BankIds)
	queryStringToIntArrayFilter(queryString, "bank", &f.BankIds)
	limit := getQueryStringLimit(queryString)

	vals, err := database.GoDB.GetAccounts(f, gettable_account_infos, limit)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	// if an id was given, we're awaiting a single element.
	if hasIdInUrl {
		if len(vals) == 0 {
			fmt.Fprintf(w, "{}")
		} else {
			fmt.Fprintf(w, generateAccountResponse(vals[0]))
		}
		return
	}

	if len(vals) == 0 {
		fmt.Fprintf(w, "[]")
	} else {
		fmt.Fprintf(w, generateAccountsResponse(vals))
	}
}

func handleAccountCreate(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	if _, hasId := getApiId(r.URL.Path); hasId {
		handleNotSupportedMethod(w, r.Method)
		return
	}

	bankIds, err := getBankIdsForUserId(t.UserId)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	accountElem, err := httpRequestToAccount(r)
	if err != nil {
		handleError(w, genericOperationError{})
		return
	}

	switch {
	case accountElem.Name == "":
		handleError(w, missingParameterError{"name"})
	case accountElem.BankId == 0:
		handleError(w, missingParameterError{"bankId"})
	case !intInArray(accountElem.BankId, bankIds):
		handleError(w, notPermittedOperationError{})
	default:
		acc, err := database.GoDB.AddAccount(accountElem)
		if err != nil {
			handleError(w, queryOperationError{})
			return
		}
		fmt.Fprintf(w, generateAccountResponse(acc))
	}
}

func handleAccountUpdate(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	var id, hasId = getApiId(r.URL.Path)

	if !hasId {
		handleAccountReplace(w, r, t)
		return
	}

	bankIds, err := getBankIdsForUserId(t.UserId)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}
	accountIds, err := getAccountIdsForBankIds(bankIds)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	if !intInArray(id, accountIds) {
		handleError(w, notPermittedOperationError{})
		return
	}

	jsonElem, err := httpRequestToMap(r)
	if err != nil {
		handleError(w, genericOperationError{})
		return
	}

	var fields []string
	var accountElem database.DBAccountParams

	if val, ok := jsonElem["name"]; ok {
		if str, ok := val.(string); !ok {
			handleError(w, bodyParsingError{})
			return
		} else {
			accountElem.Name = str
			fields = append(fields, "Name")
		}
	}

	if val, ok := jsonElem["description"]; ok {
		if str, ok := val.(string); !ok {
			handleError(w, bodyParsingError{})
			return
		} else {
			accountElem.Description = str
			fields = append(fields, "Description")
		}
	}

	var f database.DBAccountFilters
	f.Ids.Activated = true
	f.Ids.Value = []int{id}

	if err = database.GoDB.UpdateAccounts(f, fields,
		accountElem); err != nil {
		handleError(w, queryOperationError{})
		return
	}

	handleSuccess(w, r)
}

func handleAccountDelete(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	var id, hasId = getApiId(r.URL.Path)
	var f database.DBAccountFilters

	bankIds, err := getBankIdsForUserId(t.UserId)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	accountIds, err := getAccountIdsForBankIds(bankIds)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	if hasId {
		if !intInArray(id, accountIds) {
			handleError(w, notPermittedOperationError{})
			return
		}

		f.Ids.Activated = true
		f.Ids.Value = []int{id}
	}

	f.BankIds.Activated = true
	f.BankIds.Value = bankIds

	if err := database.GoDB.RemoveAccounts(f); err != nil {
		handleError(w, queryOperationError{})
		return
	}

	handleSuccess(w, r)
}

func handleAccountReplace(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	var accs, err = httpRequestToAccounts(r)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	bankIds, err := getBankIdsForUserId(t.UserId)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	var f database.DBAccountFilters
	f.BankIds.Activated = true
	f.BankIds.Value = bankIds

	if err := database.GoDB.RemoveAccounts(f); err != nil {
		handleError(w, queryOperationError{})
		return
	}

	for _, acc := range accs {
		if !intInArray(acc.BankId, bankIds) {
			handleError(w, notPermittedOperationError{})
			return
		}

		if _, err := database.GoDB.AddAccount(acc); err != nil {
			handleError(w, queryOperationError{})
			return
		}
	}
	handleSuccess(w, r)
}

func generateAccountResponse(acc database.DBAccount) string {
	var resJson = accountToAccountJson(acc)

	resBytes, err := json.Marshal(resJson)
	if err != nil || resBytes == nil {
		return "{}"
	}
	return string(resBytes)
}

func generateAccountsResponse(acc []database.DBAccount) string {
	var resJson []AccountJSON
	for _, t := range acc {
		resJson = append(resJson, accountToAccountJson(t))
	}
	resBytes, err := json.Marshal(resJson)
	if err != nil || resBytes == nil {
		return "[]"
	}
	return string(resBytes)
}

func httpRequestToAccount(r *http.Request) (database.DBAccountParams,
	error) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return database.DBAccountParams{}, err
	}
	return parseAccountJson(string(bodyBytes))
}

func httpRequestToAccounts(r *http.Request) ([]database.DBAccountParams,
	error) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []database.DBAccountParams{}, err
	}
	return parseAccountsJson(string(bodyBytes))
}

func accountToAccountJson(acc database.DBAccount) AccountJSON {
	return AccountJSON{
		Id:          acc.Id,
		Name:        acc.Name,
		Description: acc.Description,
		BankId:      acc.BankId,
	}
}

func accountJsonToAccount(accj AccountJSON) database.DBAccountParams {
	return database.DBAccountParams{
		Name:        accj.Name,
		Description: accj.Description,
		BankId:      accj.BankId,
	}
}

func parseAccountJson(input string) (database.DBAccountParams, error) {
	var inputJson AccountJSON
	if err := json.Unmarshal([]byte(input), &inputJson); err != nil {
		return database.DBAccountParams{}, err
	}
	return accountJsonToAccount(inputJson), nil
}

func parseAccountsJson(input string) ([]database.DBAccountParams, error) {
	var res []database.DBAccountParams
	var inputJson []AccountJSON
	if err := json.Unmarshal([]byte(input), &inputJson); err != nil {
		return []database.DBAccountParams{}, err
	}
	for _, val := range inputJson {
		res = append(res, accountJsonToAccount(val))
	}
	return res, nil
}
