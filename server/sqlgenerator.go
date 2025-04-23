package server

import (
	"fmt"
	"strings"

	"github.com/marketconnect/queryguard"
)

// SQLGeneratorRequest represents a request to generate a SQL query
type SQLGeneratorRequest struct {
	TableName string   `json:"table_name"`
	Fields    []string `json:"fields"`
	Filters   []Filter `json:"filters"`
	Limit     int      `json:"limit"`
}

// Filter represents a condition in a WHERE clause
type Filter struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// SQLGenerator generates and validates SQL queries
type SQLGenerator struct{}

// NewSQLGenerator creates a new SQLGenerator
func NewSQLGenerator() *SQLGenerator {
	return &SQLGenerator{}
}

// GenerateSelectQuery generates a SELECT query from a request and validates it with queryguard
func (g *SQLGenerator) GenerateSelectQuery(req SQLGeneratorRequest) (string, error) {
	if req.TableName == "" {
		return "", fmt.Errorf("table name is required")
	}

	if len(req.Fields) == 0 {
		return "", fmt.Errorf("at least one field is required")
	}

	// Validate table is allowed
	if _, ok := queryguard.AllowedTables[req.TableName]; !ok {
		return "", fmt.Errorf("table '%s' is not allowed", req.TableName)
	}

	// Validate fields are allowed
	for _, field := range req.Fields {
		if !g.isFieldAllowed(req.TableName, field) {
			return "", fmt.Errorf("field '%s' is not allowed for table '%s'", field, req.TableName)
		}
	}

	// Validate filters fields and operators
	for _, filter := range req.Filters {
		if !g.isFieldAllowed(req.TableName, filter.Field) {
			return "", fmt.Errorf("filter field '%s' is not allowed for table '%s'", filter.Field, req.TableName)
		}

		// Convert operator to uppercase
		upperOp := strings.ToUpper(filter.Operator)
		if !isAllowedOperator(upperOp) {
			return "", fmt.Errorf("operator '%s' is not allowed", filter.Operator)
		}
	}

	// Build the query
	query := fmt.Sprintf("SELECT %s FROM %s", strings.Join(req.Fields, ", "), req.TableName)

	// Add WHERE clause if filters exist
	if len(req.Filters) > 0 {
		whereClause := []string{}
		for _, filter := range req.Filters {
			condition := fmt.Sprintf("%s %s '%s'", filter.Field, filter.Operator, filter.Value)
			whereClause = append(whereClause, condition)
		}
		query += " WHERE " + strings.Join(whereClause, " AND ")
	}

	// Add LIMIT clause if specified
	if req.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	// Validate the final query with queryguard
	if err := queryguard.IsSafeSelectQuery(query, 0, false); err != nil {
		return "", fmt.Errorf("query validation failed: %w", err)
	}

	return query, nil
}

// isFieldAllowed checks if a field is allowed for a table
func (g *SQLGenerator) isFieldAllowed(tableName, fieldName string) bool {
	fields, ok := queryguard.AllowedTables[tableName]
	if !ok {
		return false
	}

	for _, field := range fields {
		if field == fieldName {
			return true
		}
	}
	return false
}

// isAllowedOperator checks if an operator is allowed
func isAllowedOperator(op string) bool {
	allowedOperators := map[string]bool{
		"=":      true,
		"<":      true,
		">":      true,
		"<=":     true,
		">=":     true,
		"!=":     true,
		"<>":     true,
		"IN":     true,
		"LIKE":   true,
		"IS":     true,
		"IS NOT": true,
	}
	return allowedOperators[op]
}
