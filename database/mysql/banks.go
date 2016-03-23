package mysql

import "database/sql"
import "errors"

import dbt "github.com/peaberberian/GoBanks/database/types"

func (gbs *GoBanksSql) AddBank(bnk dbt.Bank) (id int, err error) {
	if bnk.LinkedUserDbId == 0 {
		return 0, errors.New("The linked user must be added to the" +
			" database before the bank.")
	}
	var res sql.Result
	gbs.mutex.Lock()
	res, err = gbs.db.Exec("INSERT INTO bank "+
		"(user_id, name) "+
		"values (?, ?)",
		bnk.LinkedUserDbId,
		bnk.Name,
	)
	gbs.mutex.Unlock()
	if err != nil {
		return 0, err
	}
	var id64 int64
	id64, err = res.LastInsertId()
	id = int(id64)
	bnk.DbId = id
	return
}

func (gbs *GoBanksSql) RemoveBank(bankId int) (err error) {
	gbs.mutex.Lock()
	_, err = gbs.db.Exec("DELETE FROM bank "+
		"WHERE id=?", bankId)
	gbs.mutex.Unlock()
	return
}

func (gbs *GoBanksSql) UpdateBank(bankId int, bnk dbt.Bank) (err error) {
	gbs.mutex.Lock()
	_, err = gbs.db.Exec("UPDATE bank "+
		"SET bank_id=?, name=? "+
		"WHERE id=?",
		bnk.LinkedUserDbId,
		bnk.Name,
		bnk.DbId,
	)
	gbs.mutex.Unlock()
	if err != nil {
		return
	}
	return nil
}

func (gbs *GoBanksSql) GetBank(bankId int) (t dbt.Bank, err error) {
	gbs.mutex.Lock()
	row := gbs.db.QueryRow("select bank_id, name from"+
		" bank where id=?", bankId)
	gbs.mutex.Unlock()
	t.DbId = bankId
	err = row.Scan(&t.LinkedUserDbId, &t.Name)
	if err != nil {
		return dbt.Bank{}, err
	}
	return
}
