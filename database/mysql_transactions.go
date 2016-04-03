package database

func (gbs *goBanksSql) AddTransaction(trn DBTransactionParams) (DBTransaction,
	error) {
	if trn.AccountId == 0 {
		return DBTransaction{}, missingInformationsError{"AccountId"}
	}

	values := make([]interface{}, 0)
	values = append(values, trn.AccountId, trn.Label,
		trn.Debit, trn.Credit, trn.TransactionDate, trn.RecordDate)

	id, err := gbs.insertInTable(transaction_table,
		stripIdField(transaction_fields), values)
	if err != nil {
		return DBTransaction{}, err
	}
	return DBTransaction{
		Id:              id,
		AccountId:       trn.AccountId,
		Label:           trn.Label,
		CategoryId:      trn.CategoryId,
		Description:     trn.Description,
		TransactionDate: trn.TransactionDate,
		RecordDate:      trn.RecordDate,
		Debit:           trn.Debit,
		Credit:          trn.Credit,
		Reference:       trn.Reference,
	}, nil
}

func (gbs *goBanksSql) UpdateTransactions(f DBTransactionFilters,
	fields []string, trn DBTransactionParams) error {

	var whereString, args, valid = constructTransactionFilterQuery(f)
	if !valid {
		return nil
	}

	var values = make([]interface{}, 0)
	var filteredFields = make([]string, 0)

	// TODO do that with reflection
	for _, field := range fields {
		switch field {
		case "AccountId":
			values = append(values, trn.AccountId)
			filteredFields = append(filteredFields, transaction_fields["AccountId"])
		case "Label":
			values = append(values, trn.Label)
			filteredFields = append(filteredFields, transaction_fields["Label"])
		case "CategoryId":
			values = append(values, trn.CategoryId)
			filteredFields = append(filteredFields, transaction_fields["CategoryId"])
		case "Description":
			values = append(values, trn.Description)
			filteredFields = append(filteredFields, transaction_fields["Description"])
		case "TransactionDate":
			values = append(values, trn.TransactionDate)
			filteredFields = append(filteredFields, transaction_fields["TransactionDate"])
		case "RecordDate":
			values = append(values, trn.RecordDate)
			filteredFields = append(filteredFields, transaction_fields["RecordDate"])
		case "Debit":
			values = append(values, trn.Debit)
			filteredFields = append(filteredFields, transaction_fields["Debit"])
		case "Credit":
			values = append(values, trn.Credit)
			filteredFields = append(filteredFields, transaction_fields["Credit"])
		case "Reference":
			values = append(values, trn.Reference)
			filteredFields = append(filteredFields, transaction_fields["Reference"])
		}
	}

	return gbs.updateTable(transaction_table, whereString, args, filteredFields, values)
}

func (gbs *goBanksSql) RemoveTransactions(f DBTransactionFilters) error {
	var deleteString = constructDeleteString(transaction_table)
	var whereString, args, valid = constructTransactionFilterQuery(f)
	if !valid {
		return nil
	}

	var queryString = joinStringsWithSpace(deleteString, whereString)
	_, err := gbs.execQuery(queryString, args...)
	return err
}

func (gbs *goBanksSql) GetTransactions(f DBTransactionFilters,
	fields []string, limit uint) ([]DBTransaction,
	error) {

	var selectString = constructSelectString(transaction_table,
		filterFields(fields, transaction_fields))

	var whereString, args, valid = constructTransactionFilterQuery(f)
	if !valid {
		return []DBTransaction{}, nil
	}

	var queryString = joinStringsWithSpace(selectString, whereString)
	if limit != 0 {
		joinStringsWithSpace(queryString, constructLimitString(limit))
	}

	rows, err := gbs.getRows(queryString, args...)
	if err != nil {
		return []DBTransaction{}, databaseQueryError{err.Error()}
	}

	var trns []DBTransaction

	for rows.Next() {
		var trn DBTransaction

		var values = make([]interface{}, 0)

		// TODO do that with reflection
		for _, field := range fields {
			switch field {
			case "Id":
				values = append(values, &trn.Id)
			case "AccountId":
				values = append(values, &trn.AccountId)
			case "Label":
				values = append(values, &trn.Label)
			case "CategoryId":
				values = append(values, &trn.CategoryId)
			case "Description":
				values = append(values, &trn.Description)
			case "TransactionDate":
				values = append(values, &trn.TransactionDate)
			case "RecordDate":
				values = append(values, &trn.RecordDate)
			case "Debit":
				values = append(values, &trn.Debit)
			case "Credit":
				values = append(values, &trn.Credit)
			case "Reference":
				values = append(values, &trn.Reference)
			}
		}

		if err = rows.Scan(values...); err != nil {
			return []DBTransaction{}, err
		}

		trns = append(trns, trn)
	}
	return trns, nil
}

// constructTransactionFilterQuery takes your filters and returns two elements
// usable for the final sql query:
// - The "WHERE" string
//   For example -> "WHERE id=? AND name=?"
// - An array on interfaces for the sql arguments.
//   For example -> 3, "toto"
// Also returns a boolean if the resulting query is not doable (ex: trying to
// filter users with an empty array of int).
func constructTransactionFilterQuery(f DBTransactionFilters) (string,
	[]interface{}, bool) {

	var conditionString string
	var args = make([]interface{}, 0)

	var fieldsOneOf = []string{
		account_fields["Id"],
		account_fields["AccountId"],
		account_fields["CategoryId"]}

	ok := addFiltersOneOf(&conditionString, &args, fieldsOneOf,
		f.Ids,
		f.AccountIds,
		f.CategoryIds)

	addFilterGEq(&conditionString, &args,
		transaction_fields["TransactionDate"], f.FromTransactionDate)

	addFilterLEq(&conditionString, &args,
		transaction_fields["TransactionDate"], f.ToTransactionDate)

	addFilterGEq(&conditionString, &args,
		transaction_fields["RecordDate"], f.FromRecordDate)

	addFilterLEq(&conditionString, &args,
		transaction_fields["RecordDate"], f.ToRecordDate)

	addFilterGEq(&conditionString, &args,
		transaction_fields["Debit"], f.MinDebit)

	addFilterLEq(&conditionString, &args,
		transaction_fields["Debit"], f.MaxDebit)

	addFilterGEq(&conditionString, &args,
		transaction_fields["Credit"], f.MinCredit)

	addFilterLEq(&conditionString, &args,
		transaction_fields["Credit"], f.MaxCredit)

	return processFilterQuery(conditionString, args, ok)
}
