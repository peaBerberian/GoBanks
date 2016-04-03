package database

import "database/sql"
import "errors"

const TRANSACTION_TABLE = "transaction"

var TRANSACTION_FIELDS = []string{"account_id", "label", "debit", "credit",
	"date_of_transaction", "date_of_record"}

func (gbs *goBanksSql) AddTransaction(trns Transaction) (id int,
	err error) {
	if trns.LinkedAccountDbId == 0 {
		return 0, errors.New("The linked account must be added to the" +
			" database before the transaction.")
	}

	values := make([]interface{}, 0)
	values = append(values, trns.LinkedAccountDbId, trns.Label,
		trns.Debit, trns.Credit, trns.TransactionDate, trns.RecordDate)

	id, err = gbs.insertInTable(TRANSACTION_TABLE, TRANSACTION_FIELDS, values)
	if err != nil {
		return 0, err
	}
	trns.DbId = id
	return
}

func (gbs *goBanksSql) RemoveTransaction(id int) (err error) {
	return gbs.removeIdFromTable(TRANSACTION_TABLE, id)
}

func (gbs *goBanksSql) UpdateTransaction(trns Transaction,
) (err error) {
	values := make([]interface{}, 0)
	values = append(values, trns.LinkedAccountDbId, trns.Label,
		trns.Debit, trns.Credit, trns.TransactionDate, trns.RecordDate)

	return gbs.updateInTableFromId(TRANSACTION_TABLE, trns.DbId,
		TRANSACTION_FIELDS, values)
}

func (gbs *goBanksSql) GetTransaction(id int) (t Transaction,
	err error) {
	row := gbs.getFromTable(TRANSACTION_TABLE, id, TRANSACTION_FIELDS)
	t.DbId = id
	err = row.Scan(&t.LinkedAccountDbId, &t.Label, &t.Debit, &t.Credit,
		&t.TransactionDate, &t.RecordDate)
	return
}

// TODO search from string (last 3 filters)
func (gbs *goBanksSql) GetTransactions(f TransactionFilters,
) (ts []Transaction, err error) {
	var queryString string
	var selectString = "select id, account_id, label, debit, credit," +
		" date_of_transaction, date_of_record from " + TRANSACTION_TABLE + " "
	var whereString = "WHERE "
	var atLeastOneFilter = false
	var sqlArguments = make([]interface{}, 0)

	if f.Filters.Types {
		if len(f.Values.Types) > 0 {
			str, arg := addSqlFilterStringArray("name", f.Values.Types...)
			whereString += str + " "
			sqlArguments = append(sqlArguments, arg...)
			atLeastOneFilter = true
		} else {
			return ts, nil
		}
	}
	if f.Filters.Accounts {
		if len(f.Values.Accounts) > 0 {
			if atLeastOneFilter {
				whereString += "AND "
			}
			addSqlFilterIntArray("account", f.Values.Accounts...)
			atLeastOneFilter = true
		} else {
			return ts, nil
		}
	}

	// Now that's what I call ugly! vol. 74
	// Ugly but nice though...

	if f.Filters.FromTransactionDate {
		if atLeastOneFilter {
			whereString += "AND "
		}
		whereString += "date_of_transaction >= ? "
		sqlArguments = append(sqlArguments, f.Values.FromTransactionDate)
		atLeastOneFilter = true
	}

	if f.Filters.ToTransactionDate {
		if atLeastOneFilter {
			whereString += "AND "
		}
		whereString += "date_of_transaction <= ? "
		sqlArguments = append(sqlArguments, f.Values.ToTransactionDate)
		atLeastOneFilter = true
	}

	if f.Filters.FromRecordDate {
		if atLeastOneFilter {
			whereString += "AND "
		}
		whereString += "date_of_record >= ? "
		sqlArguments = append(sqlArguments, f.Values.FromRecordDate)
		atLeastOneFilter = true
	}

	if f.Filters.ToRecordDate {
		if atLeastOneFilter {
			whereString += "AND "
		}
		whereString += "date_of_record < ? "
		sqlArguments = append(sqlArguments, f.Values.ToRecordDate)
		atLeastOneFilter = true
	}

	if f.Filters.MinDebit {
		if atLeastOneFilter {
			whereString += "AND "
		}
		whereString += "debit >= ? "
		sqlArguments = append(sqlArguments, f.Values.MinDebit)
		atLeastOneFilter = true
	}

	if f.Filters.MaxDebit {
		if atLeastOneFilter {
			whereString += "AND "
		}
		whereString += "debit <= ? "
		sqlArguments = append(sqlArguments, f.Values.MaxDebit)
		atLeastOneFilter = true
	}

	if f.Filters.MinCredit {
		if atLeastOneFilter {
			whereString += "AND "
		}
		whereString += "credit >= ? "
		sqlArguments = append(sqlArguments, f.Values.MinCredit)
		atLeastOneFilter = true
	}

	if f.Filters.MaxCredit {
		if atLeastOneFilter {
			whereString += "AND "
		}
		whereString += "credit <= ? "
		sqlArguments = append(sqlArguments, f.Values.MaxCredit)
		atLeastOneFilter = true
	}

	if atLeastOneFilter {
		queryString = selectString + whereString
	} else {
		queryString = selectString
	}

	gbs.mutex.Lock()
	var rows = new(sql.Rows)
	rows, err = gbs.db.Query(queryString, sqlArguments...)
	gbs.mutex.Unlock()

	if err != nil {
		return ts, err
	}

	for rows.Next() {
		var trn Transaction
		err = rows.Scan(&trn.DbId, &trn.LinkedAccountDbId, &trn.Label,
			&trn.Debit, &trn.Credit, &trn.TransactionDate, &trn.RecordDate)
		if err != nil {
			return
		}
		ts = append(ts, trn)
	}

	return
}
