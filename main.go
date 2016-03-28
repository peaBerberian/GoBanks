package main

import "github.com/peaberberian/GoBanks/config"
import "github.com/peaberberian/GoBanks/database"

// used for tests
import "fmt"
import def "github.com/peaberberian/GoBanks/database/definitions"

// import "github.com/peaberberian/GoBanks/login"

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	// Create global database object
	db, err := database.New(conf.Databases)
	if err != nil {
		panic(err)
	}

	var f def.TransactionFilters
	f.Filters.MinCredit = true
	f.Values.MinCredit = 1000
	// f.Filters.MinDebit = true
	// f.Values.MinDebit = -12

	debs, err := db.GetTransactions(f)
	// random tests
	// usr, err := login.LoginUser(db, "abraham", "Simpson")
	fmt.Println(debs, err)
	defer db.Close()
}
