package rdb

import (
	"fmt"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
)

const rdbPasswordPlaceholder = "YOUR_PASSWORD"

type rdbConfigType string

const (
	rdbConfigTypePHP        rdbConfigType = "php"
	rdbConfigTypeNode       rdbConfigType = "node"
	rdbConfigTypeTypeScript rdbConfigType = "typescript"
	rdbConfigTypePython     rdbConfigType = "python"
	rdbConfigTypeGo         rdbConfigType = "go"
	rdbConfigTypeRust       rdbConfigType = "rust"
)

func renderRDBConfig(configType rdbConfigType, info *ConnectionInfo) (core.RawResult, error) {
	switch configType {
	case rdbConfigTypePHP:
		return RenderPHPConfig(info), nil
	case rdbConfigTypeNode:
		return RenderNodeConfig(info), nil
	case rdbConfigTypeTypeScript:
		return RenderTypeScriptConfig(info), nil
	case rdbConfigTypePython:
		return RenderPythonConfig(info), nil
	case rdbConfigTypeGo:
		return RenderGoConfig(info), nil
	case rdbConfigTypeRust:
		return RenderRustConfig(info), nil
	default:
		return core.RawResult(""), fmt.Errorf("unsupported config type %q", configType)
	}
}

func privateNetworkComment(info *ConnectionInfo) string {
	if info.PrivateNetworkID == "" {
		return ""
	}

	return fmt.Sprintf(
		"// Connect from a resource attached to private network %s\n",
		info.PrivateNetworkID,
	)
}

// RenderPHPConfig renders a PHP database connection snippet.
func RenderPHPConfig(info *ConnectionInfo) core.RawResult {
	lines := []string{privateNetworkComment(info)}

	switch info.EngineFamily {
	case PostgreSQL:
		lines = append(lines,
			"<?php",
			"",
			fmt.Sprintf(
				`$connection = pg_connect("host=%s port=%d dbname=%s user=%s password=%s");`,
				info.Host,
				info.Port,
				info.Database,
				info.User,
				rdbPasswordPlaceholder,
			),
			"if (!$connection) {",
			`    echo "Connection failed\n";`,
			"    exit(1);",
			"}",
			"",
		)
	case MySQL:
		lines = append(lines,
			"<?php",
			"",
			fmt.Sprintf(
				`$dsn = "mysql:host=%s;port=%d;dbname=%s";`,
				info.Host,
				info.Port,
				info.Database,
			),
			fmt.Sprintf(
				`$pdo = new PDO($dsn, "%s", "%s");`,
				info.User,
				rdbPasswordPlaceholder,
			),
			"$pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);",
			"",
		)
	}

	return core.RawResult(strings.Join(lines, "\n"))
}

// RenderNodeConfig renders a Node.js database connection snippet.
func RenderNodeConfig(info *ConnectionInfo) core.RawResult {
	lines := []string{privateNetworkComment(info)}

	switch info.EngineFamily {
	case PostgreSQL:
		lines = append(lines,
			`const { Pool } = require("pg");`,
			"",
			"const pool = new Pool({",
			fmt.Sprintf(`  host: "%s",`, info.Host),
			fmt.Sprintf(`  port: %d,`, info.Port),
			fmt.Sprintf(`  user: "%s",`, info.User),
			`  password: process.env.RDB_PASSWORD,`,
			fmt.Sprintf(`  database: "%s",`, info.Database),
			`  ssl: { rejectUnauthorized: true },`,
			"});",
			"",
			"module.exports = { pool };",
			"",
		)
	case MySQL:
		lines = append(lines,
			`const mysql = require("mysql2/promise");`,
			"",
			"const pool = mysql.createPool({",
			fmt.Sprintf(`  host: "%s",`, info.Host),
			fmt.Sprintf(`  port: %d,`, info.Port),
			fmt.Sprintf(`  user: "%s",`, info.User),
			`  password: process.env.RDB_PASSWORD,`,
			fmt.Sprintf(`  database: "%s",`, info.Database),
			`  ssl: { rejectUnauthorized: true },`,
			"});",
			"",
			"module.exports = { pool };",
			"",
		)
	}

	return core.RawResult(strings.Join(lines, "\n"))
}

// RenderTypeScriptConfig renders a TypeScript database connection snippet.
func RenderTypeScriptConfig(info *ConnectionInfo) core.RawResult {
	lines := []string{privateNetworkComment(info)}

	switch info.EngineFamily {
	case PostgreSQL:
		lines = append(lines,
			`import { Pool } from "pg";`,
			"",
			"export const pool = new Pool({",
			fmt.Sprintf(`  host: "%s",`, info.Host),
			fmt.Sprintf(`  port: %d,`, info.Port),
			fmt.Sprintf(`  user: "%s",`, info.User),
			`  password: process.env.RDB_PASSWORD,`,
			fmt.Sprintf(`  database: "%s",`, info.Database),
			`  ssl: { rejectUnauthorized: true },`,
			"});",
			"",
		)
	case MySQL:
		lines = append(lines,
			`import mysql from "mysql2/promise";`,
			"",
			"export const pool = mysql.createPool({",
			fmt.Sprintf(`  host: "%s",`, info.Host),
			fmt.Sprintf(`  port: %d,`, info.Port),
			fmt.Sprintf(`  user: "%s",`, info.User),
			`  password: process.env.RDB_PASSWORD,`,
			fmt.Sprintf(`  database: "%s",`, info.Database),
			`  ssl: { rejectUnauthorized: true },`,
			"});",
			"",
		)
	}

	return core.RawResult(strings.Join(lines, "\n"))
}

// RenderPythonConfig renders a Python database connection snippet.
func RenderPythonConfig(info *ConnectionInfo) core.RawResult {
	lines := []string{privateNetworkComment(info)}

	switch info.EngineFamily {
	case PostgreSQL:
		lines = append(lines,
			"import psycopg2",
			"",
			"conn = psycopg2.connect(",
			fmt.Sprintf(`    host="%s",`, info.Host),
			fmt.Sprintf(`    port=%d,`, info.Port),
			fmt.Sprintf(`    user="%s",`, info.User),
			fmt.Sprintf(`    password="%s",`, rdbPasswordPlaceholder),
			fmt.Sprintf(`    dbname="%s",`, info.Database),
			`    sslmode="require",`,
			")",
			"",
		)
	case MySQL:
		lines = append(lines,
			"import mysql.connector",
			"",
			"conn = mysql.connector.connect(",
			fmt.Sprintf(`    host="%s",`, info.Host),
			fmt.Sprintf(`    port=%d,`, info.Port),
			fmt.Sprintf(`    user="%s",`, info.User),
			fmt.Sprintf(`    password="%s",`, rdbPasswordPlaceholder),
			fmt.Sprintf(`    database="%s",`, info.Database),
			`    ssl_disabled=False,`,
			")",
			"",
		)
	}

	return core.RawResult(strings.Join(lines, "\n"))
}

// RenderGoConfig renders a Go database connection snippet.
func RenderGoConfig(info *ConnectionInfo) core.RawResult {
	lines := []string{privateNetworkComment(info)}

	switch info.EngineFamily {
	case PostgreSQL:
		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=require",
			info.User,
			rdbPasswordPlaceholder,
			info.Host,
			info.Port,
			info.Database,
		)
		lines = append(lines,
			"package main",
			"",
			"import (",
			`    "database/sql"`,
			`    _ "github.com/jackc/pgx/v5/stdlib"`,
			")",
			"",
			"func main() {",
			fmt.Sprintf(`    db, err := sql.Open("pgx", "%s")`, dsn),
			"    if err != nil {",
			`        panic(err)`,
			"    }",
			"    defer db.Close()",
			"}",
			"",
		)
	case MySQL:
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?tls=true",
			info.User,
			rdbPasswordPlaceholder,
			info.Host,
			info.Port,
			info.Database,
		)
		lines = append(lines,
			"package main",
			"",
			"import (",
			`    "database/sql"`,
			`    _ "github.com/go-sql-driver/mysql"`,
			")",
			"",
			"func main() {",
			fmt.Sprintf(`    db, err := sql.Open("mysql", "%s")`, dsn),
			"    if err != nil {",
			`        panic(err)`,
			"    }",
			"    defer db.Close()",
			"}",
			"",
		)
	}

	return core.RawResult(strings.Join(lines, "\n"))
}

// RenderRustConfig renders a Rust database connection snippet.
func RenderRustConfig(info *ConnectionInfo) core.RawResult {
	lines := []string{privateNetworkComment(info)}

	switch info.EngineFamily {
	case PostgreSQL:
		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=require",
			info.User,
			rdbPasswordPlaceholder,
			info.Host,
			info.Port,
			info.Database,
		)
		lines = append(lines,
			"use sqlx::postgres::PgPoolOptions;",
			"",
			"#[tokio::main]",
			"async fn main() -> Result<(), sqlx::Error> {",
			fmt.Sprintf(`    let pool = PgPoolOptions::new().connect("%s").await?;`, dsn),
			"    let _ = pool;",
			"    Ok(())",
			"}",
			"",
		)
	case MySQL:
		dsn := fmt.Sprintf(
			"mysql://%s:%s@%s:%d/%s",
			info.User,
			rdbPasswordPlaceholder,
			info.Host,
			info.Port,
			info.Database,
		)
		lines = append(lines,
			"use sqlx::mysql::MySqlPoolOptions;",
			"",
			"#[tokio::main]",
			"async fn main() -> Result<(), sqlx::Error> {",
			fmt.Sprintf(`    let pool = MySqlPoolOptions::new().connect("%s").await?;`, dsn),
			"    let _ = pool;",
			"    Ok(())",
			"}",
			"",
		)
	}

	return core.RawResult(strings.Join(lines, "\n"))
}
