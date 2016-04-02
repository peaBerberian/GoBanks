package types

// used on json.marshall for constructing the API response
type BankJSON struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// used on json.marshall for constructing the API response
type BankAccountJSON struct {
	Id     int    `json:"id"`
	BankId int    `json:"bankId"`
	Name   string `json:"name"`
}

// used on json.marshall for constructing the API response
type TransactionJSON struct {
	Id        int     `json:"id"`
	AccountId int     `json:"accountId"`
	Label     string  `json:"label"`
	Debit     float32 `json:"debit"`
	Credit    float32 `json:"credit"`
	Category  string  `json:"category"`
}

type CategoryJSON struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
