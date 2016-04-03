package main

import "github.com/peaberberian/GoBanks/auth"
import "github.com/peaberberian/GoBanks/database"
import "github.com/peaberberian/GoBanks/api"

func main() {
	// Read the config file
	conf, err := getConfig()
	if err != nil {
		panic(err)
	}

	// Initialize Database (database.GoDB)
	if err := database.Connect(conf.Database); err != nil {
		panic(err)
	}

	// Update token expiration from config
	auth.SetTokenExpiration(conf.TokenExpiration)

	api.Start(8080)
	database.GoDB.Close()
}
