package mysql

import "database/sql"
import "errors"

import dbt "github.com/peaberberian/GoBanks/database/types"

func (gbs *GoBanksSql) AddTransaction(transac dbt.Transaction) (id int,
	err error) {
	if transac.LinkedAccountDbId == 0 {
		return 0, errors.New("The linked account must be added to the" +
			" database before the transaction.")
	}
	var res sql.Result

	gbs.mutex.Lock()
	res, err = gbs.db.Exec("insert into transaction "+
		"(account_id, label, debit, credit, date_of_transaction,"+
		" date_of_record) "+
		"values (?, ?, ?, ?, ?, ?)",
		transac.LinkedAccountDbId,
		transac.Label,
		transac.Debit,
		transac.Credit,
		transac.TransactionDate,
		transac.RecordDate)
	gbs.mutex.Unlock()

	if err != nil {
		return 0, err
	}
	var id64 int64
	id64, err = res.LastInsertId()
	id = int(id64)
	transac.DbId = id
	return
}

func (gbs *GoBanksSql) RemoveTransaction(transacId int) (err error) {
	gbs.mutex.Lock()
	_, err = gbs.db.Exec("DELETE FROM transaction "+
		"WHERE id=?", transacId)
	gbs.mutex.Unlock()

	return
}

func (gbs *GoBanksSql) UpdateTransaction(transac dbt.Transaction,
) (err error) {
	gbs.mutex.Lock()
	_, err = gbs.db.Exec("UPDATE transaction "+
		"SET account_id=?, label=?, debit=?, credit=?,"+
		" date_of_transaction=?, date_of_record=? "+
		"WHERE id=?",
		transac.LinkedAccountDbId,
		transac.Label,
		transac.Debit,
		transac.Credit,
		transac.TransactionDate,
		transac.RecordDate,
		transac.DbId,
	)
	gbs.mutex.Unlock()

	if err != nil {
		return
	}
	return nil
}

func (gbs *GoBanksSql) GetTransaction(transacId int) (t dbt.Transaction,
	err error) {
	gbs.mutex.Lock()
	row := gbs.db.QueryRow("select account_id, label, debit, credit,"+
		" date_of_transaction, date_of_record from transaction where id=?",
		transacId)
	gbs.mutex.Unlock()

	t.DbId = transacId
	err = row.Scan(&t.LinkedAccountDbId, &t.Label, &t.Debit, &t.Credit,
		&t.TransactionDate, &t.RecordDate)

	if err != nil {
		return dbt.Transaction{}, err
	}
	return
}

// TODO
func (gbs *GoBanksSql) GetTransactions(filters dbt.TransactionFilters,
) (ts []dbt.Transaction, err error) {
	gbs.mutex.Lock()
	var rows = new(sql.Rows)
	rows, err = gbs.db.Query("select id, account_id, label, debit, credit,"+
		" date_of_transaction, date_of_record from transaction"+
		" WHERE account_id=?",
		filters.Values.Accounts[0])
	gbs.mutex.Unlock()

	if err != nil {
		return []dbt.Transaction{}, err
	}

	for rows.Next() {
		var t dbt.Transaction
		err = rows.Scan(&t.DbId, &t.LinkedAccountDbId, &t.Label, &t.Debit,
			&t.Credit, &t.TransactionDate, &t.RecordDate)
		if err != nil {
			return
		}
		ts = append(ts, t)
	}

	return
}
