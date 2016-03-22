package main

import "os"

import "github.com/peaberberian/GoBanks/config"
import "github.com/peaberberian/GoBanks/database"
import "github.com/peaberberian/GoBanks/database/types"

// just for tests
import "github.com/peaberberian/GoBanks/file/qif"

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	db, err := database.New(conf.Databases)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// // tests
	// err = AddQifFile("./toto.qif", "DD/MM/YY", db, 1)
	// if err != nil {
	// 	panic(err)
	// }
}

// just for tests
func AddQifFile(filePath string, dateFormat string, db types.GoBanksDataBase, accountId int) (err error) {
	var f *os.File
	f, err = os.Open(filePath)
	if err != nil {
		return
	}
	ts, err := qif.ParseFile(f, accountId, dateFormat)
	if err != nil {
		return
	}
	for _, t := range ts {
		_, err = db.AddTransaction(t)
		if err != nil {
			return
		}
	}
	return nil
}

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
