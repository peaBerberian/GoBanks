package mysql

import "database/sql"
import "errors"

import def "github.com/peaberberian/GoBanks/database/definitions"

const ACCOUNT_TABLE = "account"

var ACCOUNT_FIELDS = []string{"bank_id", "name", "base_amount", "description"}

func (gbs *goBanksSql) AddBankAccount(acc def.BankAccount) (id int,
	err error) {
	if acc.LinkedBankDbId == 0 {
		return 0, errors.New("The linked bank must be added to the" +
			" database before the bank account.")
	}

	values := make([]interface{}, 0)
	values = append(values, acc.LinkedBankDbId, acc.Name,
		acc.BaseAmount, acc.Description)

	id, err = gbs.insertInTable(ACCOUNT_TABLE, ACCOUNT_FIELDS, values)
	if err != nil {
		return 0, err
	}
	acc.DbId = id
	return
}

func (gbs *goBanksSql) RemoveBankAccount(id int) (err error) {
	return gbs.removeIdFromTable(ACCOUNT_TABLE, id)
}

func (gbs *goBanksSql) UpdateBankAccount(acc def.BankAccount) (err error) {
	values := make([]interface{}, 0)
	values = append(values, acc.LinkedBankDbId, acc.Name,
		acc.BaseAmount, acc.Description)

	return gbs.updateInTableFromId(ACCOUNT_TABLE, acc.DbId,
		TRANSACTION_FIELDS, values)
}

func (gbs *goBanksSql) GetBankAccount(id int) (t def.BankAccount,
	err error) {
	row := gbs.getFromTable(TRANSACTION_TABLE, id, TRANSACTION_FIELDS)
	t.DbId = id
	err = row.Scan(&t.LinkedBankDbId, &t.Name, &t.BaseAmount, &t.Description)
	return
}

func (gbs *goBanksSql) GetBankAccounts(f def.BankAccountFilters,
) (accounts []def.BankAccount, err error) {
	var queryString string
	var selectString = "select id, bank_id, name, base_amount," +
		" description from account "
	var whereString = "WHERE "
	var atLeastOneFilter = false
	var sqlArguments = make([]interface{}, 0)

	if f.Filters.Names {
		if len(f.Values.Names) > 0 {
			str, arg := addSqlFilterStringArray("name", f.Values.Names...)
			whereString += str + " "
			sqlArguments = append(sqlArguments, arg...)
			atLeastOneFilter = true
		} else {
			return
		}
	}
	if f.Filters.Banks {
		if len(f.Values.Banks) > 0 {
			if atLeastOneFilter {
				whereString += "AND "
			}
			addSqlFilterIntArray("token", f.Values.Banks...)
			atLeastOneFilter = true
		} else {
			return
		}
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
		return
	}

	for rows.Next() {
		var acnt def.BankAccount
		err = rows.Scan(&acnt.DbId, &acnt.LinkedBankDbId, &acnt.Name,
			&acnt.BaseAmount, &acnt.Description)
		if err != nil {
			return
		}
		accounts = append(accounts, acnt)
	}

	return
}
