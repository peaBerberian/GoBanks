package types

import "time"

// Interface that any implementation of GoBanks databases must respect
type GoBanksDataBase interface {
	// Free/close the db if needed
	Close() error

	// Add a single bank
	// Returns bank's id / possible error
	AddBank(Bank) (int, error)

	// Remove a single bank from its id
	RemoveBank(int) error

	// Update the attributes of a single bank, from its id
	UpdateBank(int, Bank) error

	// Get a single bank, from its id
	GetBank(int) (Bank, error)

	// Returns every single banks (/!\ may be too much, used
	// now, but will be removed with limits, filters...)
	GetAllBanks() ([]Bank, error)

	// Add a single account for a specific bank id
	// Returns transaction's id / possible error
	AddBankAccount(BankAccount) (int, error)

	// Remove a single account from its id
	RemoveBankAccount(int) error

	// Update the attributes of a single account, from its id
	UpdateBankAccount(int, BankAccount) error

	// Get a single account, from its id
	GetBankAccount(int) (BankAccount, error)

	// Returns every single accounts (/!\ may be too much, used
	// now, but will be removed with limits, filters...)
	GetAllBankAccounts() ([]BankAccount, error)

	// Add a single transaction
	// Returns transaction's id / possible error
	AddTransaction(Transaction) (int, error)

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
	AccountId       int
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
	Name        string
	Description string
}

type BankAccount struct {
	BankId      int
	Name        string
	BaseAmount  float32
	Description string
}
