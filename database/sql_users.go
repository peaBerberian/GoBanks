package database

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

func (gbs *goBanksSql) AddUser(usr DBUserParams) (DBUser, error) {
	values := make([]interface{}, 0)
	values = append(values, usr.Name, usr.PasswordHash,
		usr.Salt, usr.Administrator)

	var id, err = gbs.insertInTable(user_table,
		stripIdField(user_fields), values)

	if err != nil {
		return DBUser{}, databaseQueryError{err.Error()}
	}

	return DBUser{
		Id:            id,
		Name:          usr.Name,
		PasswordHash:  usr.PasswordHash,
		Salt:          usr.Salt,
		Administrator: usr.Administrator,
	}, nil
}

func (gbs *goBanksSql) UpdateUser(id int, fields []string,
	usr DBUserParams) error {

	var values = make([]interface{}, 0)
	var filteredFields = make([]string, 0)

	for _, field := range fields {
		switch field {
		case "Name":
			values = append(values, usr.Name)
			filteredFields = append(filteredFields, bank_fields["Name"])
		case "PasswordHash":
			values = append(values, usr.PasswordHash)
			filteredFields = append(filteredFields, bank_fields["PasswordHash"])
		case "Salt":
			values = append(values, usr.Salt)
			filteredFields = append(filteredFields, bank_fields["Salt"])
		case "Administrator":
			values = append(values, usr.Administrator)
			filteredFields = append(filteredFields, bank_fields["Administrator"])
		}
	}

	return gbs.updateTable(user_table, "", make([]interface{}, 0),
		filteredFields, values)
}

func (gbs *goBanksSql) RemoveUser(id int) error {
	// TODO move that code elsewhere
	_, err := gbs.execQuery("DELETE FROM "+user_table+" WHERE id=?", id)
	return err
}

func (gbs *goBanksSql) GetUser(f DBUserFilters, fields []string) (DBUser, error) {
	var selectString = constructSelectString(user_table,
		filterFields(fields, user_fields))

	// var whereString = "WHERE id=?"
	// var args = make([]interface{}, 0)
	// args = append(args, id)

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

		for _, field := range fields {
			switch field {
			case "Id":
				values = append(values, &usr.Id)
			case "Name":
				values = append(values, &usr.Name)
			case "PasswordHash":
				values = append(values, &usr.PasswordHash)
			case "Salt":
				values = append(values, &usr.Salt)
			case "Administrator":
				values = append(values, &usr.Administrator)
			}
		}

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
	addFilterEq(&conditionString, &args, user_fields["Administrator"],
		f.Administrator)

	return processFilterQuery(conditionString, args, true)
}
