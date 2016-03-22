package mysql

import "database/sql"
import "time"

import dbt "github.com/peaberberian/GoBanks/database/types"

func (gbs *GoBanksSql) AddBankAccount(bnkA dbt.BankAccount) (id int, err error) {
	var res sql.Result
	res, err = gbs.db.Exec("insert into account "+
		"(bank_id, name, date_added, base_amount) "+
		"values (?, ?)",
		bnkA.BankId,
		bnkA.Name,
		time.Now(),
		bnkA.BaseAmount,
	)
	if err != nil {
		return
	}
	var id64 int64
	id64, err = res.LastInsertId()
	id = int(id64)
	return
}

func (gbs *GoBanksSql) RemoveBankAccount(accountId int) (err error) {
	_, err = gbs.db.Exec("DELETE FROM account "+
		"WHERE id=?", accountId)
	return
}

// TODO
func (gbs *GoBanksSql) UpdateBankAccount(accountId int, bnkA dbt.BankAccount) (err error) {
	return nil
}

// TODO
func (gbs *GoBanksSql) GetBankAccount(accountId int) (t dbt.BankAccount, err error) {
	return t, nil
}

// TODO
func (gbs *GoBanksSql) GetAllBankAccounts() (ts []dbt.BankAccount, err error) {
	return ts, nil
}
