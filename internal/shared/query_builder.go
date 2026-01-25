package shared

import (
	"fmt"
	"sort"
	"strings"
)

type QueryBuilder struct {
	builder *strings.Builder
}

func NewQueryBuilder() *QueryBuilder {

	return &QueryBuilder{
		builder: &strings.Builder{},
	}
}

/*
* Builds SQL filters from a map of field names to values.
* Only includes filters for non-nil values.
* Appends the values to the provided args slice and updates the argIndex accordingly.
 */
func (qb *QueryBuilder) BuildFilters(fieldMap map[string]interface{}, argIndex *int, args *[]interface{}) string {
	keys := make([]string, 0, len(fieldMap))
	for field := range fieldMap {
		keys = append(keys, field)
	}
	sort.Strings(keys)

	for _, field := range keys {
		if value := fieldMap[field]; value != nil {
			fmt.Fprintf(qb.builder, " AND %s = $%d", field, *argIndex)
			*args = append(*args, value)
			*argIndex++
		}
	}
	return qb.builder.String()
}

/**
* Builds the SET clause for an UPDATE query from a map of field names to values.
* Appends the values to the provided args slice and updates the argIndex accordingly.
* Returns the SET clause as a string.
 */
func (qb *QueryBuilder) BuildUpdates(fieldMap map[string]interface{}, argIndex *int, args *[]interface{}) string {
	keys := make([]string, 0, len(fieldMap))
	for field := range fieldMap {
		keys = append(keys, field)
	}
	sort.Strings(keys)

	var setParts []string
	for _, field := range keys {
		if value := fieldMap[field]; value != nil {
			setParts = append(setParts, fmt.Sprintf("%s = $%d", field, *argIndex))
			*args = append(*args, value)
			*argIndex++
		}
	}
	return strings.Join(setParts, ", ")
}

type OffsetPaginationResult struct {
	TotalCount  int
	TotalPages  int
	CurrentPage int
}

/**
* Builds the LIMIT and OFFSET clauses for an OFFSET pagination query.
* Appends the values to the provided args slice and updates the argIndex accordingly.
* Returns the LIMIT and OFFSET clauses as a string.
 */
func (qb *QueryBuilder) BuildOffsetPagination(offset, limit int, argIndex *int, args *[]interface{}) string {
	var clauses []string

	if limit > 0 {
		clauses = append(clauses, fmt.Sprintf("LIMIT $%d", *argIndex))
		*args = append(*args, limit)
		*argIndex++
	}

	if offset > 0 {
		clauses = append(clauses, fmt.Sprintf("OFFSET $%d", *argIndex))
		*args = append(*args, offset)
		*argIndex++
	}

	if len(clauses) > 0 {
		return " " + strings.Join(clauses, " ")
	}
	return ""
}
