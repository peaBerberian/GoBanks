package database

import "database/sql"
import "sync"
import _ "github.com/go-sql-driver/mysql"

// must respect the goBanksDatabase interface
type goBanksSql struct {
	db    *sql.DB
	mutex sync.Mutex
}

func newMySqlDB(user string, pw string, access string,
	database string) (*goBanksSql, databaseError) {

	var gbs *goBanksSql

	var db, err = sql.Open("mysql", user+":"+pw+"@"+access+"/"+database+
		"?parseTime=true")
	if err != nil {
		return gbs, genericDatabaseError{
			err:  err.Error(),
			code: DatabaseConnectionErrorCode,
		}
	}

	gbs = new(goBanksSql)
	gbs.db = db
	return gbs, nil
}

func (gbs *goBanksSql) Close() (err error) {
	return gbs.db.Close()
}

// setMysqlDB connect to the mysql database from the given config.
func setMysqlDB(c map[string]interface{}) (GoBanksDataBase, databaseError) {

	var parseField = func(field string) (string, error) {
		if val, ok := c[field]; ok {
			if val, ok := val.(string); ok {
				return val, nil
			}
			var err = genericDatabaseError{err: "The value for the field \"" +
				field + "\" in the given mysql configuration is not in the" +
				" right format.",
				code: DatabaseConfigurationErrorCode}
			return "", err
		}
		var err = genericDatabaseError{err: "No field \"" + field +
			"\" in the given mysql configuration.",
			code: DatabaseConfigurationErrorCode}
		return "", err
	}

	var args = make([]string, 4)
	for i, field := range []string{"user", "password",
		"access", "database"} {

		val, err := parseField(field)
		if err != nil {
			return nil, genericDatabaseError{err: err.Error()}
		}
		args[i] = val
	}

	var gdb, err = newMySqlDB(args[0], args[1], args[2], args[3])
	return gdb, err
}
