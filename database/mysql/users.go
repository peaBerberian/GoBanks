package mysql

import "database/sql"

import dbt "github.com/peaberberian/GoBanks/database/types"

func (gbs *GoBanksSql) UserLength() (len int, err error) {
	gbs.mutex.Lock()
	row := gbs.db.QueryRow("select count(*) from" +
		" user")
	gbs.mutex.Unlock()
	err = row.Scan(&len)
	if err != nil {
		return
	}
	return
}

func (gbs *GoBanksSql) AddUser(user dbt.User) (id int, err error) {
	var res sql.Result
	gbs.mutex.Lock()
	res, err = gbs.db.Exec("INSERT INTO user "+
		"(name, password, salt, permanent) "+
		"values (?, ?, ?, ?)",
		user.Name,
		user.PasswordHash,
		user.Salt,
		user.Permanent,
	)
	gbs.mutex.Unlock()
	if err != nil {
		return 0, err
	}
	var id64 int64
	id64, err = res.LastInsertId()
	id = int(id64)
	user.DbId = id
	return
}

func (gbs *GoBanksSql) RemoveUser(userId int) (err error) {
	gbs.mutex.Lock()
	_, err = gbs.db.Exec("DELETE FROM user "+
		"WHERE id=?", userId)
	gbs.mutex.Unlock()
	return
}

func (gbs *GoBanksSql) UpdateUser(user dbt.User) (err error) {
	gbs.mutex.Lock()
	_, err = gbs.db.Exec("UPDATE user "+
		"SET name=?, password=?, salt=?, token=? "+
		"WHERE id=?",
		user.Name,
		user.PasswordHash,
		user.Salt,
		user.Token,
		user.DbId,
	)
	gbs.mutex.Unlock()
	if err != nil {
		return
	}
	return nil
}

func (gbs *GoBanksSql) GetUser(userId int) (t dbt.User, err error) {
	gbs.mutex.Lock()
	row := gbs.db.QueryRow("select name, password, salt from"+
		" user where id=?", userId)
	gbs.mutex.Unlock()
	t.DbId = userId
	err = row.Scan(&t.Name, &t.PasswordHash, &t.Salt)
	if err != nil {
		return dbt.User{}, err
	}
	return
}

// TODO filters on Name / Permanent status
func (gbs *GoBanksSql) GetUsers(filters dbt.UserFilters,
) (usrs []dbt.User, err error) {
	gbs.mutex.Lock()
	var rows = new(sql.Rows)
	rows, err = gbs.db.Query("select id, name, password, salt, permanent from user"+
		" WHERE name=?",
		filters.Values.Names[0])
	gbs.mutex.Unlock()

	if err != nil {
		return []dbt.User{}, err
	}

	for rows.Next() {
		var usr dbt.User
		err = rows.Scan(&usr.DbId, &usr.Name, &usr.PasswordHash, &usr.Salt,
			&usr.Permanent)
		if err != nil {
			return
		}
		usrs = append(usrs, usr)
	}

	return
}
