package api

// used on json.marshall for constructing the API response
type UserJSON struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// used on json.marshall for constructing the API response
type BankJSON struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// used on json.marshall for constructing the API response
type AccountJSON struct {
	Id          int    `json:"id"`
	BankId      int    `json:"bankId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// used on json.marshall for constructing the API response
type TransactionJSON struct {
	Id              int     `json:"id"`
	AccountId       int     `json:"accountId"`
	Description     string  `json:"description"`
	Label           string  `json:"label"`
	Debit           float32 `json:"debit"`
	Credit          float32 `json:"credit"`
	CategoryId      string  `json:"category"`
	TransactionDate int64   `json:"transactionDate"`
	RecordDate      int64   `json:"recordDate"`
}

type CategoryJSON struct {
	Id     int           `json:"id"`
	Name   string        `json:"name"`
	Childs *CategoryJSON `json"childs"`
}

type TokenJSON struct {
	Token     string `json:"access_token"`
	TokenType string `json:"token_type"`
	Expires   int    `json:"expires_in"`
}

type AuthenticationJSON struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type ErrorJSON struct {
	Error string `json:"error"`
	Code  uint32 `json:"code"`
}

type GoBanksError interface {
	error
	ErrorCode() uint32
}
