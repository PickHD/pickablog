package helper

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// QueryUpdateBuilder generate update query from map
func QueryUpdateBuilder(tableName string, mapArgs map[string]interface{}, wheres []string) (query string, args []interface{}, err error) {
	query = "UPDATE " + tableName + " SET "

	whereMap := make(map[string]bool)
	for _, where := range wheres {
		whereClean := strings.Replace(where, "AND", "", -1)
		whereClean = strings.Replace(whereClean, "OR", "", -1)
		whereClean = strings.Replace(whereClean, " ", "", -1)
		whereMap[whereClean] = true
	}

	prefix := ""
	identifierNumberMap := make(map[string]int)
	identifier := 1

	// sorted map
	keys := make([]string, 0, len(mapArgs))
	for key := range mapArgs {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		// don't update filter value
		if _, ok := whereMap[key]; ok {
			continue
		}
		id, ok := identifierNumberMap[key]
		if !ok {
			identifierNumberMap[key] = identifier
			id = identifier
			identifier++
		}

		query += prefix + key + " = " + fmt.Sprintf("$%d", id)
		args = append(args, mapArgs[key])

		prefix = ", "
	}

	prefix = ""
	query += " WHERE "
	for _, where := range wheres {
		whereClean := strings.Replace(where, "and", "", -1)
		whereClean = strings.Replace(whereClean, "or", "", -1)
		whereClean = strings.Replace(whereClean, " ", "", -1)
		id, ok := identifierNumberMap[where]
		if !ok {
			identifierNumberMap[whereClean] = identifier
			id = identifier
			identifier++
		}

		if _, ok := mapArgs[whereClean]; !ok {
			return query, args, errors.New("missing filter")
		}

		query += prefix + where + " = " + fmt.Sprintf("$%d", id)
		args = append(args, mapArgs[whereClean])
		prefix = " and "
	}

	return
}
