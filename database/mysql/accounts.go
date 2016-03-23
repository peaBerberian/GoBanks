package mysql

import "database/sql"
import "errors"

import dbt "github.com/peaberberian/GoBanks/database/types"

func (gbs *GoBanksSql) AddBankAccount(acc dbt.BankAccount) (id int, err error) {
	if acc.LinkedBankDbId == 0 {
		return 0, errors.New("The linked bank must be added to the" +
			" database before the bank account.")
	}
	var res sql.Result
	gbs.mutex.Lock()
	res, err = gbs.db.Exec("INSERT INTO account "+
		"(bank_id, name, base_amount) "+
		"values (?, ?)",
		acc.LinkedBankDbId,
		acc.Name,
		acc.BaseAmount,
	)
	gbs.mutex.Unlock()
	if err != nil {
		return 0, err
	}
	var id64 int64
	id64, err = res.LastInsertId()
	id = int(id64)
	acc.DbId = id
	return
}

func (gbs *GoBanksSql) RemoveBankAccount(accountId int) (err error) {
	gbs.mutex.Lock()
	_, err = gbs.db.Exec("DELETE FROM account "+
		"WHERE id=?", accountId)
	gbs.mutex.Unlock()
	return
}

func (gbs *GoBanksSql) UpdateBankAccount(accountId int, acc dbt.BankAccount) (err error) {
	gbs.mutex.Lock()
	_, err = gbs.db.Exec("UPDATE account "+
		"SET bank_id=?, name=?, base_amount=? "+
		"WHERE id=?",
		acc.LinkedBankDbId,
		acc.Name,
		acc.BaseAmount,
		acc.DbId,
	)
	gbs.mutex.Unlock()
	if err != nil {
		return
	}
	return nil
}

func (gbs *GoBanksSql) GetBankAccount(accountId int) (t dbt.BankAccount, err error) {
	gbs.mutex.Lock()
	row := gbs.db.QueryRow("select bank_id, name, base_amount from"+
		" account where id=?", accountId)
	gbs.mutex.Unlock()
	t.DbId = accountId
	err = row.Scan(&t.LinkedBankDbId, &t.Name, &t.BaseAmount)
	if err != nil {
		return dbt.BankAccount{}, err
	}
	return
}
