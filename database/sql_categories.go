package database

func (gbs *goBanksSql) AddCategory(ctg DBCategoryParams) (DBCategory, error) {
	if ctg.UserId == 0 {
		return DBCategory{}, missingInformationsError{"UserId"}
	}

	values := make([]interface{}, 0)
	values = append(values, ctg.UserId, ctg.Name, ctg.Description, ctg.ParentId)

	id, err := gbs.insertInTable(category_table, stripIdField(category_fields), values)

	if err != nil {
		return DBCategory{}, databaseQueryError{err.Error()}
	}

	return DBCategory{
		Id:          id,
		UserId:      ctg.UserId,
		Name:        ctg.Name,
		Description: ctg.Description,
		ParentId:    ctg.ParentId,
	}, nil
}

func (gbs *goBanksSql) UpdateCategories(f DBCategoryFilters, fields []string,
	ctg DBCategoryParams) error {

	var whereString, args, valid = constructCategoryFilterQuery(f)
	if !valid {
		return nil
	}

	var filteredFields = make([]string, 0)
	var values = make([]interface{}, 0)

	// update only wanted fields
	for _, field := range fields {
		switch field {
		case "UserId":
			values = append(values, ctg.UserId)
			filteredFields = append(filteredFields, category_fields["UserId"])
		case "Name":
			values = append(values, ctg.Name)
			filteredFields = append(filteredFields, category_fields["Name"])
		case "Description":
			values = append(values, ctg.Description)
			filteredFields = append(filteredFields, category_fields["Description"])
		case "ParentId":
			values = append(values, ctg.ParentId)
			filteredFields = append(filteredFields, category_fields["ParentId"])
		}
	}

	return gbs.updateTable(category_table, whereString, args, filteredFields, values)
}

func (gbs *goBanksSql) RemoveCategories(f DBCategoryFilters) error {
	var deleteString = constructDeleteString(category_table)
	var whereString, args, valid = constructCategoryFilterQuery(f)
	if !valid {
		return nil
	}

	queryString := joinStringsWithSpace(deleteString, whereString)
	_, err := gbs.execQuery(queryString, args...)
	return err
}

func (gbs *goBanksSql) GetCategories(f DBCategoryFilters, fields []string,
	limit uint) ([]DBCategory, error) {

	var selectString = constructSelectString(category_table,
		filterFields(fields, category_fields))

	var whereString, args, valid = constructCategoryFilterQuery(f)
	if !valid {
		return []DBCategory{}, nil
	}

	var queryString = joinStringsWithSpace(selectString, whereString)
	if limit != 0 {
		queryString = joinStringsWithSpace(queryString, "LIMIT ?")
		args = append(args, limit)
	}

	rows, err := gbs.getRows(queryString, args...)
	if err != nil {
		return []DBCategory{}, databaseQueryError{err.Error()}
	}

	var bnks []DBCategory

	for rows.Next() {
		var ctg DBCategory

		var values = make([]interface{}, 0)

		for _, field := range fields {
			switch field {
			case "Id":
				values = append(values, &ctg.Id)
			case "UserId":
				values = append(values, &ctg.UserId)
			case "Name":
				values = append(values, &ctg.Name)
			case "Description":
				values = append(values, &ctg.Description)
			case "ParentId":
				values = append(values, &ctg.ParentId)
			}
		}

		if err = rows.Scan(values...); err != nil {
			return []DBCategory{}, err
		}

		bnks = append(bnks, ctg)
	}
	return bnks, nil
}

// constructCategoryFilterQuery takes your filters and returns two elements
// usable for the final sql query:
// - The "WHERE" string
//   For example -> "WHERE id=? AND name=?"
// - An array on interfaces for the sql arguments.
//   For example -> 3, "toto"
// Also returns a boolean if the resulting query is not doable (ex: trying to
// filter users with an empty array of int).
func constructCategoryFilterQuery(f DBCategoryFilters) (string, []interface{}, bool) {
	var conditionString string
	var args = make([]interface{}, 0)

	addFilterEq(&conditionString, &args, category_fields["UserId"], f.UserId)

	// TODO re-normalize that... It's to hard to comprehend for not much advantages
	var fieldsOneOf = []string{
		category_fields["Id"],
		category_fields["Name"],
		category_fields["ParentId"],
	}
	ok := addFiltersOneOf(&conditionString, &args, fieldsOneOf,
		f.Ids,
		f.Names,
		f.ParentIds,
	)

	return processFilterQuery(conditionString, args, ok)
}
