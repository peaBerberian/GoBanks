package api

import "net/http"

import "encoding/json"

// import "fmt"
import "strings"
import "strconv"
import "time"

import "github.com/peaberberian/GoBanks/auth"
import "github.com/peaberberian/GoBanks/database"

func handleTransactions(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	var authorizedBankIds []int
	var authorizedAccountIds []int
	var err error

	authorizedBankIds, err = getBankIdsForUserId(t.UserId)
	if err != nil {
		handleError(w, err)
		return
	}

	authorizedAccountIds, err = getAccountIdsForBankIds(authorizedBankIds)
	if err != nil {
		handleError(w, err)
		return
	}

	var ft database.DBTransactionFilters
	var queryString = r.URL.Query()
	var intArray []int

	// if only some accounts are wanted, construct filter
	if val := queryString.Get("accounts"); val != "" {
		accIdsStr := strings.Split(val, ",")
		for _, accIdStr := range accIdsStr {
			var id, err = strconv.Atoi(accIdStr)
			if err == nil && intInArray(id, authorizedAccountIds) {
				intArray = append(intArray, id)
			}
		}
		if len(intArray) > 0 {
			ft.AccountIds.Activated = true
			ft.AccountIds.Value = intArray
		}
	} else {
		ft.AccountIds.Activated = true
		ft.AccountIds.Value = authorizedAccountIds
	}

	// vals, err := database.GoDB.GetTransactions(ft)
	// if err != nil {
	// 	handleError(w, err)
	// 	return
	// }

	// if len(vals) == 0 {
	// 	fmt.Fprintf(w, "[]")
	// 	return
	// }

	// fmt.Fprintf(w, generateTransactionsResponse(vals))
}

func generateTransactionsResponse(ts []database.DBTransaction) string {
	var resJson []TransactionJSON
	for _, t := range ts {
		resJson = append(resJson, TransactionJSON{
			Id:              t.Id,
			AccountId:       t.AccountId,
			Label:           t.Label,
			Debit:           t.Debit,
			Credit:          t.Credit,
			Description:     t.Description,
			CategoryId:      t.CategoryId,
			TransactionDate: t.TransactionDate.UnixNano() / 1e6,
			RecordDate:      t.RecordDate.UnixNano() / 1e6,
		})
	}
	resBytes, err := json.Marshal(resJson)
	if err != nil || resBytes == nil {
		return "[]"
	}
	return string(resBytes)
}

func parseTransactionJson(input string) (database.DBTransactionParams, error) {
	var t database.DBTransactionParams
	var inputJson TransactionJSON
	var err = json.Unmarshal([]byte(input), inputJson)
	if err != nil {
		return t, err
	}

	var tDate time.Time
	var rDate time.Time
	tDate = time.Unix(0, inputJson.TransactionDate*1e6)
	rDate = time.Unix(0, inputJson.RecordDate*1e6)

	t = database.DBTransactionParams{
		AccountId:       inputJson.AccountId,
		Label:           inputJson.Label,
		Debit:           inputJson.Debit,
		Credit:          inputJson.Credit,
		Description:     inputJson.Description,
		CategoryId:      inputJson.CategoryId,
		TransactionDate: tDate,
		RecordDate:      rDate,
	}
	return t, nil
}
