package database

import "errors"

import ms "github.com/peaberberian/GoBanks/database/mysql"
import "github.com/peaberberian/GoBanks/config"
import "github.com/peaberberian/GoBanks/database/types"

func New(dbConfig config.DatabasesConfig) (gdb types.GoBanksDataBase, err error) {
	switch dbConfig.DatabaseType {
	case "mySql":
		var m = dbConfig.MySql
		gdb, err = ms.New(m.User, m.Password, m.Access, m.Database)
		if err != nil {
			return
		}
	default:
		err = errors.New("Can't manage a(n) " + dbConfig.DatabaseType + " database.")
		return
	}
	return gdb, nil
}
