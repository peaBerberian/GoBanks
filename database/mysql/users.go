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
		"(name, password, salt) "+
		"values (?, ?)",
		user.Name,
		user.PasswordHash,
		user.Salt,
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

func (gbs *GoBanksSql) UpdateUser(userId int, user dbt.User) (err error) {
	gbs.mutex.Lock()
	_, err = gbs.db.Exec("UPDATE user "+
		"SET name=?, password=?, salt=? "+
		"WHERE id=?",
		user.Name,
		user.PasswordHash,
		user.Salt,
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
) (us []dbt.User, err error) {
	return
}
