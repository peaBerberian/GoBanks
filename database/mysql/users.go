package mysql

import "database/sql"

import def "github.com/peaberberian/GoBanks/database/definitions"

const USER_TABLE = "user"

var USER_FIELDS = []string{"name", "password", "salt", "admin"}

func (gbs *goBanksSql) UserLength() (len int, err error) {
	gbs.mutex.Lock()
	row := gbs.db.QueryRow("select count(*) from user")
	gbs.mutex.Unlock()
	err = row.Scan(&len)
	if err != nil {
		return
	}
	return
}

func (gbs *goBanksSql) AddUser(user def.User) (id int, err error) {
	values := make([]interface{}, 0)
	values = append(values, user.Name, user.PasswordHash,
		user.Salt, user.Administrator)

	id, err = gbs.insertInTable(TRANSACTION_TABLE, TRANSACTION_FIELDS, values)
	if err != nil {
		return 0, err
	}
	user.DbId = id
	return
}

func (gbs *goBanksSql) RemoveUser(id int) (err error) {
	return gbs.removeIdFromTable(USER_TABLE, id)
}

func (gbs *goBanksSql) UpdateUser(user def.User) (err error) {
	values := make([]interface{}, 0)
	values = append(values, user.Name, user.PasswordHash,
		user.Salt, user.Administrator)

	return gbs.updateInTableFromId(USER_TABLE, user.DbId,
		USER_FIELDS, values)
}

func (gbs *goBanksSql) GetUser(id int) (t def.User, err error) {
	row := gbs.getFromTable(USER_TABLE, id, USER_FIELDS)
	t.DbId = id
	err = row.Scan(&t.Name, &t.PasswordHash, &t.Salt, &t.Administrator)
	return
}

func (gbs *goBanksSql) GetUsers(f def.UserFilters,
) (usrs []def.User, err error) {

	var queryString string
	var selectString = "select id, name, password, salt, administrator" +
		" from user "
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
			return usrs, nil
		}
	}
	if f.Filters.Administrator {
		if atLeastOneFilter {
			whereString += "AND "
		}
		whereString += "administrator=? "
		sqlArguments = append(sqlArguments, f.Values.Administrator)
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
		return []def.User{}, err
	}

	for rows.Next() {
		var usr def.User
		err = rows.Scan(&usr.DbId, &usr.Name, &usr.PasswordHash, &usr.Salt,
			&usr.Administrator)
		if err != nil {
			return
		}
		usrs = append(usrs, usr)
	}

	return
}
