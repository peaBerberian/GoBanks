package mysql

import "database/sql"

import dbt "github.com/peaberberian/GoBanks/database/types"

func (gbs *GoBanksSql) AddTransaction(transac dbt.Transaction) (id int, err error) {
	var res sql.Result
	res, err = gbs.db.Exec("insert into transaction "+
		"(account_id, label, debit, credit, date_of_transaction, date_of_record) "+
		"values (?, ?, ?, ?, ?, ?)",
		transac.AccountId,
		transac.Label,
		transac.Debit,
		transac.Credit,
		transac.TransactionDate,
		transac.RecordDate)
	if err != nil {
		return
	}
	var id64 int64
	id64, err = res.LastInsertId()
	id = int(id64)
	return
}

func (gbs *GoBanksSql) RemoveTransaction(transacId int) (err error) {
	_, err = gbs.db.Exec("DELETE FROM transaction "+
		"WHERE id=?", transacId)
	return
}

// TODO
func (gbs *GoBanksSql) UpdateTransaction(transacId int, transac dbt.Transaction) (err error) {
	return nil
}

// TODO
func (gbs *GoBanksSql) GetTransaction(transacId int) (t dbt.Transaction, err error) {
	return t, nil
}

// TODO
func (gbs *GoBanksSql) GetAllTransactions() (ts []dbt.Transaction, err error) {
	return ts, nil
}
