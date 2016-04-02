package main

import "github.com/peaberberian/GoBanks/config"
import "github.com/peaberberian/GoBanks/database"

// used for tests
import "fmt"
import "github.com/peaberberian/GoBanks/login"

// TODO Read config here and setup db choice etc. from here
func init() {
}

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

	token, err := login.CreateToken("oscarito", db)
	fmt.Println("token: ", token, err)
	myToken, err := login.ParseToken(token)
	fmt.Println("token: ", myToken, err)

	defer db.Close()
}
