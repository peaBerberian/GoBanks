package api

import (
	"fmt"
	"net/http"

	"encoding/json"

	"github.com/peaberberian/GoBanks/auth"
	"github.com/peaberberian/GoBanks/database"
)

// DBTransaction properties gettable through this handler
var gettable_transaction_fields = []string{
	"Id",
	"AccountId",
	"Label",
	"CategoryId",
	"Description",
	"TransactionDate",
	"RecordDate",
	"Debit",
	"Credit",
	"Reference",
}

// TODO
// Make relation between DBTransactionFilters names and their property name in
// the query string
var query_string_properties = map[string]string{
	"Ids":                 "id",
	"AccountIds":          "aid",
	"CategoryIds":         "cid",
	"FromTransactionDate": "tfrom",
	"ToTransactionDate":   "tto",
	"FromRecordDate":      "rfrom",
	"ToRecordDate":        "rto",
	"MinDebit":            "mind",
	"MaxDebit":            "maxd",
	"MinCredit":           "minc",
	"MaxCredit":           "maxc",
	"References":          "ref",
}

// madatory_transaction_json_fields list all mandatory parameters when the user
// wants to create a new transaction.
// If one of those field is not present as a request is received, an error
// message with the corresponding field in argument will be returned.
var mandatory_transaction_json_fields = []string{
	"accountId",
}

// handleTransactions is the main handler for call on the /accounts api. It
// dispatches to other function based on the HTTP method used the typical
// REST CRUD naming scheme.
func handleTransactions(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	switch r.Method {
	case "GET":
		handleTransactionRead(w, r, t)
	case "POST":
		handleTransactionCreate(w, r, t)
	case "PUT":
		handleTransactionUpdate(w, r, t)
	case "DELETE":
		handleTransactionDelete(w, r, t)
	default:
		handleNotSupportedMethod(w, r.Method)
	}
}

// handleAccountRead handle GET requests on the /accounts API
func handleTransactionRead(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	// look if we have an id (GET /transactions/35 => id == 35)
	var id, hasIdInUrl = getApiId(r.URL.Path)

	var queryString = r.URL.Query()
	var f database.DBTransactionFilters
	var limit int

	// recuperate every bank attached to this user.
	// (blocking database request here :(, TODO see what I can do, cache?)
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

	f.AccountIds.SetFilter(accountIds)

	// if an id was set in the url, filter to the record corresponding to it
	if hasIdInUrl {
		f.Ids.SetFilter([]int{id})
	} else {
		// if only some transaction ids are wanted, filter
		wantedTransactionIds, _ := queryStringPropertyToIntArray(queryString, "id")
		if len(wantedTransactionIds) > 0 {
			f.Ids.SetFilter(wantedTransactionIds)
		}

		// if only some category ids are wanted, filter
		wantedCategoryIds, _ :=
			queryStringPropertyToIntArray(queryString, "category")
		if len(wantedCategoryIds) > 0 {
			f.CategoryIds.SetFilter(wantedCategoryIds)
		}

		// if a FromTransactionDate timestamp has been provided, filter
		if wantedFromTDate, isDefined :=
			queryStringPropertyToTime(queryString, "tfrom"); isDefined {
			f.FromTransactionDate.SetFilter(wantedFromTDate)
		}

		// if a ToTransactionDate timestamp has been provided, filter
		if wantedToTDate, isDefined :=
			queryStringPropertyToTime(queryString, "tto"); isDefined {
			f.ToTransactionDate.SetFilter(wantedToTDate)
		}

		// if a FromRecordDate timestamp has been provided, filter
		if wantedFromRDate, isDefined :=
			queryStringPropertyToTime(queryString, "rfrom"); isDefined {
			f.FromRecordDate.SetFilter(wantedFromRDate)
		}

		// if a ToRecordDate timestamp has been provided, filter
		if wantedToRDate, isDefined :=
			queryStringPropertyToTime(queryString, "rto"); isDefined {
			f.ToRecordDate.SetFilter(wantedToRDate)
		}

		if wantedMinDebit, isDefined :=
			queryStringPropertyToFloat32(queryString, "min_debit"); isDefined {
			f.MinDebit.SetFilter(wantedMinDebit)
		}

		if wantedMaxDebit, isDefined :=
			queryStringPropertyToFloat32(queryString, "max_debit"); isDefined {
			f.MaxDebit.SetFilter(wantedMaxDebit)
		}

		if wantedMinCredit, isDefined :=
			queryStringPropertyToFloat32(queryString, "min_debit"); isDefined {
			f.MinCredit.SetFilter(wantedMinCredit)
		}

		if wantedMaxCredit, isDefined :=
			queryStringPropertyToFloat32(queryString, "max_debit"); isDefined {
			f.MaxCredit.SetFilter(wantedMaxCredit)
		}

		// if only some transaction names are wanted, filter
		if wantedTransactionReferences, isDefined :=
			queryStringPropertyToStringArray(queryString, "reference"); isDefined {
			f.References.SetFilter(wantedTransactionReferences)
		}

		// obtain limit of wanted records, if set
		limit, _ = queryStringPropertyToInt(queryString, "limit")
	}

	// perform the database request
	vals, err := database.GoDB.GetTransactions(
		f,
		gettable_transaction_fields,
		uint(limit))

	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	// if an id was given, we're awaiting an object, not an array.
	if hasIdInUrl {
		if len(vals) == 0 {
			fmt.Fprintf(w, "{}")
		} else {
			fmt.Fprintf(w, dbTransactionToJSONString(vals[0]))
		}
		return
	}

	// else respond directly with the result
	if len(vals) == 0 {
		fmt.Fprintf(w, "[]")
	} else {
		fmt.Fprintf(w, dbTransactionsToJSONString(vals))
	}
}

// handleTransactionCreate handle POST requests on the /accounts API
func handleTransactionCreate(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	// you cannot post on a specific id, reject if you want to do that
	if _, hasId := getApiId(r.URL.Path); hasId {
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

	accountIds, err := getAccountIdsForBankIds(bankIds)
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

	// translate data into a DBTransactionParams element
	// (also check mandatory fields)
	transactionElem, err := stringMapInputToDBTransactionParams(bodyMap)
	if err != nil {
		handleError(w, err)
		return
	}

	// check if the user has the account he tries to add to
	if !intInArray(transactionElem.AccountId, accountIds) {
		handleError(w, notPermittedOperationError{})
		return
	}

	// perform database add request
	transaction, err := database.GoDB.AddTransaction(transactionElem)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	fmt.Fprintf(w, dbTransactionToJSONString(transaction))
}

// handleTransactionUpdate handle PUT requests on the /transactions API
func handleTransactionUpdate(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	// look if we have an id (GET /transactions/35 => id == 35)
	var id, hasId = getApiId(r.URL.Path)

	// if an id was found, it means that we want to replace an element
	// redirect to the right function
	if !hasId {
		handleTransactionReplace(w, r, t)
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

	// recuperate every transaction ids associated to this user
	// (blocking database request here :(, TODO see what I can do, cache?)
	transactionIds, err := getTransactionIdsForAccountIds(accountIds)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	// if the wanted transaction does not belong to the user, reject
	if !intInArray(id, transactionIds) {
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

	// TODO

	var fields []string
	var transactionElem database.DBTransactionParams

	if val, ok := bodyMap["label"]; ok {
		if str, ok := val.(string); !ok {
			handleError(w, bodyParsingError{})
			return
		} else {
			transactionElem.Label = str
			fields = append(fields, "Label")
		}
	}
	if val, ok := bodyMap["description"]; ok {
		if str, ok := val.(string); !ok {
			handleError(w, bodyParsingError{})
			return
		} else {
			transactionElem.Description = str
			fields = append(fields, "Description")
		}
	}

	// Filter the transaction id
	var f database.DBTransactionFilters
	f.Ids.SetFilter([]int{id})

	// perform the database request
	if err = database.GoDB.UpdateTransactions(f, fields,
		transactionElem); err != nil {
		handleError(w, queryOperationError{})
		return
	}

	handleSuccess(w, r)
}

// handleTransactionDelete handle DELETE requests on the /transactions API
func handleTransactionDelete(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	// look if we have an id (GET /banks/35 => id == 35)
	var id, hasId = getApiId(r.URL.Path)

	var f database.DBTransactionFilters

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

	// if we have an id, check permission and set filter
	if hasId {
		// recuperate every transaction ids associated to this user
		// (blocking database request here :(, TODO see what I can do, cache?)
		transactionIds, err := getTransactionIdsForAccountIds(accountIds)
		if err != nil {
			handleError(w, queryOperationError{})
			return
		}

		if !intInArray(id, transactionIds) {
			handleError(w, notPermittedOperationError{})
			return
		}
		f.Ids.SetFilter([]int{id})
	} else {
		// filter by accountId
		f.AccountIds.SetFilter(accountIds)
	}

	// perform the database request
	if err := database.GoDB.RemoveTransactions(f); err != nil {
		handleError(w, queryOperationError{})
		return
	}
	handleSuccess(w, r)
}

// handleTransactionReplace handle specifically PUT requests on the main /transactions
// API.
// (not restricted to a certain id).
func handleTransactionReplace(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	var bodyMaps, err = readBodyAsArrayOfStringMap(r.Body)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	// recuperate every transaction ids associated to this user
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

	var accs []database.DBTransactionParams

	// translate data into DBBankParams elements
	// (also check mandatory fields)
	for _, bodyMap := range bodyMaps {
		transElem, err := stringMapInputToDBTransactionParams(bodyMap)
		if err != nil {
			handleError(w, err)
			return
		}

		// check that the bankId indicated is attached to this user
		if !intInArray(transElem.AccountId, accountIds) {
			handleError(w, notPermittedOperationError{})
			return
		}
		accs = append(accs, transElem)
	}

	// Remove old transactions linked to this user
	var f database.DBTransactionFilters
	f.AccountIds.SetFilter(accountIds)

	if err := database.GoDB.RemoveTransactions(f); err != nil {
		handleError(w, queryOperationError{})
		return
	}

	// add each transaction indicated to the database
	for _, acc := range accs {
		if _, err := database.GoDB.AddTransaction(acc); err != nil {
			handleError(w, queryOperationError{})
			return
		}
	}
	handleSuccess(w, r)
}

// dbTransactionToJSONString generates a JSON string representing the DBAccount
// struct provided for the API user. If the marshalling fails or if the
// result is nil, an empty JSON object is returned ('{}')
func dbTransactionToJSONString(trn database.DBTransaction) string {
	var resJson = dbTransactionToTransactionJSON(trn)

	resBytes, err := json.Marshal(resJson)
	if err != nil || resBytes == nil {
		return "{}"
	}
	return string(resBytes)
}

// dbTransactionToJSONString generates a JSON string representing a collection
// of DBTransaction structs provided for the API user. If the marshalling fails or
// if the result is nil, an empty JSON array is returned ('[]')
func dbTransactionsToJSONString(trn []database.DBTransaction) string {
	var resJson []TransactionJSON
	for _, t := range trn {
		resJson = append(resJson, dbTransactionToTransactionJSON(t))
	}
	resBytes, err := json.Marshal(resJson)
	if err != nil || resBytes == nil {
		return "[]"
	}
	return string(resBytes)
}

// dbTransactionToTransactionJSON takes a DBTransaction and convert it to its
// corresponding TransactionJSON struct.
func dbTransactionToTransactionJSON(trn database.DBTransaction) TransactionJSON {
	return TransactionJSON{
		Id:              trn.Id,
		AccountId:       trn.AccountId,
		Label:           trn.Label,
		CategoryId:      trn.CategoryId,
		Description:     trn.Description,
		TransactionDate: trn.TransactionDate.UnixNano() / 1e6,
		RecordDate:      trn.RecordDate.UnixNano() / 1e6,
		Debit:           trn.Debit,
		Credit:          trn.Credit,
		Reference:       trn.Reference,
	}
}

// stringMapInputToDBTransactionParams process a map[string]interface{} input,
// normally received on the payload of a POST/PUT request, to create a
// DBTransactionParams object. if mandatory fields are not found, this function
// returns an error.
func stringMapInputToDBTransactionParams(
	// -- args --
	input map[string]interface{},
) (
	// -- returns --
	database.DBTransactionParams,
	error,
) {
	var res database.DBTransactionParams
	var field string

	// TODO check why can't cast directly to int
	// (don't ask me why. It just werks...)
	field = "accountId"
	accountIdStr, valid := input[field].(float64)
	if stringInArray(field, mandatory_transaction_json_fields) && !valid {
		return res, missingParameterError{field}
	}
	res.AccountId = int(accountIdStr)

	field = "label"
	res.Description, valid = input[field].(string)
	if stringInArray(field, mandatory_transaction_json_fields) && !valid {
		return res, missingParameterError{field}
	}

	field = "categoryId"
	categoryIdStr, valid := input[field].(float64)
	if stringInArray(field, mandatory_transaction_json_fields) && !valid {
		return res, missingParameterError{field}
	}
	res.CategoryId = int(categoryIdStr)

	field = "description"
	res.Description, valid = input[field].(string)
	if stringInArray(field, mandatory_transaction_json_fields) && !valid {
		return res, missingParameterError{field}
	}

	field = "transactionDate"
	tDateStr, valid := input[field].(float64)
	if stringInArray(field, mandatory_transaction_json_fields) && !valid {
		return res, missingParameterError{field}
	}
	res.TransactionDate = int64TimeStampToTime(int64(tDateStr))

	field = "recordDate"
	rDateStr, valid := input[field].(float64)
	if stringInArray(field, mandatory_transaction_json_fields) && !valid {
		return res, missingParameterError{field}
	}
	res.RecordDate = int64TimeStampToTime(int64(rDateStr))

	field = "debit"
	debitStr, valid := input[field].(float64)
	if stringInArray(field, mandatory_transaction_json_fields) && !valid {
		return res, missingParameterError{field}
	}
	res.Debit = float32(debitStr)

	field = "credit"
	creditStr, valid := input[field].(float64)
	if stringInArray(field, mandatory_transaction_json_fields) && !valid {
		return res, missingParameterError{field}
	}
	res.Credit = float32(creditStr)

	return res, nil
}
