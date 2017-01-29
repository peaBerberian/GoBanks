package database

import (
	"database/sql"
	"strings"

	"fmt"
)

// execQuery is a simple wrapper for the "Exec" sql method.
func (gbs *goBanksSql) execQuery(query string, args ...interface{},
) (sql.Result, error) {
	fmt.Printf("%s : %+v\n", query, args)
	return gbs.db.Exec(query, args...)
}

// getRows is a simple wrapper for the "Query" sql method.
func (gbs *goBanksSql) getRows(query string,
	args ...interface{}) (*sql.Rows, error) {

	fmt.Printf("%s : %+v\n", query, args)
	return gbs.db.Query(query, args...)
}

// joinStringsWithSpace ... joins strings with a space character.
// example: joinStringsWithSpace("foo", "bar", "baz") => "foo bar baz"
func joinStringsWithSpace(queries ...string) string {
	return strings.Join(queries, " ")
}

// constructSelectString constructs the beginning of a SELECT sql request
// based on the wanted fields.
// example: constructSelectString("foo", []{"aa","bb") ->
// "Select aa, bb FROM "foo"
func constructSelectString(tablename string, fields []string) string {
	var str = "SELECT "

	var fieldLength = len(fields)
	if fieldLength > 0 {
		str += fields[0] + " "
	}

	for i := 1; i < fieldLength; i++ {
		str += ", " + fields[i]
	}

	str += " FROM " + tablename
	return str
}

// constructDeleteString constructs the simple beginning of a DELETE sql
// request.
// example: constructDeleteString("foo") -> "DELETE FROM foo"
func constructDeleteString(tablename string) string {
	return "DELETE FROM " + tablename
}

// constructLimitString constructs a simple SQL LIMIT instruction.
// example: constructDeleteString(100) -> "LIMIT 100"
func constructLimitString(limit uint) string {
	return "LIMIT " + string(limit)
}

// updateTable performs a UPDATE request on any database table.
// Here are the arguments it needs:
//    - tablename: the name of the database table
//    - conditions: the "where string", optional
//    - conditionsArgs: arguments for the where string, optional
//    - fields: the wanted fields to update
//    - args: the new values for the wanted fields
func (gbs *goBanksSql) updateTable(tablename string,
	conditions string, conditionsArgs []interface{}, fields []string,
	args []interface{}) (err error) {

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
	valuesStr += " " + conditions

	args = append(args, conditionsArgs...)

	sqlQuery = requestStr + " " + valuesStr

	_, err = gbs.execQuery(sqlQuery, args...)
	return
}

// removeElemFromTable remove sql row(s) based on the table name, a single
// field and its value.
// example: removeElemFromTable("my_table", "id", 5)
func (gbs *goBanksSql) removeElemFromTable(tablename string, field string,
	value interface{}) error {

	_, err := gbs.execQuery("DELETE FROM "+tablename+" WHERE "+field+
		"=?", value)
	return err
}

// insertInTable performs an INSERT INTO SQL call based on the given
// table name, fields and values. It returns the Id of the new given row
// and a possible sql error.
// example: insertInTable("my_table", []string{"foo", "biz"},
// []int{55, 31) => 2, nil
func (gbs *goBanksSql) insertInTable(tablename string, fields []string,
	values []interface{}) (int, error) {
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

	res, err := gbs.execQuery(sqlQuery, values...)

	if err != nil {
		return 0, err
	}
	var id64 int64
	id64, err = res.LastInsertId()
	return int(id64), err
}

// stripIdField browse a map[string]string and delete the one(s) having the
// value "id".
// example: stripIdField(map[string]string{"toto":"foo", "tutu":"id",
// "titi": "bar"}) => map[string]string{"toto": "foo", "titi": "bar"}
func stripIdField(fields map[string]string) []string {
	var res = make([]string, 0)
	for _, field := range fields {
		if field != "id" {
			res = append(res, field)
		}
	}
	return res
}

// Returns a []string by obtaining values from a map[string]string while
// filtering the keys through a []string
func filterFields(fields []string, fieldsMap map[string]string) []string {
	var res = make([]string, 0)
	for _, field := range fields {
		if val, ok := fieldsMap[field]; ok {
			res = append(res, val)
		}
	}
	return res
}

func addConditionEq(cString *string, args *[]interface{}, field string,
	x interface{}) {

	addConditionOperator(cString, args, field, x, "=")
}

func addConditionGEq(cString *string, args *[]interface{}, field string,
	x interface{}) {

	addConditionOperator(cString, args, field, x, ">=")
}

func addConditionLEq(cString *string, args *[]interface{}, field string,
	x interface{}) {

	addConditionOperator(cString, args, field, x, "<=")
}

func addConditionOperator(cString *string, args *[]interface{}, field string,
	x interface{}, operator string) {

	if len(*cString) > 0 {
		*cString += "AND "
	}
	*cString += field + " " + operator + " ? "
	*args = append(*args, x)
}

// /!\ Only works with []int and []string right now
func addConditionOneOf(cString *string, args *[]interface{}, field string,
	x interface{}) bool {

	var temp_str string
	var temp_args []interface{}

	switch x.(type) {
	case []int:
		val, _ := x.([]int)

		if len(val) <= 0 {
			return false
		}
		temp_str, temp_args = addSqlFilterIntArray(field, val...)
	case []string:
		val, _ := x.([]string)

		if len(val) <= 0 {
			return false
		}
		temp_str, temp_args = addSqlFilterStringArray(field, val...)
	}

	if len(*cString) > 0 {
		*cString += "AND "
	}
	*cString += temp_str + " "
	*args = append(*args, temp_args...)
	return true
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
func addSqlFilterStringArray(fieldName string,
	elems ...string) (fstring string, farg []interface{}) {

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
	return fstring, farg
}

func processFilterQuery(conditionString string,
	args []interface{}, ok bool) (string, []interface{}, bool) {

	if !ok {
		return "", nil, false
	}

	if len(conditionString) < 0 {
		return "", nil, true
	}

	return joinStringsWithSpace("WHERE", conditionString), args, true
}

func addFilterOneOf(cString *string, args *[]interface{}, field string,
	filter dbFilterInterface) bool {
	if filter.isFilterActivated() {
		return addConditionOneOf(cString, args, field, filter.getFilterValue())
	}
	return true
}

func addFiltersOneOf(cString *string, args *[]interface{}, fields []string,
	filters ...dbFilterInterface) bool {
	for i, filter := range filters {
		if filter.isFilterActivated() {
			if ok := addConditionOneOf(cString, args, fields[i],
				filter.getFilterValue()); !ok {
				return false
			}
		}
	}
	return true
}

func addFilterGEq(cString *string, args *[]interface{}, field string,
	filter dbFilterInterface) {
	if filter.isFilterActivated() {
		addConditionGEq(cString, args, field, filter.getFilterValue())
	}
}

func addFilterLEq(cString *string, args *[]interface{}, field string,
	filter dbFilterInterface) {
	if filter.isFilterActivated() {
		addConditionLEq(cString, args, field, filter.getFilterValue())
	}
}

func addFilterEq(cString *string, args *[]interface{}, field string,
	filter dbFilterInterface) {
	if filter.isFilterActivated() {
		addConditionEq(cString, args, field, filter.getFilterValue())
	}
}
