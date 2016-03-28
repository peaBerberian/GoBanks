package mysql

import "database/sql"

func (gbs *goBanksSql) execSqlQuery(query string, args ...interface{},
) (res sql.Result, err error) {
	gbs.mutex.Lock()
	res, err = gbs.db.Exec(query, args...)
	gbs.mutex.Unlock()
	return
}

func (gbs *goBanksSql) getFromTable(tablename string, id int,
	fields []string) (row *sql.Row) {
	var requestStr = "SELECT "

	var fieldsLength = len(fields)
	for i, field := range fields {
		requestStr += field
		if i < fieldsLength-1 {
			requestStr += ", "
		}
	}
	requestStr = " FROM " + tablename + " WHERE id=?"

	gbs.mutex.Lock()
	row = gbs.db.QueryRow(requestStr, id)
	gbs.mutex.Unlock()
	return row
}

func (gbs *goBanksSql) updateInTableFromId(tablename string,
	id int, fields []string, values []interface{}) (err error) {
	var sqlQuery string
	var requestStr string = "UPDATE " + tablename + " SET "
	var valuesStr string = ""

	var fieldsLength = len(fields)
	for i, field := range fields {
		valuesStr += field + "=?"
		if i < fieldsLength-1 {
			valuesStr += ", "
		}
	}
	valuesStr += " WHERE Id=?"
	values = append(values, id)

	sqlQuery = requestStr + " " + valuesStr

	_, err = gbs.execSqlQuery(sqlQuery, values...)
	return
}

func (gbs *goBanksSql) removeIdFromTable(tablename string, id int,
) (err error) {
	_, err = gbs.execSqlQuery("DELETE FROM "+tablename+" WHERE id=?", id)
	return
}

func (gbs *goBanksSql) insertInTable(tablename string, fields []string,
	values []interface{}) (id int, err error) {
	var sqlQuery string
	var requestStr string = "INSERT INTO " + tablename + "("
	var valuesStr string = "values ("
	var res sql.Result

	var fieldsLength = len(fields)
	for i, field := range fields {
		requestStr += field
		valuesStr += "?"
		if i < fieldsLength-1 {
			requestStr += ", "
			valuesStr += ", "
		}
	}
	requestStr += ")"
	valuesStr += ")"
	sqlQuery = requestStr + " " + valuesStr

	res, err = gbs.execSqlQuery(sqlQuery, values...)

	if err != nil {
		return 0, err
	}
	var id64 int64
	id64, err = res.LastInsertId()
	id = int(id64)
	return
}

// addSqlFilterArray construct both the "condition string" and the sql
// arguments for the db.Query method, from the wanted array

// Example of a returned "condition string":
// ( myField = ? OR myField = ? )
// Example of a returned sql arguments array:
// []{myValue, myValue}
//
// Warning: Does not work with array of strings, you need to call
// addSqlStringArray for that.
//
// "Repeated' for lack of generic. Surely not the best solution ever.
// Will see about that later.
func addSqlFilterArray(fieldName string,
	elems ...interface{}) (fstring string, farg []interface{}) {

	var len = len(elems)
	fstring += "( "
	for i, name := range elems {
		farg = append(farg, name)
		fstring += fieldName + " = ? "
		if i < len-1 {
			fstring += "OR "
		}
	}
	fstring += ")"
	return
}

// see addSqlFilterArray
func addSqlFilterIntArray(fieldName string,
	elems ...int) (fstring string, farg []interface{}) {

	var len = len(elems)
	fstring += "( "
	for i, name := range elems {
		farg = append(farg, name)
		fstring += fieldName + " = ? "
		if i < len-1 {
			fstring += "OR "
		}
	}
	fstring += ")"
	return
}

// see addSqlFilterArray
func addSqlFilterStringArray(fieldName string, elems ...string,
) (fstring string, farg []interface{}) {
	var len = len(elems)
	fstring += "( "
	for i, name := range elems {
		farg = append(farg, name)
		fstring += fieldName + " = ? "
		if i < len-1 {
			fstring += "OR "
		}
	}
	fstring += ")"
	return
}
