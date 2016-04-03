package main

import "github.com/peaberberian/GoBanks/database"
import "github.com/peaberberian/GoBanks/api"

// used for tests
import "fmt"

type GoBanksError interface {
	error
	ErrorCode() uint32
}

// set configuration before starting
func init() {
	conf, err := getConfig()
	if err != nil {
		panic(err)
	}

	var typ = conf.Databases.Typ
	if val, ok := conf.Databases.Config.(map[string]interface{}); ok {
		err = database.Set(typ, val[typ])
		if err != nil {
			panic(err)
		}
	} else {
		panic("The configuration file is not valid " +
			"(The database configuration could not be understood).")
	}

	api.SetTokenExpiration(conf.TokenExpiration)
}

// create database and launch the application
func main() {
	db, err := database.New()
	if err != nil {
		panic(err)
	}

	token, err := api.CreateToken("oscarito", db)
	fmt.Println("token: ", token, err)
	myToken, err := api.ParseToken(token)
	fmt.Println("token: ", myToken, err)

	defer db.Close()
}
