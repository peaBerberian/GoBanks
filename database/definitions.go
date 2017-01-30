package database

import "time"

// Perform operations on the DataBase relative to Users
// TODO Allow only one user? Removes a lot of complexity
type UserDataBase interface {
	// Number of users stored in the database
	UserLength() (len int, err error)

	// Add a single user
	AddUser(DBUserParams) (DBUser, error)

	// Remove multiple users, based on UserId
	RemoveUser(int) error

	// Update the attributes of multiple users, based on UserId and params
	UpdateUser(int, DBUserParams) error

	// Get a single user based on filters.
	GetUser(DBUserFilters) (DBUser, error)
}

// Perform operations on the DataBase relative to Categories
type CategoryDataBase interface {
	// Add a single Category
	AddCategory(DBCategoryParams) (DBCategory, error)

	// Remove multiple Categories, based on filters
	UpdateCategories(DBCategoryFilters, DBCategoryParams) error

	// Update the attributes of multiple categories, based on filters and field
	// names.
	RemoveCategories(DBCategoryFilters) error

	// Get a multiple categories based on filters.
	// The second param is the max number of item you wish to receive (0 = no limit)
	GetCategories(DBCategoryFilters, uint) ([]DBCategory, error)
}

// Perform operations on the DataBase relative to Bank Accounts
type BankAccountDataBase interface {
	// Add a single account for a specific Id
	AddAccount(DBAccountParams) (DBAccount, error)

	// Update the attributes of multiple accounts,
	// based on filters and field names.
	UpdateAccounts(DBAccountFilters, DBAccountParams) error

	// Remove multiple accounts, based on filters
	RemoveAccounts(DBAccountFilters) error

	// Get multiple accounts, based on filters
	// The second param is the max number of item you wish to receive (0 = no limit)
	GetAccounts(DBAccountFilters, uint) ([]DBAccount, error)
}

// Perform operations on the DataBase relative to Bank Accounts
type BankDatabase interface {
	// Add a single bank
	AddBank(DBBankParams) (DBBank, error)

	// Update the attributes of multiple banks, based on filters and field names.
	UpdateBanks(DBBankFilters, DBBankParams) error

	// Remove multiple banks, based on filters
	RemoveBanks(DBBankFilters) error

	// Get multiple banks, based on filters
	// The second param is the max number of item you wish to receive (0 = no limit)
	GetBanks(DBBankFilters, uint) ([]DBBank, error)
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
	// The second param is the max number of item you wish to receive (0 = no limit)
	GetTransactions(DBTransactionFilters, uint) ([]DBTransaction, error)

	// Get the sum of all debits for the given filters
	// GetDebit(DBTransactionFilters)

	// Get the sum of all credits for the given filters
	// GetCredit(DBTransactionFilters)

	// Get Report for the given filters (debit and credit)
	// GetReport(DBTransactionFilters)
}

// Interface GoBanks databases must implement
type GoBanksDataBase interface {
	Close() error // Free/close the db if needed
	UserDataBase
	CategoryDataBase
	BankAccountDataBase
	BankDatabase
	TransactionDataBase
}

// Representation of a single User as returned by the UserDatabase
type DBUser struct {
	Id           int    // Id for the user in the database
	Name         string // User's Name
	PasswordHash string // Hash of the user's password
	Salt         string // Password's salt
}

// Representation of a single Category as returned by the CategoryDatabase
type DBCategory struct {
	Id          int    // Id of the category in the database
	UserId      int    // User linked to this category
	Name        string // Name of the category
	Description string // Optional description
	ParentId    int    // Id of the parent category, 0 if none (TODO remove)
}

// Representation of a single Account as returned by the BankAccountDatabase
type DBAccount struct {
	Id          int    // Id of the bank account in the database
	BankId      int    // Bank Id linked to this account
	Name        string // Name of the bank account
	Description string // Optional description
}

// Representation of a single Bank as returned by the BankDatabase
type DBBank struct {
	Id          int    // Id of the bank in the database
	UserId      int    // User linked to this Bank
	Name        string // Name of the bank
	Description string // Optional description
}

// Representation of a single Transaction as returned by the TransactionDatabase
type DBTransaction struct {
	Id              int       // Id of the transaction in the database
	AccountId       int       // Id for the account concerned by this transaction
	Label           string    // Label describing the transaction
	CategoryId      int       // Category of the transaction
	Description     string    // Details on the transaction
	TransactionDate time.Time // Date at which the transaction was done
	RecordDate      time.Time // Date at which the transaction was recorded
	Debit           float32   // Amount of money going out of your pocket
	Credit          float32   // Amount of money going in your pocket
	Reference       string    // Bank Reference (id)
}

// --------- PARAMS

// Parameters awaited to create a new User in the UserDatabase
type DBUserParams struct {
	Name         DBStringParam // User's Name
	PasswordHash DBStringParam // Hash of the user's password
	Salt         DBStringParam // Password's salt
}

// Parameters awaited to create a new Category in the CategoryDatabase
type DBCategoryParams struct {
	UserId      DBIntParam    // The user adding the category
	Name        DBStringParam // The 'name' of the category
	Description DBStringParam // Optional description
	ParentId    DBIntParam    // Id of the category's parent, for nesting (TODO Remove)
}

// Parameters awaited to create a new BankAccount in the BankAccountDatabase
type DBAccountParams struct {
	BankId      DBIntParam    // Bank Id linked to this account
	Name        DBStringParam // Name of the bank account
	Description DBStringParam // Optional description
}

// Parameters awaited to create a new Bank in the BankDatabase
type DBBankParams struct {
	UserId      DBIntParam    // The user adding the bank
	Name        DBStringParam // The 'name' of the bank
	Description DBStringParam // Optional description
}

// Parameters awaited to create a new Transaction in the TransactionDatabase
type DBTransactionParams struct {
	AccountId       DBIntParam    // Id for the account concerned by this transaction
	Label           DBStringParam // Label describing the transaction
	CategoryId      DBIntParam    // Category of the transaction
	Description     DBStringParam // Details on the transaction
	TransactionDate DBTimeParam   // Date on which the transaction was done
	RecordDate      DBTimeParam   // Date on which the transaction was recorded
	Debit           DBFloatParam  // Amount of money going out of your pocket
	Credit          DBFloatParam  // Amount of money going in your pocket
	Reference       DBStringParam // Bank Reference (id)
}

// ------ FILTERS

// Filters that can be used to filter Users when doing operations on the
// UserDatabase.
//
// example:
// ```go
// var filters DBUserFilters
// filters.Id.SetFilter(5) // => filters users having the Id '5'
// filters.Name.SetFilter("to") // => filters users having the name "to"
// ```
type DBUserFilters struct {
	Id   DBIntFilter    // by User Id
	Name DBStringFilter // by user's name
}

// Filters that can be used to filter Categories when doing operations on the
// CategoryDatabase
//
// example:
// ```go
// var filters DBCategoryFilters
// filters.Ids.SetFilter([]{5, 3})
// ```
type DBCategoryFilters struct {
	Ids       DBIntArrayFilter    // by Categories Ids
	Names     DBStringArrayFilter // by Categories names
	UserId    DBIntFilter         // by User Id
	ParentIds DBIntArrayFilter    // by Parent Categories Ids (TODO Remove)
}

// Filters that can be used to filter Bank Accounts when doing operations on the
// BankAccountDataBase.
//
// example:
// ```go
// var filters DBAccountFilters
// filters.Ids.SetFilter([]{5, 3})
// ```
type DBAccountFilters struct {
	Ids     DBIntArrayFilter    // by Bank Account Ids
	UserId  DBIntFilter         // by User Id (TODO Remove?)
	BankIds DBIntArrayFilter    // by Bank Ids corresponding to the accounts
	Names   DBStringArrayFilter // by Bank Account names
}

// Filters that can be used to filter Banks when doing operations on the
// BankDataBase
//
// example:
// ```go
// var filters DBBankFilters
// filters.Ids.SetFilter([]{5, 3})
// ```
type DBBankFilters struct {
	Ids    DBIntArrayFilter    // by Bank Ids
	UserId DBIntFilter         // by User Id
	Names  DBStringArrayFilter // by Bank names
}

// Filters that can be used to filter Transactions when doing operations on the
// TransactionDataBase.
//
// example:
// ```go
// var filters DBTransactionFilters
// filters.Ids.SetFilter([]{5, 3})
// ```
type DBTransactionFilters struct {
	Ids                 DBIntArrayFilter    // by Transactions Ids
	UserId              DBIntFilter         // by User Id (TODO Remove)
	BankIds             DBIntArrayFilter    // by Bank Ids (TODO Remove)
	AccountIds          DBIntArrayFilter    // by Bank Accounts Ids
	CategoryIds         DBIntArrayFilter    // by Categories Ids
	FromTransactionDate DBTimeFilter        // by minimum transaction date
	ToTransactionDate   DBTimeFilter        // by maximum transaction date
	FromRecordDate      DBTimeFilter        // by minimum record date
	ToRecordDate        DBTimeFilter        // by maximum record date
	MinDebit            DBFloatFilter       // by minimum debit
	MaxDebit            DBFloatFilter       // by maximum debit
	MinCredit           DBFloatFilter       // by minimum credit
	MaxCredit           DBFloatFilter       // by maximum credit
	References          DBStringArrayFilter // by bank's reference
}

type dbBaseParam struct{ activated bool }

type DBIntParam struct {
	dbBaseParam
	value int
}

type DBFloatParam struct {
	dbBaseParam
	value float32
}

type DBStringParam struct {
	dbBaseParam
	value string
}

type DBTimeParam struct {
	dbBaseFilter
	value string
}

func (d *DBIntParam) SetParam(val int) {
	d.activated = true
	d.value = val
}

func (d *DBFloatParam) SetParam(val float32) {
	d.activated = true
	d.value = val
}

func (d *DBStringParam) SetParam(val string) {
	d.activated = true
	d.value = val
}

func (d *DBTimeParam) SetParam(val time.Time) {
	d.activated = true
	d.value = val
}

// Common base of filters
type dbBaseFilter struct{ activated bool }

func (d dbBaseFilter) isFilterActivated() bool { return d.activated }

// ugly generic dbFilter interface for using them in generic helpers
type dbParamInterface interface {
	isParamActivated() bool
	getParamValue() interface{}
}

func (d dbBaseParam) isParamActivated() bool  { return d.activated }
func (d DBIntParam) getParamValue() int       { return d.value }
func (d DBStringParam) getParamValue() string { return d.value }
func (d DBFloatParam) getParamValue() float32 { return d.value }
func (d DBTimeParam) getParamValue() string   { return d.value }

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
