package database

// table used
const user_table = "user"

// fields used, in the right order
// TODO only []string?
var user_fields = map[string]string{
	"Id":           "id",
	"Name":         "name",
	"PasswordHash": "password",
	"Salt":         "salt",
}

const bank_table = "bank"

var bank_fields = map[string]string{
	"Id":          "id",
	"UserId":      "user_id",
	"Name":        "name",
	"Description": "description",
}

const account_table = "account"

var account_fields = map[string]string{
	"Id":          "id",
	"BankId":      "bank_id",
	"Name":        "name",
	"Description": "description",
}

const category_table = "category"

var category_fields = map[string]string{
	"Id":          "id",
	"UserId":      "user_id",
	"Name":        "name",
	"Description": "description",
	"ParentId":    "parent_id",
}

const transaction_table = "transaction"

var transaction_fields = map[string]string{
	"Id":              "id",
	"AccountId":       "account_id",
	"Label":           "label",
	"CategoryId":      "category_id",
	"Description":     "description",
	"TransactionDate": "transaction_date",
	"RecordDate":      "record_date",
	"Debit":           "debit",
	"Credit":          "credit",
	"Reference":       "reference",
}
