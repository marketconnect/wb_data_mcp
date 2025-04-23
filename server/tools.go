package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

type contextKey string

const databaseContextKey contextKey = "database"

// QueryTool handles database queries
var QueryTool = mcp.NewTool("query",
	mcp.WithDescription("Execute a database query"),
	mcp.WithString("table_name",
		mcp.Required(),
		mcp.Description("Table to query"),
	),
	mcp.WithString("fields",
		mcp.Required(),
		mcp.Description("Comma-separated list of fields to select"),
	),
	mcp.WithString("filters",
		mcp.Description("JSON array of filter conditions"),
	),
	mcp.WithString("limit",
		mcp.Description("Limit number of results"),
	),
)

// HandleQueryTool implements the query tool handler
func HandleQueryTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get database from context
	db, ok := ctx.Value(databaseContextKey).(*Database)
	if !ok {
		return mcp.NewToolResultError("database not found in context"), nil
	}

	// Parse request parameters
	tableName := request.Params.Arguments["table_name"].(string)
	fieldsStr := request.Params.Arguments["fields"].(string)
	filtersStr := request.Params.Arguments["filters"].(string)
	limitStr := request.Params.Arguments["limit"].(string)

	// Convert fields string to array
	fields := strings.Split(fieldsStr, ",")
	for i := range fields {
		fields[i] = strings.TrimSpace(fields[i])
	}

	// Parse filters JSON
	var filters []Filter
	if filtersStr != "" {
		if err := json.Unmarshal([]byte(filtersStr), &filters); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("invalid filters format: %v", err)), nil
		}
	}

	// Parse limit
	limit := 0
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("invalid limit: %v", err)), nil
		}
	}

	// Generate SQL query
	generator := NewSQLGenerator()
	query, err := generator.GenerateSelectQuery(SQLGeneratorRequest{
		TableName: tableName,
		Fields:    fields,
		Filters:   filters,
		Limit:     limit,
	})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to generate query: %v", err)), nil
	}

	// Execute query based on table type
	var result interface{}
	if IsClickHouseTable(tableName) {
		result, err = db.ClickHouse.Query(ctx, query)
	} else {
		result, err = db.PostgreSQL.Query(ctx, query)
	}

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("query failed: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("%v", result)), nil
}

// RegisterTools registers all tools with the MCP server
func RegisterTools(server *mcpserver.MCPServer, db *Database) {
	// Add query tool
	server.AddTool(QueryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Add database to context
		ctx = context.WithValue(ctx, databaseContextKey, db)
		return HandleQueryTool(ctx, request)
	})
}
