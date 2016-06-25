package database

import "time"

// Perform operations on the DataBase relative to Users
type UserDataBase interface {
	// Number of users stored in the database
	UserLength() (len int, err error)

	// Add a single user
	AddUser(DBUserParams) (DBUser, error)

	// Remove multiple users, based on filters
	RemoveUser(int) error

	// Update the attributes of multiple users, based on filters and field
	// names.
	UpdateUser(int, []string, DBUserParams) error

	// Get a single user based on filters.
	// The second param is the wanted field
	GetUser(DBUserFilters, []string) (DBUser, error)
}

// Perform operations on the DataBase relative to Categories
type CategoryDataBase interface {
	// Add a single Category
	AddCategory(DBCategoryParams) (DBCategory, error)

	// Remove multiple Categories, based on filters
	UpdateCategories(DBCategoryFilters, []string, DBCategoryParams) error

	// Update the attributes of multiple categories, based on filters and field
	// names.
	RemoveCategories(DBCategoryFilters) error

	// Get a multiple categories based on filters.
	// The second param is  the wanted fields
	GetCategories(DBCategoryFilters, []string, uint) ([]DBCategory, error)
}

// Perform operations on the DataBase relative to Bank Accounts
type BankAccountDataBase interface {
	// Add a single account for a specific Id
	AddAccount(DBAccountParams) (DBAccount, error)

	// Update the attributes of multiple accounts,
	// based on filters and field names.
	UpdateAccounts(DBAccountFilters, []string, DBAccountParams) error

	// Remove multiple accounts, based on filters
	RemoveAccounts(DBAccountFilters) error

	// Get multiple accounts, based on filters
	// The second param is  the wanted fields
	GetAccounts(DBAccountFilters, []string, uint) ([]DBAccount, error)
}

// Perform operations on the DataBase relative to Bank Accounts
type BankDatabase interface {
	// Add a single bank
	AddBank(DBBankParams) (DBBank, error)

	// Update the attributes of multiple banks, based on filters and field names.
	UpdateBanks(DBBankFilters, []string, DBBankParams) error

	// Remove multiple banks, based on filters
	RemoveBanks(DBBankFilters) error

	// Get multiple banks, based on filters
	// The second param is  the wanted fields
	GetBanks(DBBankFilters, []string, uint) ([]DBBank, error)
}

// Perform operations on the DataBase relative to Transactions
type TransactionDataBase interface {
	// Add a single transaction for a specific Id
	AddTransaction(DBTransactionParams) (DBTransaction, error)

	// Update the attributes of multiple transactions,
	// based on filters and field names.
	UpdateTransactions(DBTransactionFilters, []string, DBTransactionParams) error

	// Remove multiple transactions, based on filters
	RemoveTransactions(DBTransactionFilters) error

	// Get multiple transactions, based on filters
	// The second param is  the wanted fields
	GetTransactions(DBTransactionFilters, []string, uint) ([]DBTransaction, error)
}

// Interface that any implementation of GoBanks databases must implement
type GoBanksDataBase interface {
	// Free/close the db if needed
	Close() error

	UserDataBase
	CategoryDataBase
	BankAccountDataBase
	BankDatabase
	TransactionDataBase
}

// Representation of a single User as returned by the UserDatabase
type DBUser struct {
	// Id for the user in the database
	// Set by the database's methods.
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

// Representation of a single Category as returned by the CategoryDatabase
type DBCategory struct {
	// Id of the category in the database
	// Set by the database's methods.
	Id int

	// User linked to this category
	UserId int

	// Name of the category
	Name string

	// Optional description
	Description string

	// Id of the parent category
	// 0 if none
	ParentId int
}

// Representation of a single Account as returned by the BankAccountDatabase
type DBAccount struct {
	// Id of the bank account in the database
	// Set by the database's methods.
	Id int

	// Bank Id linked to this account
	BankId int

	// Name of the bank account
	Name string

	// Optional description
	Description string
}

// Representation of a single Bank as returned by the BankDatabase
type DBBank struct {
	// Id of the bank in the database
	// Set by the database's methods.
	Id int

	// User linked to this Bank
	// LinkedUser User
	UserId int

	// Name of the bank
	Name string

	// Optional description
	Description string
}

// Representation of a single Transaction as returned by the TransactionDatabase
type DBTransaction struct {
	// Id of the transaction in the database
	// Set by the database's methods.
	Id int

	// Id for the account concerned by this transaction
	AccountId int

	// Label describing the transaction
	Label string

	// Category of the transaction
	CategoryId int

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

// Parameters awaited to create a new User in the UserDatabase
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

// Parameters awaited to create a new Category in the CategoryDatabase
type DBCategoryParams struct {
	// The user adding the category
	UserId int

	// The 'name' of the category
	Name string

	// Optional description
	Description string

	// Id of the category's parent, for nesting
	ParentId int
}

// Parameters awaited to create a new BankAccount in the BankAccountDatabase
type DBAccountParams struct {
	// Bank Id linked to this account
	BankId int

	// Name of the bank account
	Name string

	// Optional description
	Description string
}

// Parameters awaited to create a new Bank in the BankDatabase
type DBBankParams struct {
	// The user adding the bank
	UserId int

	// The 'name' of the bank
	Name string

	// Optional description
	Description string
}

// Parameters awaited to create a new Transaction in the TransactionDatabase
type DBTransactionParams struct {
	// Id for the account concerned by this transaction
	AccountId int

	// Label describing the transaction
	Label string

	// Category of the transaction
	CategoryId int

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

// Filters that can be used to filter Users when doing operations on the
// UserDatabase
// example: filters.Id.SetValue(5)
type DBUserFilters struct {
	// Filter by using the User Id
	Id DBIntFilter

	// Filter by using the user's name
	Name DBStringFilter

	// Filters only Administrators TODO what? No!
	Administrator DBBoolFilter
}

// Filters that can be used to filter Categories when doing operations on the
// CategoryDatabase
// example: filters.Ids.SetValue([]int{5})
type DBCategoryFilters struct {
	// Filter by using the Categories Ids
	Ids DBIntArrayFilter

	// Filter by using the Categories names
	Names DBStringArrayFilter

	// Filter by using the User Id
	UserId DBIntFilter

	// Filter by using the Parent Categories Ids
	ParentIds DBIntArrayFilter
}

// Filters that can be used to filter Bank Accounts when doing operations on the
// BankAccountDataBase
// example: filters.Ids.SetValue([]int{5})
type DBAccountFilters struct {
	// Filter by using the Bank Account Ids
	Ids DBIntArrayFilter

	// Filter by using the User Id
	// TODO REMOVE?
	UserId DBIntFilter

	// Filter by using the Bank Ids corresponding to the accounts
	BankIds DBIntArrayFilter

	// Filter by using the Bank Account names
	Names DBStringArrayFilter
}

// Filters that can be used to filter Banks when doing operations on the
// BankDataBase
// example: filters.Ids.SetValue([]int{5})
type DBBankFilters struct {
	// Filter by using the Bank Ids
	Ids DBIntArrayFilter

	// Filter by using the User Id
	UserId DBIntFilter

	// Filter by using the Bank names
	Names DBStringArrayFilter
}

// Filters that can be used to filter Transactions when doing operations on the
// TransactionDataBase
// example: filters.Ids.SetValue([]int{5})
type DBTransactionFilters struct {
	// Filter by using the Transactions Ids
	Ids DBIntArrayFilter

	// Filter by using the User Id
	// TODO REMOVE?
	UserId DBIntFilter

	// Filter by using the Bank Accounts Ids corresponding to the transactions
	AccountIds DBIntArrayFilter

	// Filter by using the Bank Ids corresponding to the transactions
	// TODO REMOVE?
	BankIds DBIntArrayFilter

	// Filter by using the Categories Ids corresponding to the transactions
	CategoryIds DBIntArrayFilter

	// Filter by setting the minimum transaction date
	FromTransactionDate DBTimeFilter

	// Filter by setting the maximum transaction date
	ToTransactionDate DBTimeFilter

	// Filter by setting the minimum record date
	FromRecordDate DBTimeFilter

	// Filter by setting the maximum record date
	ToRecordDate DBTimeFilter

	// Filter by setting the minimum debit
	MinDebit DBFloatFilter

	// Filter by setting the maximum debit
	MaxDebit DBFloatFilter

	// Filter by setting the minimum credit
	MinCredit DBFloatFilter

	// Filter by setting the maximum credit
	MaxCredit DBFloatFilter

	// Filter by setting the bank's reference
	References DBStringArrayFilter
}

// Common base of filters
type dbBaseFilter struct {
	activated bool
}

func (d dbBaseFilter) isFilterActivated() bool { return d.activated }

// Filter by setting an int value
type DBIntFilter struct {
	dbBaseFilter
	value int
}

// Filter by setting a []int value
type DBIntArrayFilter struct {
	dbBaseFilter
	value []int
}

// Filter by setting a string value
type DBStringFilter struct {
	dbBaseFilter
	value string
}

// Filter by setting a []string value
type DBStringArrayFilter struct {
	dbBaseFilter
	value []string
}

// Filter by setting a bool value
type DBBoolFilter struct {
	dbBaseFilter
	value bool
}

// Filter by setting a float32 value
type DBFloatFilter struct {
	dbBaseFilter
	value float32
}

// Filter by setting a time.Time value
type DBTimeFilter struct {
	dbBaseFilter
	value time.Time
}

// About to get really ugly

// Activate and set the value for a DBIntFilter
func (d *DBIntFilter) SetFilter(val int) {
	d.activated = true
	d.value = val
}

// Activate and set the value for a DBIntArrayFilter
func (d *DBIntArrayFilter) SetFilter(val []int) {
	d.activated = true
	d.value = val
}

// Activate and set the value for a DBStringFilter
func (d *DBStringFilter) SetFilter(val string) {
	d.activated = true
	d.value = val
}

// Activate and set the value for a DBStringArrayFilter
func (d *DBStringArrayFilter) SetFilter(val []string) {
	d.activated = true
	d.value = val
}

// Activate and set the value for a DBBoolFilter
func (d *DBBoolFilter) SetFilter(val bool) {
	d.activated = true
	d.value = val
}

// Activate and set the value for a DBFloatFilter
func (d *DBFloatFilter) SetFilter(val float32) {
	d.activated = true
	d.value = val
}

// Activate and set the value for a DBTimeFilter
func (d *DBTimeFilter) SetFilter(val time.Time) {
	d.activated = true
	d.value = val
}

// ugly generic dbFilter interface for using them in generic helpers
type dbFilterInterface interface {
	isFilterActivated() bool
	getFilterValue() interface{}
}

func (d DBIntFilter) getFilterValue() interface{}         { return d.value }
func (d DBIntArrayFilter) getFilterValue() interface{}    { return d.value }
func (d DBStringFilter) getFilterValue() interface{}      { return d.value }
func (d DBStringArrayFilter) getFilterValue() interface{} { return d.value }
func (d DBBoolFilter) getFilterValue() interface{}        { return d.value }
func (d DBFloatFilter) getFilterValue() interface{}       { return d.value }
func (d DBTimeFilter) getFilterValue() interface{}        { return d.value }
