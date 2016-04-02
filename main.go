package main

import "github.com/peaberberian/GoBanks/config"
import "github.com/peaberberian/GoBanks/database"
import "github.com/peaberberian/GoBanks/login"

// used for tests
import "fmt"

// set configuration before starting
func init() {
	conf, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	var typ = conf.Databases.Type
	if val, ok := conf.Databases.Config.(map[string]interface{}); ok {
		err = database.Set(typ, val[typ])
		if err != nil {
			panic(err)
		}
	} else {
		panic("The configuration file is not valid" +
			"(The database configuration could not be understood).")
	}

	login.SetTokenExpiration(conf.TokenExpiration)
}

// create database and launch the application
func main() {
	db, err := database.New()
	if err != nil {
		panic(err)
	}

	token, err := login.CreateToken("oscarito", db)
	fmt.Println("token: ", token, err)
	myToken, err := login.ParseToken(token)
	fmt.Println("token: ", myToken, err)

	defer db.Close()
}
