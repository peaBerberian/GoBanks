package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/peaberberian/GoBanks/database"
)

// Respond to the request with the given error message
func handleError(w http.ResponseWriter, err error) {

	fmt.Fprintf(w, generateErrorResponse(err))
}

// Respond to the request with a message indicating success
func handleSuccess(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{\"success\":true}")
}

// Respond to the request with a message indicating that the method
// ("GET"/"POST"/"PUT"/"DELETE") is not supported for the route wanted
func handleNotSupportedMethod(w http.ResponseWriter, method string) {
	http.Error(w, "Method "+method+
		" not supported for this route.", 405)
}

// generateErrorResponse generates a JSON ready to be sent describing
// the error given in argument. See ErrorJSON for more information.
func generateErrorResponse(err error) string {
	var errJson ErrorJSON

	if val, ok := err.(GoBanksError); ok {
		errJson.Code = val.ErrorCode()
		errJson.Error = val.Error()
	} else {
		errJson.Code = 0
		errJson.Error = err.Error()
	}

	errBytes, err := json.Marshal(errJson)
	if err != nil {
		return "{\"error\":\"internal error\",\"code\":0}"
	} else {
		return string(errBytes)
	}
}

// getApiRoute reads the url and just returns the API wanted.
func getApiRoute(url string) string {
	var paths = strings.Split(url, "/")
	if len(paths) < 3 {
		return ""
	}
	return paths[2]
}

// getApiId
// TODO
func getApiId(url string) (int, bool) {
	var paths = strings.Split(url, "/")
	if len(paths) < 4 {
		return 0, false
	}
	if val, err := strconv.Atoi(paths[3]); err == nil {
		return val, true
	}
	return 0, false
}

// read io.Reader (presumably from a request body) into a map[string]
// returns an error if:
//   - the reader could not be read (genericOperationError)
//   - the content could not be translated into a map[string] (bodyParsingError)
func readBodyAsStringMap(r io.Reader) (map[string]interface{}, error) {
	var xMap map[string]interface{}

	// parse the request's body
	bodyBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return xMap, genericOperationError{}
	}

	// get body as a map[string]interface{}
	xMap, err = parseJsonToStringMap(bodyBytes)
	if err != nil {
		return xMap, bodyParsingError{}
	}
	return xMap, nil
}

// read io.Reader (presumably from a request body) into a []map[string]
// returns an error if:
//   - the reader could not be read (genericOperationError)
//   - the content could not be translated into a []map[string] (bodyParsingError)
func readBodyAsArrayOfStringMap(r io.Reader) ([]map[string]interface{}, error) {
	var xMap []map[string]interface{}

	// parse the request's body
	bodyBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return xMap, genericOperationError{}
	}

	// get body as a map[string]interface{}
	xMap, err = parseJsonToArrayOfStringMap(bodyBytes)
	if err != nil {
		return xMap, bodyParsingError{}
	}
	return xMap, nil
}

// Parse JSON into a string map. Returns the unmarshalling error if one.
func parseJsonToStringMap(input []byte) (map[string]interface{}, error) {
	var inputJson map[string]interface{}
	var err = json.Unmarshal(input, &inputJson)
	return inputJson, err
}

// Parse JSON into an array of string map. Returns the unmarshalling error if one.
func parseJsonToArrayOfStringMap(input []byte) ([]map[string]interface{}, error) {
	var inputJson []map[string]interface{}
	var err = json.Unmarshal(input, &inputJson)
	return inputJson, err
}

// Returns true if an int was found in an array of int
// intInArray(4, []int{1,4} -> true
// intInArray(4, []int{1,3} -> false
func intInArray(val int, arr []int) bool {
	for _, av := range arr {
		if av == val {
			return true
		}
	}
	return false
}

// Returns true if an string was found in an array of string
// stringInArray(4, []string{1,4} -> true
// stringInArray(4, []string{1,3} -> false
func stringInArray(val string, arr []string) bool {
	for _, av := range arr {
		if av == val {
			return true
		}
	}
	return false
}

// Read a specific query string property and try to convert it into
// an array of string
// The returned boolean is false when the property content is empty.
func queryStringPropertyToStringArray(
	qs url.Values,
	str string,
) ([]string, bool) {
	if ctnt := qs.Get(str); ctnt != "" {
		var strArr []string
		return append(strArr, strings.Split(ctnt, ",")...), true
	}
	return []string{}, false
}

// Read a specific query string property and try to convert it into
// an array of int
// The returned boolean is false when the property content is empty.
func queryStringPropertyToIntArray(
	qs url.Values,
	str string,
) ([]int, bool) {
	if ctnt := qs.Get(str); ctnt != "" {
		strVals := strings.Split(ctnt, ",")
		var intArr []int
		for _, strVal := range strVals {
			if toInt, err := strconv.Atoi(strVal); err == nil {
				intArr = append(intArr, toInt)
			}
		}
		return intArr, true
	}
	return []int{}, false
}

// Read a specific query string property and try to convert it into
// a time.Time (consider ms timestamp as querystring values)
// The returned boolean is false when the property content is empty.
func queryStringPropertyToTime(qs url.Values, str string) (time.Time, bool) {
	toInt, worked := queryStringPropertyToInt(qs, str)
	if !worked {
		return time.Time{}, false
	}
	return int64TimeStampToTime(int64(toInt)), true
}

// Read a specific query string property and try to convert it into
// an int
// The returned boolean is false when the property content is empty.
func queryStringPropertyToInt(qs url.Values, str string) (int, bool) {
	if ctnt := qs.Get(str); ctnt != "" {
		if toInt, err := strconv.Atoi(ctnt); err == nil {
			return toInt, true
		}
	}
	return 0, false
}

// Read a specific query string property and try to convert it into
// a float32
// The returned boolean is false when the property content is empty.
func queryStringPropertyToFloat32(qs url.Values, str string) (float32, bool) {
	if ctnt := qs.Get(str); ctnt != "" {
		if toInt, err := strconv.ParseFloat(ctnt, 32); err == nil {
			return float32(toInt), true
		}
	}
	return 0, false
}

func int64TimeStampToTime(ts int64) time.Time {
	return time.Unix(0, ts*1e6)
}

func getBankIdsForUserId(userId int) ([]int, error) {
	var banksFilter database.DBBankFilters
	banksFilter.UserId.SetFilter(userId)
	bnks, err := database.GoDB.GetBanks(banksFilter, []string{"Id"}, 0)
	if err != nil {
		return []int{}, err
	}
	var bnkIds []int
	for _, bnk := range bnks {
		bnkIds = append(bnkIds, bnk.Id)
	}
	return bnkIds, nil
}

func getAccountIdsForBankIds(bankIds []int) ([]int, error) {

	var accountsFilter database.DBAccountFilters
	accountsFilter.BankIds.SetFilter(bankIds)
	accs, err := database.GoDB.GetAccounts(accountsFilter, []string{"Id"}, 0)
	if err != nil {
		return []int{}, err
	}
	var accIds []int
	for _, acc := range accs {
		accIds = append(accIds, acc.Id)
	}
	return accIds, nil
}

func getTransactionIdsForAccountIds(accountIds []int) ([]int, error) {

	var transactionsFilter database.DBTransactionFilters
	transactionsFilter.AccountIds.SetFilter(accountIds)
	accs, err := database.GoDB.GetTransactions(transactionsFilter, []string{"Id"}, 0)
	if err != nil {
		return []int{}, err
	}
	var accIds []int
	for _, acc := range accs {
		accIds = append(accIds, acc.Id)
	}
	return accIds, nil
}
