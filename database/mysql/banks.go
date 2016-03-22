package mysql

import "database/sql"
import "time"

import dbt "github.com/peaberberian/GoBanks/database/types"

func (gbs *GoBanksSql) AddBank(bnk dbt.Bank) (id int, err error) {
	var res sql.Result
	res, err = gbs.db.Exec("insert into bank "+
		"(name, date_added) "+
		"values (?, ?)",
		bnk.Name,
		time.Now(),
	)
	if err != nil {
		return
	}
	var id64 int64
	id64, err = res.LastInsertId()
	id = int(id64)
	return
}

func (gbs *GoBanksSql) RemoveBank(bankId int) (err error) {
	_, err = gbs.db.Exec("DELETE FROM bank "+
		"WHERE id=?", bankId)
	return
}

// TODO
func (gbs *GoBanksSql) UpdateBank(bankId int, bnk dbt.Bank) (err error) {
	return nil
}

// TODO
func (gbs *GoBanksSql) GetBank(bankId int) (t dbt.Bank, err error) {
	return t, nil
}

// TODO
func (gbs *GoBanksSql) GetAllBanks() (ts []dbt.Bank, err error) {
	return ts, nil
}
