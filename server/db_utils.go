package server

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

// IsClickHouseTable checks if a table is stored in ClickHouse
func IsClickHouseTable(tableName string) bool {
	clickHouseTables := map[string]bool{
		"stocks":    true,
		"orders":    true,
		"orders30d": true,
	}
	return clickHouseTables[strings.ToLower(tableName)]
}

// QueryClickHouse executes a query against ClickHouse
func QueryClickHouse(ctx context.Context, db *sql.DB, query string) ([]map[string]interface{}, error) {
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return processRows(rows)
}

// QueryPostgreSQL executes a query against PostgreSQL
func QueryPostgreSQL(ctx context.Context, db *sql.DB, query string) ([]map[string]interface{}, error) {
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return processRows(rows)
}

// processRows converts database rows to a slice of maps
func processRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	// Extract column names
	columnNames := make([]string, len(columnTypes))
	for i, ct := range columnTypes {
		columnNames[i] = ct.Name()
	}

	var result []map[string]interface{}

	// Iterate through rows
	for rows.Next() {
		// Create a slice of interface{} to hold the values
		values := make([]interface{}, len(columnNames))
		valuePointers := make([]interface{}, len(columnNames))
		for i := range values {
			valuePointers[i] = &values[i]
		}

		if err := rows.Scan(valuePointers...); err != nil {
			return nil, err
		}

		// Create a map for this row
		row := make(map[string]interface{})
		for i, name := range columnNames {
			row[name] = convertValue(values[i])
		}

		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// convertValue converts database values to appropriate Go types
func convertValue(v interface{}) interface{} {
	if v == nil {
		return nil
	}

	switch value := v.(type) {
	case []byte:
		// Try to convert to string
		return string(value)
	case int64:
		return strconv.FormatInt(value, 10)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", value)
	}
}
