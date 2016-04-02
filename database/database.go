package database

import "errors"

import config "github.com/peaberberian/GoBanks/config"
import mysql "github.com/peaberberian/GoBanks/database/mysql"
import def "github.com/peaberberian/GoBanks/database/definitions"

func New(dbConfig config.DatabasesConfig) (gdb def.GoBanksDataBase,
	err error) {
	switch dbConfig.DatabaseType {
	case "mySql":
		var m = dbConfig.MySql
		gdb, err = mysql.New(m.User, m.Password, m.Access, m.Database)
	default:
		err = errors.New("Can't manage a(n) " + dbConfig.DatabaseType +
			" database.")
		return
	}
	return
}
