package types

import "time"

// Interface that any implementation of GoBanks databases must respect
type GoBanksDataBase interface {
	// Free/close the db if needed
	Close() error

	// // Add a single bank for a specific account id
	// AddBank(Bank) error

	// // Remove a single bank from its id
	// RemoveBank(int) error

	// // Update the attributes of a single bank, from its id
	// UpdateBank(int, Bank) error

	// // Get a single bank, from its id
	// GetBank(int) (Bank, error)

	// // Returns every single banks (/!\ may be too much, used
	// // now, but will be removed with limits, filters...)
	// GetAllBanks() ([]Bank, error)

	// // Add a single account for a specific bank id
	// AddAccount(Account, int) error

	// // Remove a single account from its id
	// RemoveAccount(int) error

	// // Update the attributes of a single account, from its id
	// UpdateAccount(int, Account) error

	// // Get a single account, from its id
	// GetAccount(int) (Account, error)

	// // Returns every single accounts (/!\ may be too much, used
	// // now, but will be removed with limits, filters...)
	// GetAllAccounts() ([]Account, error)

	// Add a single transaction for a specific account id
	AddTransaction(Transaction, int) error

	// Remove a single transaction from its id
	RemoveTransaction(int) error

	// Update the attributes of a single transaction, from its id
	UpdateTransaction(int, Transaction) error

	// Get a single transaction, from its id
	GetTransaction(int) (Transaction, error)

	// Returns every single transactions (/!\ may be too much, used
	// now, but will be removed with limits, filters...)
	GetAllTransactions() ([]Transaction, error)
}

type Transaction struct {
	// AccountId       int
	Label           string
	Type            string
	Description     string
	TransactionDate time.Time
	RecordDate      time.Time
	Debit           float32
	Credit          float32
	Reference       string
}

type Bank struct {
	name        string
	description string
}

type Account struct {
	// BankId      int
	name        string
	baseAmount  float32
	description string
}
