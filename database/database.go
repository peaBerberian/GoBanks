package database

// Globally accessible configuration
var dbConfig struct {
	name   string
	config map[string]interface{}
}

// New creates a new database with the configuration given in the Set
// function.
func New() (gdb GoBanksDataBase, err error) {
	if dbConfig.name == "" {
		return gdb, databaseError{err: "Database is not yet configured. " +
			"Please call the Set method.", code: DatabaseConfigurationError}
	}

	switch dbConfig.name {
	case "mySql":
		return newMysql()
	default:
		var err = databaseError{err: "Can't manage a(n) " + dbConfig.name +
			" database.", code: DatabaseConfigurationError}
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
	dbConfig.name = typ

	if val, ok := config.(map[string]interface{}); ok {
		dbConfig.config = val
	} else {
		var err = databaseError{err: "Wrong database configuration.",
			code: DatabaseConfigurationError}
		return err
	}
	return nil
}

// newMysql creates specifically mysql databases.
func newMysql() (gdb GoBanksDataBase, err error) {
	var c = dbConfig.config

	var parseField = func(field string) (string, error) {
		if val, ok := c[field]; ok {
			if val, ok := val.(string); ok {
				return val, nil
			}
			var err = databaseError{err: "The value for the field \"" + field +
				"\" in the given mysql configuration is not in the right format.",
				code: DatabaseConfigurationError}
			return "", err
		}
		var err = databaseError{err: "No field \"" + field +
			"\" in the given mysql configuration.",
			code: DatabaseConfigurationError}
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

	gdb, err = newMySqlDB(args[0], args[1], args[2], args[3])
	return
}
