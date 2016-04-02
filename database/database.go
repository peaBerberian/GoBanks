package database

import mysql "github.com/peaberberian/GoBanks/database/mysql"
import def "github.com/peaberberian/GoBanks/database/definitions"

type databaseConfiguration struct {
	name   string
	config map[string]interface{}
}

var dbConf databaseConfiguration

// New creates a new database with the configuration given in the Set
// function.
func New() (gdb def.GoBanksDataBase, err error) {
	switch dbConf.name {
	case "mySql":
		return newMysql()
	default:
		var err = DatabaseError{err: "Can't manage a(n) " + dbConf.name +
			" database.", ErrorCode: DatabaseConfigurationError}
		return gdb, err
	}
	return
}

// Set initialize configuration values. Please call this before calling
// the New function.
// It takes in argument the database type (mysql, mongo...) and the
// awaited config for it (in the form of a map[string]interface{}).
// The configuration values depends on which database type was chosen.
// If any problem is detected while setting this config, an error will
// be returned.
func Set(typ string, config interface{}) error {
	dbConf.name = typ

	if val, ok := config.(map[string]interface{}); ok {
		dbConf.config = val
	} else {
		var err = DatabaseError{err: "Wrong database configuration.",
			ErrorCode: DatabaseConfigurationError}
		return err
	}
	return nil
}

// newMysql creates specifically mysql databases.
func newMysql() (gdb def.GoBanksDataBase, err error) {
	var c = dbConf.config

	var parseField = func(field string) (string, error) {
		if val, ok := c[field]; ok {
			if val, ok := val.(string); ok {
				return val, nil
			}
			var err = DatabaseError{err: "The value for the field \"" + field +
				"\" in the given mysql configuration is not in the right format.",
				ErrorCode: DatabaseConfigurationError}
			return "", err
		}
		var err = DatabaseError{err: "No field \"" + field +
			"\" in the given mysql configuration.",
			ErrorCode: DatabaseConfigurationError}
		return "", err
	}

	var args = make([]string, 4)
	for i, field := range []string{"user", "password",
		"access", "database"} {

		val, err := parseField(field)
		if err != nil {
			return gdb, err
		}
		args[i] = val
	}

	gdb, err = mysql.New(args[0], args[1], args[2], args[3])
	return
}
