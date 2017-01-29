package database

func (gbs *goBanksSql) AddBank(bnk DBBankParams) (DBBank, error) {
	if bnk.UserId == 0 {
		return DBBank{}, missingInformationsError{"UserId"}
	}

	values := make([]interface{}, 0)
	values = append(values, bnk.UserId, bnk.Name, bnk.Description)

	id, err := gbs.insertInTable(bank_table, stripIdField(bank_fields), values)

	if err != nil {
		return DBBank{}, databaseQueryError{err.Error()}
	}

	return DBBank{
		Id:          id,
		UserId:      bnk.UserId,
		Name:        bnk.Name,
		Description: bnk.Description,
	}, nil
}

// func (gbs *goBanksSql) UpdataBanks(linkedUserId int, name string, description string) (Bank, error) {
func (gbs *goBanksSql) UpdateBanks(f DBBankFilters, fields []string,
	bnk DBBankParams) error {

	var whereString, args, valid = constructBankFilterQuery(f)
	if !valid {
		return nil
	}

	var values = make([]interface{}, 0)
	var filteredFields = make([]string, 0)

	for _, field := range fields {
		switch field {
		case "UserId":
			values = append(values, bnk.UserId)
			filteredFields = append(filteredFields, bank_fields["UserId"])
		case "Name":
			values = append(values, bnk.Name)
			filteredFields = append(filteredFields, bank_fields["Name"])
		case "Description":
			values = append(values, bnk.Description)
			filteredFields = append(filteredFields, bank_fields["Description"])
		}
	}

	return gbs.updateTable(bank_table, whereString, args, filteredFields, values)
}

func (gbs *goBanksSql) RemoveBanks(f DBBankFilters) error {
	var deleteString = constructDeleteString(bank_table)
	var whereString, args, valid = constructBankFilterQuery(f)
	if !valid {
		return nil
	}

	queryString := joinStringsWithSpace(deleteString, whereString)
	_, err := gbs.execQuery(queryString, args...)
	return err
}

func (gbs *goBanksSql) GetBanks(f DBBankFilters, fields []string,
	limit uint) ([]DBBank, error) {

	var selectString = constructSelectString(bank_table,
		filterFields(fields, bank_fields))

	var whereString, args, valid = constructBankFilterQuery(f)
	if !valid {
		return []DBBank{}, nil
	}

	var queryString = joinStringsWithSpace(selectString, whereString)
	if limit != 0 {
		queryString = joinStringsWithSpace(queryString, "LIMIT ?")
		args = append(args, limit)
	}

	rows, err := gbs.getRows(queryString, args...)
	if err != nil {
		return []DBBank{}, databaseQueryError{err.Error()}
	}

	var bnks []DBBank

	for rows.Next() {
		var bnk DBBank

		var values = make([]interface{}, 0)

		for _, field := range fields {
			switch field {
			case "Id":
				values = append(values, &bnk.Id)
			case "UserId":
				values = append(values, &bnk.UserId)
			case "Name":
				values = append(values, &bnk.Name)
			case "Description":
				values = append(values, &bnk.Description)
			}
		}

		if err = rows.Scan(values...); err != nil {
			return []DBBank{}, err
		}

		bnks = append(bnks, bnk)
	}
	return bnks, nil
}

// constructBankFilterQuery takes your filters and returns two elements
// usable for the final sql query:
// - The "WHERE" string
//   For example -> "WHERE id=? AND name=?"
// - An array on interfaces for the sql arguments.
//   For example -> 3, "toto"
// Also returns a boolean if the resulting query is not doable (ex: trying to
// filter users with an empty array of int).
func constructBankFilterQuery(f DBBankFilters) (string, []interface{}, bool) {
	var conditionString string
	var args = make([]interface{}, 0)

	addFilterEq(&conditionString, &args, bank_fields["UserId"], f.UserId)

	var fieldsOneOf = []string{bank_fields["Id"], bank_fields["Name"]}
	ok := addFiltersOneOf(&conditionString, &args, fieldsOneOf,
		f.Ids,
		f.Names,
	)

	return processFilterQuery(conditionString, args, ok)
}
