package mysql

import "database/sql"
import "errors"

import def "github.com/peaberberian/GoBanks/database/definitions"

const BANK_TABLE = "bank"

var BANK_FIELD = []string{"user_id", "name", "description"}

func (gbs *goBanksSql) AddBank(bnk def.Bank) (id int, err error) {
	if bnk.LinkedUserDbId == 0 {
		return 0, errors.New("The linked user must be added to the" +
			" database before the bank.")
	}

	values := make([]interface{}, 0)
	values = append(values, bnk.LinkedUserDbId, bnk.Name,
		bnk.Description)

	id, err = gbs.insertInTable(BANK_TABLE, BANK_FIELD, values)
	if err != nil {
		return 0, err
	}
	bnk.DbId = id
	return
}

func (gbs *goBanksSql) RemoveBank(id int) (err error) {
	return gbs.removeIdFromTable(BANK_TABLE, id)
}

func (gbs *goBanksSql) UpdateBank(bnk def.Bank) (err error) {
	values := make([]interface{}, 0)
	values = append(values, bnk.LinkedUserDbId, bnk.Name,
		bnk.Description)

	return gbs.updateInTableFromId(TRANSACTION_TABLE, bnk.DbId,
		TRANSACTION_FIELDS, values)
}

func (gbs *goBanksSql) GetBank(id int) (t def.Bank, err error) {
	row := gbs.getFromTable(TRANSACTION_TABLE, id, TRANSACTION_FIELDS)
	t.DbId = id
	err = row.Scan(&t.LinkedUserDbId, &t.Name, &t.Description)
	return
}

func (gbs *goBanksSql) GetBanks(f def.BankFilters) (bs []def.Bank,
	err error) {
	var queryString string
	var selectString = "select id, user_id, name, description from bank "
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
			return []def.Bank{}, nil
		}
	}
	if f.Filters.Users {
		if len(f.Values.Users) > 0 {
			if atLeastOneFilter {
				whereString += "AND "
			}
			str, arg := addSqlFilterIntArray("user_id", f.Values.Users...)
			whereString += str + " "
			sqlArguments = append(sqlArguments, arg...)
			atLeastOneFilter = true
		} else {
			return []def.Bank{}, nil
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
		return []def.Bank{}, err
	}

	for rows.Next() {
		var bnk def.Bank
		err = rows.Scan(&bnk.DbId, &bnk.LinkedUserDbId, &bnk.Name,
			&bnk.Description)
		if err != nil {
			return
		}
		bs = append(bs, bnk)
	}

	return
}
