package database

var GoDB GoBanksDataBase

// Connect connect to the chosen database.
// This is needed to then be able to use the exported GoDB afterwards.
//
// It takes in argument the awaited config for it.
// The configuration values depends on which database type was chosen.
// If any problem is detected while setting this config, an error will
// be returned.
func Connect(config interface{}) databaseError {
	var dbConfig map[string]interface{}

	if val, ok := config.(map[string]interface{}); ok {
		dbConfig = val
	} else {
		return databaseConfigurationError{}
	}

	var err databaseError
	GoDB, err = setMysqlDB(dbConfig)
	return err
}
