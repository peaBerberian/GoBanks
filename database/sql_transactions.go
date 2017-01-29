package database

// AddTransaction add a single transaction in the database
// Refer to DBTransactionParams to know the params you should provide
// TODO add userId and account checking (via cache?) here
func (gbs *goBanksSql) AddTransaction(trn DBTransactionParams) (
	DBTransaction,
	error,
) {
	// an accountId is required for every transactions
	if trn.AccountId == 0 {
		return DBTransaction{}, missingInformationsError{"AccountId"}
	}

	values := make([]interface{}, 0)
	values = append(values,
		trn.AccountId,
		trn.Label,
		trn.CategoryId,
		trn.Description,
		trn.TransactionDate,
		trn.RecordDate,
		trn.Debit,
		trn.Credit,
		trn.Reference,
	)

	id, err := gbs.insertInTable(transaction_table,
		stripIdField(transaction_fields), values)
	if err != nil {
		return DBTransaction{}, databaseQueryError{err: err.Error()}
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

// UpdateTransactions update one or multiple transactions in the database
// Basically, you provide:
//   - DBTransactionFilters to know which transaction(s) to update
//   - A list of fields, which are the specific fields you want to update
//   - DBTransactionParams, containing the corresponding fields' value
// TODO add userId and account checking (via cache?) here
func (gbs *goBanksSql) UpdateTransactions(filters DBTransactionFilters,
	updatedFields []string, trn DBTransactionParams) error {

	var whereString, args, valid = constructTransactionFilterQuery(filters)
	if !valid {
		return nil
	}

	var values = make([]interface{}, 0)
	var filteredFields = make([]string, 0)

	// iterate through updatedFields
	// This is ugly as pie but it seems the idiomatic way of doing it
	for _, field := range updatedFields {
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

// RemoveTransaction removes one or multiple transactions from the database
// based on filters
// TODO add userId and account checking (via cache?) here
func (gbs *goBanksSql) RemoveTransactions(filters DBTransactionFilters) error {
	var deleteString = constructDeleteString(transaction_table)
	var whereString, args, valid = constructTransactionFilterQuery(filters)
	if !valid {
		return nil
	}

	var queryString = joinStringsWithSpace(deleteString, whereString)

	_, err := gbs.execQuery(queryString, args...)
	return err
}

// GetTransactions returns one or multiple transactions from the database
// based on filters
// TODO add userId and account checking (via cache?) here
// TODO integrate intelligently fields into the select, maybe add things
// like category (string) / bank (string) / account (string)
func (gbs *goBanksSql) GetTransactions(filters DBTransactionFilters,
	fields []string, limit uint) ([]DBTransaction,
	error) {

	var selectString = constructSelectString(transaction_table,
		filterFields(fields, transaction_fields))

	var whereString, args, valid = constructTransactionFilterQuery(filters)
	if !valid {
		return []DBTransaction{}, nil
	}

	var queryString = joinStringsWithSpace(selectString, whereString)
	if limit != 0 {
		queryString = joinStringsWithSpace(queryString, "LIMIT ?")
		args = append(args, limit)
	}

	rows, err := gbs.getRows(queryString, args...)
	if err != nil {
		return []DBTransaction{}, databaseQueryError{err.Error()}
	}

	var trns []DBTransaction

	for rows.Next() {
		var trn DBTransaction

		var values = make([]interface{}, 0)

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
func constructTransactionFilterQuery(filters DBTransactionFilters) (string,
	[]interface{}, bool) {

	var conditionString string
	var args = make([]interface{}, 0)

	var fieldsOneOf = []string{
		transaction_fields["Id"],
		transaction_fields["AccountId"],
		transaction_fields["CategoryId"],
		transaction_fields["Reference"]}

	ok := addFiltersOneOf(&conditionString, &args, fieldsOneOf,
		filters.Ids,
		filters.AccountIds,
		filters.CategoryIds,
		filters.References)

	addFilterGEq(&conditionString, &args,
		transaction_fields["TransactionDate"], filters.FromTransactionDate)

	addFilterLEq(&conditionString, &args,
		transaction_fields["TransactionDate"], filters.ToTransactionDate)

	addFilterGEq(&conditionString, &args,
		transaction_fields["RecordDate"], filters.FromRecordDate)

	addFilterLEq(&conditionString, &args,
		transaction_fields["RecordDate"], filters.ToRecordDate)

	addFilterGEq(&conditionString, &args,
		transaction_fields["Debit"], filters.MinDebit)

	addFilterLEq(&conditionString, &args,
		transaction_fields["Debit"], filters.MaxDebit)

	addFilterGEq(&conditionString, &args,
		transaction_fields["Credit"], filters.MinCredit)

	addFilterLEq(&conditionString, &args,
		transaction_fields["Credit"], filters.MaxCredit)

	return processFilterQuery(conditionString, args, ok)
}
