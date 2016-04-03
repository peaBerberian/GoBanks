package database

func (gbs *goBanksSql) AddAccount(acc DBAccountParams) (DBAccount, error) {
	if acc.BankId == 0 {
		return DBAccount{}, missingInformationsError{"BankId"}
	}

	values := make([]interface{}, 0)
	values = append(values, acc.BankId, acc.Name, acc.Description)

	id, err := gbs.insertInTable(account_table,
		stripIdField(account_fields), values)

	if err != nil {
		return DBAccount{}, databaseQueryError{err: err.Error()}
	}

	return DBAccount{
		Id:          id,
		BankId:      acc.BankId,
		Name:        acc.Name,
		Description: acc.Description,
	}, nil
}

func (gbs *goBanksSql) UpdateAccounts(f DBAccountFilters,
	fields []string, acc DBAccountParams) error {

	var whereString, args, valid = constructAccountFilterQuery(f)
	if !valid {
		return nil
	}

	var values = make([]interface{}, 0)
	var filteredFields = make([]string, 0)

	// TODO do that with reflection
	for _, field := range fields {
		switch field {
		case "BankId":
			values = append(values, acc.BankId)
			filteredFields = append(filteredFields, account_fields["BankId"])
		case "Name":
			values = append(values, acc.Name)
			filteredFields = append(filteredFields, account_fields["Name"])
		case "Description":
			values = append(values, acc.Description)
			filteredFields = append(filteredFields, account_fields["Description"])
		}
	}

	return gbs.updateTable(account_table, whereString, args, filteredFields, values)
}

func (gbs *goBanksSql) RemoveAccounts(f DBAccountFilters) error {
	var deleteString = constructDeleteString(account_table)
	var whereString, args, valid = constructAccountFilterQuery(f)
	if !valid {
		return nil
	}
	var queryString = joinStringsWithSpace(deleteString, whereString)

	_, err := gbs.execQuery(queryString, args...)
	return err
}

func (gbs *goBanksSql) GetAccounts(f DBAccountFilters,
	fields []string, limit uint) ([]DBAccount, error) {

	var selectString = constructSelectString(account_table,
		filterFields(fields, account_fields))

	var whereString, args, valid = constructAccountFilterQuery(f)
	if !valid {
		return []DBAccount{}, nil
	}

	var queryString = joinStringsWithSpace(selectString, whereString)
	if limit != 0 {
		joinStringsWithSpace(queryString, constructLimitString(limit))
	}

	rows, err := gbs.getRows(queryString, args...)
	if err != nil {
		return []DBAccount{}, databaseQueryError{err.Error()}
	}

	var accs []DBAccount

	for rows.Next() {
		var acc DBAccount

		var values = make([]interface{}, 0)

		// TODO do that with reflection
		for _, field := range fields {
			switch field {
			case "Id":
				values = append(values, &acc.Id)
			case "BankId":
				values = append(values, &acc.BankId)
			case "Name":
				values = append(values, &acc.Name)
			case "Description":
				values = append(values, &acc.Description)
			}
		}

		if err = rows.Scan(values...); err != nil {
			return []DBAccount{}, err
		}

		accs = append(accs, acc)
	}
	return accs, nil
}

// constructAccountFilterQuery takes your filters and returns two elements
// usable for the final sql query:
// - The "WHERE" string
//   For example -> "WHERE id=? AND name=?"
// - An array on interfaces for the sql arguments.
//   For example -> 3, "toto"
// TODO userId somewhere
func constructAccountFilterQuery(f DBAccountFilters) (string,
	[]interface{}, bool) {

	var conditionString string
	var args = make([]interface{}, 0)

	var fieldsOneOf = []string{
		account_fields["Id"],
		account_fields["Name"],
		account_fields["BankId"]}

	ok := addFiltersOneOf(&conditionString, &args, fieldsOneOf,
		f.Ids,
		f.Names,
		f.BankIds)

	return processFilterQuery(conditionString, args, ok)
}
