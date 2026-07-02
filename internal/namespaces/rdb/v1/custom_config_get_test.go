package rdb_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/rdb/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func samplePostgresConnectionInfo() *rdb.ConnectionInfo {
	return &rdb.ConnectionInfo{
		EngineFamily: rdb.PostgreSQL,
		Host:         "163.172.166.66",
		Port:         21441,
		User:         "myuser",
		Database:     "mydb",
	}
}

func sampleMySQLConnectionInfo() *rdb.ConnectionInfo {
	return &rdb.ConnectionInfo{
		EngineFamily: rdb.MySQL,
		Host:         "163.172.166.66",
		Port:         3306,
		User:         "myuser",
		Database:     "mydb",
	}
}

func Test_RenderPHPConfig(t *testing.T) {
	t.Run("postgresql", func(t *testing.T) {
		got := string(rdb.RenderPHPConfig(samplePostgresConnectionInfo()))
		assert.Contains(
			t,
			got,
			`pg_connect("host=163.172.166.66 port=21441 dbname=mydb user=myuser password=YOUR_PASSWORD")`,
		)
	})

	t.Run("mysql", func(t *testing.T) {
		got := string(rdb.RenderPHPConfig(sampleMySQLConnectionInfo()))
		assert.Contains(t, got, `mysql:host=163.172.166.66;port=3306;dbname=mydb`)
		assert.Contains(t, got, `new PDO($dsn, "myuser", "YOUR_PASSWORD")`)
	})

	t.Run("private network comment", func(t *testing.T) {
		info := samplePostgresConnectionInfo()
		info.PrivateNetworkID = "11111111-1111-1111-1111-111111111111"
		got := string(rdb.RenderPHPConfig(info))
		assert.Contains(t, got, "private network 11111111-1111-1111-1111-111111111111")
	})
}

func Test_RenderNodeConfig(t *testing.T) {
	got := string(rdb.RenderNodeConfig(samplePostgresConnectionInfo()))
	assert.Contains(t, got, `require("pg")`)
	assert.Contains(t, got, `host: "163.172.166.66"`)
	assert.Contains(t, got, `password: process.env.RDB_PASSWORD`)
}

func Test_RenderTypeScriptConfig(t *testing.T) {
	got := string(rdb.RenderTypeScriptConfig(samplePostgresConnectionInfo()))
	assert.Contains(t, got, `import { Pool } from "pg"`)
	assert.Contains(t, got, `database: "mydb"`)
}

func Test_RenderPythonConfig(t *testing.T) {
	t.Run("postgresql", func(t *testing.T) {
		got := string(rdb.RenderPythonConfig(samplePostgresConnectionInfo()))
		assert.Contains(t, got, "import psycopg2")
		assert.Contains(t, got, `sslmode="require"`)
	})

	t.Run("mysql", func(t *testing.T) {
		got := string(rdb.RenderPythonConfig(sampleMySQLConnectionInfo()))
		assert.Contains(t, got, "import mysql.connector")
	})
}

func Test_RenderGoConfig(t *testing.T) {
	t.Run("postgresql", func(t *testing.T) {
		got := string(rdb.RenderGoConfig(samplePostgresConnectionInfo()))
		assert.Contains(t, got, `github.com/jackc/pgx/v5/stdlib`)
		assert.Contains(t, got, "sslmode=require")
	})

	t.Run("mysql", func(t *testing.T) {
		got := string(rdb.RenderGoConfig(sampleMySQLConnectionInfo()))
		assert.Contains(t, got, `github.com/go-sql-driver/mysql`)
	})
}

func Test_RenderRustConfig(t *testing.T) {
	got := string(rdb.RenderRustConfig(samplePostgresConnectionInfo()))
	assert.Contains(t, got, "sqlx::postgres::PgPoolOptions")
	require.Contains(
		t,
		got,
		"postgres://myuser:YOUR_PASSWORD@163.172.166.66:21441/mydb?sslmode=require",
	)
}
