package mysql

import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import dbt "github.com/peaberberian/GoBanks/database/types"

// must respect the GoBankDatabase interface
type GoBankSql struct {
	db *sql.DB
}

func New(user string, pw string, access string, database string) (gbs *GoBankSql, err error) {
	var db *sql.DB
	db, err = sql.Open("mysql", user+":"+pw+"@"+access+"/"+database)
	if err != nil {
		return
	}
	gbs = new(GoBankSql)
	gbs.db = db
	return gbs, nil
}

func (gbs *GoBankSql) Close() (err error) {
	return gbs.db.Close()
}

func (gbs *GoBankSql) AddTransaction(transac dbt.Transaction, accountId int) (err error) {
	_, err = gbs.db.Exec("insert into transaction "+
		"(account_id, label, debit, credit, date_of_transaction, date_of_record) "+
		"values (?, ?, ?, ?, ?, ?)",
		accountId,
		transac.Label,
		transac.Debit,
		transac.Credit,
		transac.TransactionDate,
		transac.RecordDate)
	return
}

func (gbs *GoBankSql) RemoveTransaction(transacId int) (err error) {
	_, err = gbs.db.Exec("DELETE FROM transaction "+
		"WHERE id=?", transacId)
	return
}

// TODO
func (gbs *GoBankSql) UpdateTransaction(transacId int, transac dbt.Transaction) (err error) {
	return nil
}

// TODO
func (gbs *GoBankSql) GetTransaction(transacId int) (t dbt.Transaction, err error) {
	return t, nil
}

// TODO
func (gbs *GoBankSql) GetAllTransactions() (ts []dbt.Transaction, err error) {
	return ts, nil
}
