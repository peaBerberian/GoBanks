package database

import "time"

// Interface that any implementation of GoBanks databases must implement
type GoBanksDataBase interface {
	// Free/close the db if needed
	Close() error

	// TODO separate user part?
	// Number of users stored in the database
	UserLength() (len int, err error)

	// Add a single user
	AddUser(DBUserParams) (DBUser, error)

	// Remove multiple users, based on filters
	RemoveUser(int) error

	// Update the attributes of multiple users, based on filters and field
	// names.
	UpdateUser(int, []string, DBUserParams) error

	// Get a single user base on its Id.
	GetUser(DBUserFilters, []string) (DBUser, error)

	// Add a single bank
	AddBank(DBBankParams) (DBBank, error)

	// Update the attributes of multiple banks, based on filters and field names.
	UpdateBanks(DBBankFilters, []string, DBBankParams) error

	// Remove multiple banks, based on filters
	RemoveBanks(DBBankFilters) error

	// Get multiple banks, based on filters
	// The second params are the wanted fields
	GetBanks(DBBankFilters, []string, uint) ([]DBBank, error)

	// Add a single account for a specific Id
	AddAccount(DBAccountParams) (DBAccount, error)

	// Update the attributes of multiple accounts,
	// based on filters and field names.
	UpdateAccounts(DBAccountFilters, []string, DBAccountParams) error

	// Remove multiple accounts, based on filters
	RemoveAccounts(DBAccountFilters) error

	// Get multiple accounts, based on filters
	// The second params are the wanted fields
	GetAccounts(DBAccountFilters, []string, uint) ([]DBAccount, error)

	// Add a single transaction for a specific Id
	AddTransaction(DBTransactionParams) (DBTransaction, error)

	// Update the attributes of multiple transactions,
	// based on filters and field names.
	UpdateTransactions(DBTransactionFilters, []string, DBTransactionParams) error

	// Remove multiple transactions, based on filters
	RemoveTransactions(DBTransactionFilters) error

	// Get multiple transactions, based on filters
	// The second params are the wanted fields
	GetTransactions(DBTransactionFilters, []string, uint) ([]DBTransaction, error)
}

type dbFilterInterface interface {
	isFilterActivated() bool
	getFilterValue() interface{}
}

type dbBaseFilter struct {
	Activated bool
}

func (d dbBaseFilter) isFilterActivated() bool { return d.Activated }

type DBGenericFilter struct {
	dbBaseFilter
	Value interface{}
}

type DBIntFilter struct {
	dbBaseFilter
	Value int
}

type DBIntArrayFilter struct {
	dbBaseFilter
	Value []int
}

type DBStringFilter struct {
	dbBaseFilter
	Value string
}

type DBStringArrayFilter struct {
	dbBaseFilter
	Value []string
}

type DBBoolFilter struct {
	dbBaseFilter
	Value bool
}

type DBFloatFilter struct {
	dbBaseFilter
	Value float32
}

type DBTimeFilter struct {
	dbBaseFilter
	Value time.Time
}

// About to get really ugly
func (d DBGenericFilter) getFilterValue() interface{}     { return d.Value }
func (d DBIntFilter) getFilterValue() interface{}         { return d.Value }
func (d DBIntArrayFilter) getFilterValue() interface{}    { return d.Value }
func (d DBStringFilter) getFilterValue() interface{}      { return d.Value }
func (d DBStringArrayFilter) getFilterValue() interface{} { return d.Value }
func (d DBBoolFilter) getFilterValue() interface{}        { return d.Value }
func (d DBFloatFilter) getFilterValue() interface{}       { return d.Value }
func (d DBTimeFilter) getFilterValue() interface{}        { return d.Value }

type DBAccountFilters struct {
	UserId  DBIntFilter
	Ids     DBIntArrayFilter
	BankIds DBIntArrayFilter
	Names   DBStringArrayFilter
}

type DBBankFilters struct {
	UserId DBIntFilter
	Ids    DBIntArrayFilter
	Names  DBStringArrayFilter
}

// TODO To implement
type DBTransactionCategoryFilters struct {
	UserId DBIntFilter
	Ids    DBIntArrayFilter
	Users  DBIntArrayFilter
	Names  DBStringArrayFilter
	Parent DBIntArrayFilter
}

type DBUserFilters struct {
	Id            DBIntFilter
	Name          DBStringFilter
	Administrator DBBoolFilter
}

type DBTransactionFilters struct {
	UserId              DBIntFilter
	Ids                 DBIntArrayFilter
	AccountIds          DBIntArrayFilter
	BankIds             DBIntArrayFilter
	CategoryIds         DBIntArrayFilter
	FromTransactionDate DBTimeFilter
	ToTransactionDate   DBTimeFilter
	FromRecordDate      DBTimeFilter
	ToRecordDate        DBTimeFilter
	MinDebit            DBFloatFilter
	MaxDebit            DBFloatFilter
	MinCredit           DBFloatFilter
	MaxCredit           DBFloatFilter
	SearchLabel         DBStringFilter
	SearchDescription   DBStringFilter
	SearchReference     DBStringFilter
}

type DBTransactionParams struct {
	// Id for the account concerned by this transaction
	AccountId int

	// Label describing the transaction
	Label string

	// Category of the transaction
	CategoryId string

	// Details on the transaction
	Description string

	// Date on which the transaction was done
	TransactionDate time.Time

	// Date on which the transaction was recorded by the bank
	RecordDate time.Time

	// Amount of money going out of your pocket
	Debit float32

	// Amount of money going in your pocket
	Credit float32

	// Bank Reference (id)
	Reference string
}

type DBTransaction struct {
	// Id of the transaction in the database
	// Set by the database's methods.
	// (starts at 1, 0 if not added to the database)
	Id int

	// Id for the account concerned by this transaction
	AccountId int

	// Label describing the transaction
	Label string

	// Category of the transaction
	CategoryId string

	// Details on the transaction
	Description string

	// Date on which the transaction was done
	TransactionDate time.Time

	// Date on which the transaction was recorded by the bank
	RecordDate time.Time

	// Amount of money going out of your pocket
	Debit float32

	// Amount of money going in your pocket
	Credit float32

	// Bank Reference (id)
	Reference string
}

type DBBankParams struct {
	UserId      int
	Name        string
	Description string
}

type DBBank struct {
	// Id of the bank in the database
	// Set by the database's methods.
	// (starts at 1, 0 if not added to the database)
	Id int

	// User linked to this Bank
	// LinkedUser User
	UserId int

	// Name of the bank
	Name string

	// Optional description
	Description string
}

// TODO To implement
type DBCategory struct {
	// Id of the category in the database
	// Set by the database's methods.
	// (starts at 1, 0 if not added to the database)
	Id int

	// User linked to this category
	LinkedUserId int

	// Name of the category
	Name string

	// Optional description
	Description string

	// Id of the parent category
	// 0 if none
	ParentId int
}

type DBAccount struct {
	// Id of the bank account in the database
	// Set by the database's methods.
	// (starts at 1, 0 if not added to the database)
	Id int

	// Bank Id linked to this account
	BankId int

	// Name of the bank account
	Name string

	// Optional description
	Description string
}

type DBAccountParams struct {
	// Bank Id linked to this account
	BankId int

	// Name of the bank account
	Name string

	// Optional description
	Description string
}

type DBUser struct {
	// Id for the user in the database
	// Set by the database's methods.
	// (starts at 1, 0 if not added to the database)
	Id int

	// User's Name
	Name string

	// Hash of the user's password
	PasswordHash string

	// Password's salt
	Salt string

	// True if the user is an administrator
	Administrator bool
}

type DBUserParams struct {
	// User's Name
	Name string

	// Hash of the user's password
	PasswordHash string

	// Password's salt
	Salt string

	// True if the user is an administrator
	Administrator bool
}
