package mysql

import "database/sql"
import "sync"
import _ "github.com/go-sql-driver/mysql"

// must respect the GoBanksDatabase interface
type GoBanksSql struct {
	db *sql.DB

	// TODO
	mutex sync.Mutex
}

func New(user string, pw string, access string, database string) (gbs *GoBanksSql, err error) {
	var db *sql.DB
	db, err = sql.Open("mysql", user+":"+pw+"@"+access+"/"+database+"?parseTime=true")
	if err != nil {
		return
	}
	gbs = new(GoBanksSql)
	gbs.db = db
	return gbs, nil
}

func (gbs *GoBanksSql) Close() (err error) {
	return gbs.db.Close()
}
