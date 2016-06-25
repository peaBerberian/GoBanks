package api

import (
	"encoding/json"
	"fmt"
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

// handleBanks is the main handler for call on the /banks api. It dispatches
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

// handleBankRead handle GET requests on the /banks API
func handleBankRead(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	// look if we have an id (GET /banks/35 => id == 35)
	var id, hasIdInUrl = getApiId(r.URL.Path)

	var queryString = r.URL.Query()
	var f database.DBBankFilters
	var limit int

	// always filter on the current user
	f.UserId.SetFilter(t.UserId)

	// if an id was set in the url, filter to the record corresponding to it
	if hasIdInUrl {
		f.Ids.SetFilter([]int{id})
	} else {
		// if only some bank ids are wanted, filter
		wantedIds, _ := queryStringPropertyToIntArray(queryString, "id")
		fmt.Println(wantedIds)
		if len(wantedIds) > 0 {
			f.Ids.SetFilter(wantedIds)
		}

		// if only some bank names are wanted, filter
		wantedBankNames, _ := queryStringPropertyToStringArray(queryString, "name")
		if len(wantedBankNames) > 0 {
			f.Names.SetFilter(wantedBankNames)
		}

		// obtain limit of wanted records, if set
		limit, _ = queryStringPropertyToInt(queryString, "limit")
	}

	// perform the database request
	vals, err := database.GoDB.GetBanks(f, gettable_bank_fields, uint(limit))
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	// if an id was given, we're awaiting an object, not an array.
	if hasIdInUrl {
		if len(vals) == 0 {
			fmt.Fprintf(w, "{}")
		} else {
			fmt.Fprintf(w, generateBankResponse(vals[0]))
		}
		return
	}

	// else respond directly with the result
	if len(vals) == 0 {
		fmt.Fprintf(w, "[]")
	} else {
		fmt.Fprintf(w, generateBanksResponse(vals))
	}
}

// handleBankCreate handle POST requests on the /banks API
func handleBankCreate(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	// you cannot post on a specific id, reject if you want to do that
	if _, hasIdInUrl := getApiId(r.URL.Path); hasIdInUrl {
		handleNotSupportedMethod(w, r.Method)
		return
	}

	// convert body to map[string]interface{}
	bodyMap, err := readBodyAsStringMap(r.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	// translate data into a DBBankParams element
	// (also check mandatory fields)
	bankElem, err := inputToBankParams(bodyMap)
	if err != nil {
		handleError(w, err)
		return
	}

	// attach elem to current user
	bankElem.UserId = t.UserId

	// perform database add request
	bank, err := database.GoDB.AddBank(bankElem)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	fmt.Fprintf(w, generateBankResponse(bank))
}

// handleBankUpdate handle PUT requests on the /banks API
func handleBankUpdate(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	// look if we have an id (GET /banks/35 => id == 35)
	var id, hasIdInUrl = getApiId(r.URL.Path)

	// if an id was found, it means that we want to replace an element
	// redirect to the right function
	if !hasIdInUrl {
		handleBankReplace(w, r, t)
		return
	}

	// check that we can modify this bank params
	// (blocking database request here :(, TODO see what I can do, jwt?)
	if err := checkPermissionForBank(t, id); err != nil {
		handleError(w, err)
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
	var bankElem database.DBBankParams

	if val, ok := bodyMap["name"]; ok {
		if str, ok := val.(string); !ok {
			handleError(w, bodyParsingError{})
			return
		} else {
			bankElem.Name = str
			fields = append(fields, "Name")
		}
	}
	if val, ok := bodyMap["description"]; ok {
		if str, ok := val.(string); !ok {
			handleError(w, bodyParsingError{})
			return
		} else {
			bankElem.Description = str
			fields = append(fields, "Description")
		}
	}

	// Filter the bank id
	var f database.DBBankFilters
	f.Ids.SetFilter([]int{id})

	// perform the database request
	if err = database.GoDB.UpdateBanks(f, fields, bankElem); err != nil {
		handleError(w, queryOperationError{})
		return
	}

	handleSuccess(w, r)
}

// handleBankDelete handle DELETE requests on the /banks API
func handleBankDelete(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	// look if we have an id (GET /banks/35 => id == 35)
	var id, hasIdInUrl = getApiId(r.URL.Path)

	var f database.DBBankFilters

	// if we have an id, check permission and set filter
	if hasIdInUrl {
		// (blocking database request here :(, TODO see what I can do, jwt?)
		if err := checkPermissionForBank(t, id); err != nil {
			handleError(w, err)
			return
		}
		f.Ids.SetFilter([]int{id})
	} else {
		// filter by userId
		f.UserId.SetFilter(t.UserId)
	}

	// perform the database request
	if err := database.GoDB.RemoveBanks(f); err != nil {
		handleError(w, queryOperationError{})
		return
	}
	handleSuccess(w, r)
}

// handleBankReplace handle specifically PUT requests on the main /banks API
// (not restricted to a certain id).
func handleBankReplace(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	var bodyMaps, err = readBodyAsArrayOfStringMap(r.Body)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	var bnks []database.DBBankParams

	// translate data into DBBankParams elements
	// (also check mandatory fields)
	for _, bodyMap := range bodyMaps {
		bankElem, err := inputToBankParams(bodyMap)
		if err != nil {
			handleError(w, err)
			return
		}
		bnks = append(bnks, bankElem)
	}

	// Remove old banks linked to this user
	var f database.DBBankFilters
	f.UserId.SetFilter(t.UserId)

	if err := database.GoDB.RemoveBanks(f); err != nil {
		handleError(w, queryOperationError{})
		return
	}

	// add each bank indicated to the database
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

// dbBankToBankJSON takes a DBBank and convert it to its corresponding
// BankJSON response.
func dbBankToBankJSON(bnk database.DBBank) BankJSON {
	return BankJSON{
		Id:          bnk.Id,
		Name:        bnk.Name,
		Description: bnk.Description,
	}
}

// process map[string]interface{} input to create a DBBankParams object.
// if mandatory fields are not found, this function returns an error.
func inputToBankParams(input map[string]interface{}) (database.DBBankParams, error) {
	var res database.DBBankParams
	var valid bool

	// The "name" field is mandatory
	res.Name, valid = input["name"].(string)
	if !valid {
		return res, missingParameterError{"name"}
	}

	res.Description, _ = input["description"].(string)

	return res, nil
}

// // readBodyAsBankParamsArray generates an array of DBBankParams structs from the
// // body passed with the given http request.
// func readBodyAsBankParamsArray(r io.Reader) ([]database.DBBankParams, error) {
// 	bodyBytes, err := ioutil.ReadAll(r)
// 	if err != nil {
// 		return []database.DBBankParams{}, err
// 	}
// 	return unmarshallToDBBankParamsArray(bodyBytes)
// }

// // unmarshallToDBBankParams unmarshall an array of bytes into a
// // DBBankParams struct
// func unmarshallToDBBankParams(input []byte) (database.DBBankParams, error) {
// 	var inputJson BankJSON
// 	if err := json.Unmarshal([]byte(input), &inputJson); err != nil {
// 		return database.DBBankParams{}, err
// 	}
// 	return bankJSONToDBBankParams(inputJson), nil
// }

// // unmarshallToDBBankParamsArray unmarshall an array of bytes into an array
// // of DBBankParams structs.
// func unmarshallToDBBankParamsArray(input []byte) ([]database.DBBankParams,
// 	error) {
// 	var res []database.DBBankParams
// 	var inputJson []BankJSON
// 	if err := json.Unmarshal([]byte(input), &inputJson); err != nil {
// 		return []database.DBBankParams{}, err
// 	}
// 	for _, val := range inputJson {
// 		res = append(res, bankJSONToDBBankParams(val))
// 	}
// 	return res, nil
// }

// // bankJSONToDBBankParams takes a BankJSON and convert it to its
// // corresponding DBBankParams struct (used for the database).
// func bankJSONToDBBankParams(bnkj BankJSON) database.DBBankParams {
// 	return database.DBBankParams{
// 		Name:        bnkj.Name,
// 		Description: bnkj.Description,
// 	}
// }
