package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/peaberberian/GoBanks/auth"
	"github.com/peaberberian/GoBanks/database"
)

// DBCategory properties gettable through this handler
var gettable_category_fields = []string{
	"Id",
	"UserId",
	"Name",
	"Description",
	// TODO
	// "ParentId",
}

// handleCategories is the main handler for call on the /categories api. It dispatches
// to other function based on the HTTP method used the typical REST CRUD
// naming scheme.
func handleCategories(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	switch r.Method {
	case "GET":
		handleCategoryRead(w, r, t)
	case "POST":
		handleCategoryCreate(w, r, t)
	case "PUT":
		handleCategoryUpdate(w, r, t)
	case "DELETE":
		handleCategoryDelete(w, r, t)
	default:
		handleNotSupportedMethod(w, r.Method)
	}
}

// handleCategoryRead handle GET requests on the /categories API
func handleCategoryRead(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	// look if we have an id (GET /categories/35 => id == 35)
	var id, hasIdInUrl = getApiId(r.URL.Path)

	var queryString = r.URL.Query()
	var f database.DBCategoryFilters
	var limit int

	// always filter on the current user
	f.UserId.SetFilter(t.UserId)

	// if an id was set in the url, filter to the record corresponding to it
	if hasIdInUrl {
		f.Ids.SetFilter([]int{id})
	} else {
		// if only some ids are wanted, filter
		wantedIds, _ := queryStringPropertyToIntArray(queryString, "id")
		if len(wantedIds) > 0 {
			f.Ids.SetFilter(wantedIds)
		}

		// if only some category names are wanted, filter
		wantedCategoryNames, _ := queryStringPropertyToStringArray(queryString, "name")
		if len(wantedCategoryNames) > 0 {
			f.Names.SetFilter(wantedCategoryNames)
		}

		// TODO
		// // if only some parent ids are wanted, filter
		// wantedParentIds, _ := queryStringPropertyToIntArray(queryString, "pid")
		// if len(wantedParentIds) > 0 {
		// 	f.ParentIds.SetFilter(wantedParentIds)
		// }

		// obtain limit of wanted records, if set
		limit, _ = queryStringPropertyToInt(queryString, "limit")
		fmt.Println("limit", limit)
	}

	// perform the database request
	vals, err := database.GoDB.GetCategories(f, gettable_category_fields, uint(limit))
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	// if an id was given, we're awaiting an object, not an array.
	if hasIdInUrl {
		if len(vals) == 0 {
			fmt.Fprintf(w, "{}")
		} else {
			fmt.Fprintf(w, generateCategoryResponse(vals[0]))
		}
		return
	}

	// else respond directly with the result
	if len(vals) == 0 {
		fmt.Fprintf(w, "[]")
	} else {
		fmt.Fprintf(w, generateCategoriesResponse(vals))
	}
}

// handleCategoryCreate handle POST requests on the /categories API
func handleCategoryCreate(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	// you cannot post on a specific id, reject if you want to do that
	if _, hasIdInUrl := getApiId(r.URL.Path); hasIdInUrl {
		handleNotSupportedMethod(w, r.Method)
		return
	}

	// convert body to map[string]interface{}
	bodyMap, err := readBodyAsStringMap(r.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	// translate data into a DBCategoryParams element
	// (also check mandatory fields)
	categoryElem, err := inputToCategoryParams(bodyMap)
	if err != nil {
		handleError(w, err)
		return
	}

	// attach elem to current user
	categoryElem.UserId = t.UserId

	// perform database add request
	category, err := database.GoDB.AddCategory(categoryElem)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	fmt.Fprintf(w, generateCategoryResponse(category))
}

// handleCategoryUpdate handle PUT requests on the /categories API
func handleCategoryUpdate(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	// look if we have an id (GET /categories/35 => id == 35)
	var id, hasIdInUrl = getApiId(r.URL.Path)

	// if an id was found, it means that we want to replace an element
	// redirect to the right function
	if !hasIdInUrl {
		handleCategoryReplace(w, r, t)
		return
	}

	// check that we can modify this category params
	// (blocking database request here :(, TODO see what I can do, jwt?)
	if err := checkPermissionForCategory(t, id); err != nil {
		handleError(w, err)
		return
	}

	// convert body to map[string]interface{}
	bodyMap, err := readBodyAsStringMap(r.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	// -- check fields and update only the ones there --

	var fields []string
	var categoryElem database.DBCategoryParams

	if val, ok := bodyMap["name"]; ok {
		if str, ok := val.(string); !ok {
			handleError(w, bodyParsingError{})
			return
		} else {
			categoryElem.Name = str
			fields = append(fields, "Name")
		}
	}
	if val, ok := bodyMap["description"]; ok {
		if str, ok := val.(string); !ok {
			handleError(w, bodyParsingError{})
			return
		} else {
			categoryElem.Description = str
			fields = append(fields, "Description")
		}
	}
	// TODO
	// if val, ok := bodyMap["parentId"]; ok {
	// }

	// Filter the category id
	var f database.DBCategoryFilters
	f.Ids.SetFilter([]int{id})

	// perform the database request
	if err = database.GoDB.UpdateCategories(f, fields, categoryElem); err != nil {
		handleError(w, queryOperationError{})
		return
	}

	handleSuccess(w, r)
}

// handleCategoryDelete handle DELETE requests on the /categories API
func handleCategoryDelete(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	// look if we have an id (GET /categories/35 => id == 35)
	var id, hasIdInUrl = getApiId(r.URL.Path)

	var f database.DBCategoryFilters

	// if we have an id, check permission and set filter
	if hasIdInUrl {
		// (blocking database request here :(, TODO see what I can do, jwt?)
		if err := checkPermissionForCategory(t, id); err != nil {
			handleError(w, err)
			return
		}
		f.Ids.SetFilter([]int{id})
	} else {
		// filter by userId
		f.UserId.SetFilter(t.UserId)
	}

	// perform the database request
	if err := database.GoDB.RemoveCategories(f); err != nil {
		handleError(w, queryOperationError{})
		return
	}
	handleSuccess(w, r)
}

// handleCategoryReplace handle specifically PUT requests on the main /categories API
// (not restricted to a certain id).
func handleCategoryReplace(w http.ResponseWriter, r *http.Request,
	t *auth.UserToken) {

	var bodyMaps, err = readBodyAsArrayOfStringMap(r.Body)
	if err != nil {
		handleError(w, queryOperationError{})
		return
	}

	var bnks []database.DBCategoryParams

	// translate data into DBCategoryParams elements
	// (also check mandatory fields)
	for _, bodyMap := range bodyMaps {
		categoryElem, err := inputToCategoryParams(bodyMap)
		if err != nil {
			handleError(w, err)
			return
		}
		bnks = append(bnks, categoryElem)
	}

	// Remove old categories linked to this user
	var f database.DBCategoryFilters
	f.UserId.SetFilter(t.UserId)

	if err := database.GoDB.RemoveCategories(f); err != nil {
		handleError(w, queryOperationError{})
		return
	}

	// add each category indicated to the database
	for _, ctg := range bnks {
		ctg.UserId = t.UserId
		if _, err := database.GoDB.AddCategory(ctg); err != nil {
			handleError(w, queryOperationError{})
			return
		}
	}
	handleSuccess(w, r)
}

// checkPermissionForCategory checks if an user related to the given token has
// the given categoryId. It returns an error if the user doesn't have this category
// or if the database query failed.
func checkPermissionForCategory(t *auth.UserToken, categoryId int) OperationError {
	if val, err := userHasCategory(t.UserId, categoryId); !val {
		return notPermittedOperationError{}
	} else if err != nil {
		return queryOperationError{}
	}
	return nil
}

// userHasCategory checks if the userId given possess the categoryId also given in
// argument. It can return an error if the database query failed.
func userHasCategory(userId int, categoryId int) (bool, error) {
	var f database.DBCategoryFilters

	f.UserId.SetFilter(userId)
	f.Ids.SetFilter([]int{categoryId})

	var fields = gettable_category_fields
	val, err := database.GoDB.GetCategories(f, fields, 0)
	if len(val) == 0 || err != nil {
		return false, err
	}
	return true, nil
}

// generateCategoryResponse generates a JSON string representing the DBCategory
// struct provided for the API user. If the marshalling fails or if the
// result is nil, an empty JSON object is returned ('{}')
func generateCategoryResponse(ctg database.DBCategory) string {
	var resJson = dbCategoryToCategoryJSON(ctg)

	resBytes, err := json.Marshal(resJson)
	if err != nil || resBytes == nil {
		return "{}"
	}
	return string(resBytes)
}

// generateCategoryResponse generates a JSON string representing a collection
// of DBCategory structs provided for the API user. If the marshalling fails or
// if the result is nil, an empty JSON array is returned ('[]')
func generateCategoriesResponse(ctg []database.DBCategory) string {
	var resJson []CategoryJSON
	for _, t := range ctg {
		resJson = append(resJson, dbCategoryToCategoryJSON(t))
	}
	resBytes, err := json.Marshal(resJson)
	if err != nil || resBytes == nil {
		return "[]"
	}
	return string(resBytes)
}

// dbCategoryToCategoryJSON takes a DBCategory and convert it to its corresponding
// CategoryJSON response.
func dbCategoryToCategoryJSON(ctg database.DBCategory) CategoryJSON {
	return CategoryJSON{
		Id:          ctg.Id,
		Name:        ctg.Name,
		Description: ctg.Description,
	}
}

// process map[string]interface{} input to create a DBCategoryParams object.
// if mandatory fields are not found, this function returns an error.
func inputToCategoryParams(input map[string]interface{}) (database.DBCategoryParams, error) {
	var res database.DBCategoryParams
	var valid bool

	// The "name" field is mandatory
	res.Name, valid = input["name"].(string)
	if !valid {
		return res, missingParameterError{"name"}
	}

	res.Description, _ = input["description"].(string)

	// TODO
	// parentIdfl64, _ := input["parendId"].(float64)
	// res.ParentId = int(parentIdfl64)

	return res, nil
}
