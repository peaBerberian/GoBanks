package database

import "database/sql"
import "sync"
import _ "github.com/go-sql-driver/mysql"

// must respect the goBanksDatabase interface
type goBanksSql struct {
	db    *sql.DB
	mutex sync.Mutex
}

func newMySqlDB(user string, pw string, access string, database string,
) (gbs *goBanksSql, err error) {
	var db *sql.DB
	db, err = sql.Open("mysql", user+":"+pw+"@"+access+"/"+database+
		"?parseTime=true")
	if err != nil {
		return
	}
	gbs = new(goBanksSql)
	gbs.db = db
	return gbs, nil
}

func (gbs *goBanksSql) Close() (err error) {
	return gbs.db.Close()
}
