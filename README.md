# WildBerries Data MCP Server

MCP (Model Context Protocol) server for WildBerries data that allows an AI to directly query the database using natural language. The server provides a secure interface for AI to query the database while ensuring query safety through the QueryGuard package.

## Features

- Provides an MCP-compliant API for AI-driven database queries
- Connects to both ClickHouse and PostgreSQL databases
- Ensures query security through QueryGuard validation
- Returns results in JSON format for easy consumption by AI

## Prerequisites

- Go 1.18+
- ClickHouse database
- PostgreSQL database
- Redis (for caching)

## Environment Variables

Configure the application using the following environment variables:

```
# Server configuration
WB_DATA_MCP_IP=0.0.0.0
WB_DATA_MCP_PORT=8081

# Clickhouse configuration
CH_USERNAME=default
CLICKHOUSE_PASSWORD=your_password
CLICKHOUSE_DATABASE=your_database
CH_HOST=localhost
CH_PORT=9000

# PostgreSQL configuration
PSQL_USERNAME=postgres
PG_HOST=localhost
PG_PORT=5432

# Redis configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Telegram logging
TELEGRAM_CHAT_ID=your_chat_id
TELEGRAM_BOT_TOKEN=your_bot_token
```

## Usage

1. Clone the repository:
```
git clone https://github.com/marketconnect/wb_data_mcp.git
```

2. Set environment variables (see above)

3. Build the application:
```
go build -o wb_data_mcp
```

4. Run the server:
```
./wb_data_mcp
```

## Query Example

The server provides a `query_data` tool that AI can use to query the database. Here's an example request:

```json
{
  "table_name": "stocks",
  "fields": ["product_id", "warehouse_id", "quantity", "basic_price"],
  "filters": [
    {"field": "product_id", "operator": "=", "value": "123456"}
  ],
  "limit": 10
}
```

## Allowed Tables

The server supports querying the following tables:

- `stocks`: Product stock information
- `orders`: Order history
- `orders30d`: 30-day order aggregation
- `subjects`: Subject categories
- Various other PostgreSQL tables defined in the system

## Security

QueryGuard ensures that only safe, read-only queries can be executed, protecting your database from potentially harmful operations. The system enforces:

- Only SELECT statements are allowed
- Only pre-approved tables and columns can be queried
- No modifications (INSERT, UPDATE, DELETE) are permitted
- No dangerous SQL constructs (subqueries, CTEs, etc.)

## License

This project is licensed under the MIT License 