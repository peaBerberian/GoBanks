package database

// UserLength get the number or users in the database
func (gbs *goBanksSql) UserLength() (int, error) {
	var len int
	row, err := gbs.getRows("select count(*) from usr")
	if err != nil {
		return 0, err
	}

	if err := row.Scan(&len); err != nil {
		return 0, err
	}

	return len, nil
}

// AddUser add a new user to the database
func (gbs *goBanksSql) AddUser(usr DBUserParams) (DBUser, error) {
	// get infos
	userName := usr.Name.getParamValue()
	passwordHash := usr.PasswordHash.getParamValue()
	salt := usr.Salt.getParamValue()

	// constructs []interface{} of inserted values in the right order
	values := make([]interface{}, 0)
	values = append(values,
		username,
		passwordHash,
		salt,
	)

	// add to table, get the UserId
	var id, err = gbs.insertInTable(user_table, stripIdField(user_fields), values)
	if err != nil {
		return DBUser{}, databaseQueryError{err.Error()}
	}

	// construct DBUser from infos
	return DBUser{
		Id:           id,
		Name:         userName,
		PasswordHash: passwordHash,
		Salt:         salt,
	}, nil
}

// UpdateUser update a single user in the database, based on its UserId
func (gbs *goBanksSql) UpdateUser(id int, usr DBUserParams) error {
	values := make([]interface{}, 0)
	filteredFields := make([]string, 0)

	if usr.Name.isParamActivated {
		filteredFields = append(filteredFields, users_fields["Name"])
		values = append(values, usr.Name.getParamValue())
	}
	if usr.PasswordHash.isParamActivated {
		filteredFields = append(filteredFields, users_fields["PasswordHash"])
		values = append(values, usr.passwordHash.getParamValue())
	}
	if usr.Salt.isParamActivated {
		filteredFields = append(filteredFields, users_fields["Salt"])
		values = append(values, usr.Salt.getParamValue())
	}

	return gbs.updateTable(user_table, "", make([]interface{}, 0),
		filteredFields, values)
}

// RemoveUser remove a single user based on its UserId
func (gbs *goBanksSql) RemoveUser(id int) error {
	_, err := gbs.execQuery("DELETE FROM "+user_table+" WHERE id=?", id)
	return err
}

// GetUser get one user based on filters
func (gbs *goBanksSql) GetUser(f DBUserFilters) (DBUser, error) {
	var selectString = constructSelectString(user_table, getFields(user_fields))

	var whereString, args, valid = constructUserFilterQuery(f)
	if !valid {
		return DBUser{}, nil
	}

	var queryString = joinStringsWithSpace(selectString, whereString)

	rows, err := gbs.getRows(queryString, args...)
	if err != nil {
		return DBUser{}, databaseQueryError{err.Error()}
	}

	for rows.Next() {
		var usr DBUser

		var values = make([]interface{}, 0)
		values = append(values,
			&usr.Id,
			&usr.Name,
			&usr.PasswordHash,
			&usr.Salt,
		)

		if err = rows.Scan(values...); err != nil {
			return DBUser{}, err
		}

		return usr, nil
	}
	return DBUser{}, nil
}

// constructUserFilterQuery takes your filters and returns two elements
// usable for the final sql query:
// - The "WHERE" string
//   For example -> "WHERE id=? AND name=?"
// - An array on interfaces for the sql arguments.
//   For example -> 3, "toto"
// Also returns a boolean if the resulting query is not doable (ex: trying to
// filter users with an empty array of int).
func constructUserFilterQuery(f DBUserFilters) (string,
	[]interface{}, bool) {

	var conditionString string
	var args = make([]interface{}, 0)

	addFilterEq(&conditionString, &args, user_fields["Id"], f.Id)
	addFilterEq(&conditionString, &args, user_fields["Name"], f.Name)

	return processFilterQuery(conditionString, args, true)
}
