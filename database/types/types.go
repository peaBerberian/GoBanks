package types

import "time"

// Interface that any implementation of GoBanks databases must respect
type GoBanksDataBase interface {
	// Free/close the db if needed
	Close() error

	// Number of users stored in the database
	UserLength() (len int, err error)

	// Add a single user
	// Returns user's DbId / possible error
	AddUser(User) (int, error)

	// Remove a single user from its DbId
	RemoveUser(int) error

	// Update the attributes of a single user, from its DbId
	UpdateUser(User) error

	// Get a single user, from its DbId
	GetUser(int) (User, error)

	// Get multiple users, based on filters
	GetUsers(UserFilters) ([]User, error)

	// Add a single
	// Returns bank's DbId / possible error
	AddBank(Bank) (int, error)

	// Remove a single bank from its DbId
	RemoveBank(int) error

	// Update the attributes of a single bank, from its DbId
	UpdateBank(Bank) error

	// Get a single bank, from its DbId
	GetBank(int) (Bank, error)

	// Get multiple banks, based on filters
	GetBanks(BankFilters) ([]Bank, error)

	// Add a single account for a specific bank DbId
	// Returns transaction's DbId / possible error
	AddBankAccount(BankAccount) (int, error)

	// Remove a single account from its DbId
	RemoveBankAccount(int) error

	// Update the attributes of a single account, from its DbId
	UpdateBankAccount(BankAccount) error

	// Get a single account, from its DbId
	GetBankAccount(int) (BankAccount, error)

	// Get multiple account, based on filters
	GetBankAccounts(BankAccountFilters) ([]BankAccount, error)

	// Add a single transaction
	// Returns transaction's DbId / possible error
	AddTransaction(Transaction) (int, error)

	// Remove a single transaction from its DbId
	RemoveTransaction(int) error

	// Update the attributes of a single transaction, from its DbId
	UpdateTransaction(Transaction) error

	// Get a single transaction, from its DbId
	GetTransaction(int) (Transaction, error)

	// Get multiple transactions, based on filters
	GetTransactions(TransactionFilters) ([]Transaction, error)
}

type BankAccountFilters struct {
	Filters struct {
		Banks bool
		Names bool
	}
	Values struct {
		Banks []int
		Names []string
	}
}

type BankFilters struct {
	Filters struct {
		Users bool
		Names bool
	}
	Values struct {
		Users []int
		Names []string
	}
}

type UserFilters struct {
	Filters struct {
		Names     bool
		Tokens    bool
		Permanent bool
	}
	Values struct {
		Names     []string
		Tokens    []string
		Permanent bool
	}
}

type TransactionFilters struct {
	Filters struct {
		Accounts            bool
		Types               bool
		FromTransactionDate bool
		ToTransactionDate   bool
		FromRecordDate      bool
		ToRecordDate        bool
		MinDebit            bool
		MaxDebit            bool
		MinCredit           bool
		MaxCredit           bool
		SearchLabel         bool
		SearchDescription   bool
		SearchReference     bool
	}
	Values struct {
		Accounts            []int
		Types               []string
		FromTransactionDate time.Time
		ToTransactionDate   time.Time
		FromRecordDate      time.Time
		ToRecordDate        time.Time
		MinDebit            float32
		MaxDebit            float32
		MinCredit           float32
		MaxCredit           float32
		SearchLabel         string
		SearchDescription   string
		SearchReference     string
	}
}

type Transaction struct {
	// Id of the transaction in the database
	// Set by the database's methods.
	// (starts at 1, 0 if not added to the database)
	DbId int

	// DbId for the account concerned by this transaction
	LinkedAccountDbId int

	// Label describing the transaction
	Label string

	// Category of the transaction
	Type string

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

type User struct {
	// Id for the user in the database
	// Set by the database's methods.
	// (starts at 1, 0 if not added to the database)
	DbId int

	// User's Name
	Name string

	// Hash of the user's password
	PasswordHash string

	// Password's salt
	Salt string

	// User's current token
	Token string

	// True if the user is a permanent one
	Permanent bool
}

type Bank struct {
	// Id of the bank in the database
	// Set by the database's methods.
	// (starts at 1, 0 if not added to the database)
	DbId int

	// User linked to this Bank
	// LinkedUser User
	LinkedUserDbId int

	// Name of the bank
	Name string

	// Optional description
	Description string
}

type BankAccount struct {
	// Id of the bank account in the database
	// Set by the database's methods.
	// (starts at 1, 0 if not added to the database)
	DbId int

	// TODO
	// Bank DbId linked to this account
	LinkedBankDbId int

	// Name of the bank account
	Name string

	// Amount of money used as a base. (Not all transactions may be
	// available, we have to set the base to circumvent this)
	BaseAmount float32

	// Optional description
	Description string
}
