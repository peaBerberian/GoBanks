package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/peaberberian/GoBanks/database"
)

func httpRequestToMap(r *http.Request) (map[string]interface{}, error) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return make(map[string]interface{}), err
	}
	return parseJson(bodyBytes)
}

func parseJson(input []byte) (map[string]interface{}, error) {
	var inputJson map[string]interface{}
	var err = json.Unmarshal(input, &inputJson)
	return inputJson, err
}

func getBankIdsForUserId(userId int) ([]int, error) {
	var banksFilter database.DBBankFilters
	banksFilter.UserId.Activated = true
	banksFilter.UserId.Value = userId
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
	accountsFilter.BankIds.Activated = true
	accountsFilter.BankIds.Value = bankIds
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

func handleError(w http.ResponseWriter, err error) {

	fmt.Fprintf(w, generateErrorResponse(err))
}

func handleSuccess(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{\"success\":true}")
}

func handleNotSupportedMethod(w http.ResponseWriter, method string) {
	http.Error(w, "Method "+method+
		" not supported for this route.", 405)
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

func intInArray(val int, arr []int) bool {
	for _, av := range arr {
		if av == val {
			return true
		}
	}
	return false
}

func checkMandatoryFields(r *http.Request, fields []string) error {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return genericOperationError{}
	}

	var xMap map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &xMap); err != nil {
		return bodyParsingError{}
	}
	for _, field := range fields {
		if xMap[field] == nil {
			return missingParameterError{field}
		}
	}
	return nil
}

func getQueryStringLimit(qs url.Values) uint {
	if val := qs.Get("limit"); val != "" {
		if valInt, err := strconv.Atoi(val); err != nil {
			return uint(valInt)
		}
	}
	return 0
}

func addIntArrayFilter(x []int, f *database.DBIntArrayFilter) {
	f.Activated = true
	f.Value = x
}

func addIntFilter(x int, f *database.DBIntFilter) {
	f.Activated = true
	fmt.Println(x)
	f.Value = x
	fmt.Printf("%+v", f)
}

func queryStringToStringArrayFilter(qs url.Values, str string,
	f *database.DBStringArrayFilter) {
	if ctnt := qs.Get(str); ctnt != "" {
		var strArr []string
		f.Activated = true
		f.Value = append(strArr, strings.Split(ctnt, ",")...)
	}
}

func queryStringToIntArrayFilter(qs url.Values, str string,
	f *database.DBIntArrayFilter) {
	if ctnt := qs.Get(str); ctnt != "" {
		strVals := strings.Split(ctnt, ",")
		var intArr []int
		for _, strVal := range strVals {
			if toInt, err := strconv.Atoi(strVal); err != nil {
				intArr = append(intArr, toInt)
			}
			if len(intArr) > 0 {
				f.Activated = true
				f.Value = intArr
			}
		}
	}
}

func queryStringToIntFilter(qs url.Values, str string,
	f *database.DBIntFilter) {
	if ctnt := qs.Get(str); ctnt != "" {
		if toInt, err := strconv.Atoi(ctnt); err == nil {
			addIntFilter(toInt, f)
		}
	}
}
