package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Init initialises the SQLite database, creates schema and returns the connection.
func Init(dataDir string) (*sql.DB, error) {
	if err := os.MkdirAll(dataDir, 0750); err != nil {
		return nil, fmt.Errorf("failed to create data dir: %w", err)
	}
	dbPath := filepath.Join(dataDir, "apguard.db")
	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	DB = db
	log.Printf("[DB] Connected: %s", dbPath)
	if err := migrate(db); err != nil {
		return nil, err
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	schema := `
	PRAGMA foreign_keys = ON;

	CREATE TABLE IF NOT EXISTS users (
		id            INTEGER PRIMARY KEY AUTOINCREMENT,
		username      TEXT    NOT NULL UNIQUE,
		email         TEXT    NOT NULL UNIQUE,
		password_hash TEXT    NOT NULL,
		role          TEXT    NOT NULL DEFAULT 'analyst',
		created_at    DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS targets (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		url         TEXT    NOT NULL UNIQUE,
		name        TEXT    NOT NULL,
		description TEXT,
		enabled     INTEGER NOT NULL DEFAULT 1,
		created_by  INTEGER REFERENCES users(id),
		created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS scans (
		id                   INTEGER PRIMARY KEY AUTOINCREMENT,
		target_id            INTEGER NOT NULL REFERENCES targets(id),
		status               TEXT    NOT NULL DEFAULT 'pending',
		start_time           DATETIME,
		end_time             DATETIME,
		scan_profile         TEXT    NOT NULL DEFAULT 'full',
		total_vulnerabilities INTEGER NOT NULL DEFAULT 0,
		created_at           DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS vulnerabilities (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		scan_id     INTEGER NOT NULL REFERENCES scans(id),
		type        TEXT    NOT NULL,
		severity    TEXT    NOT NULL,
		endpoint    TEXT    NOT NULL,
		parameter   TEXT,
		payload     TEXT,
		description TEXT,
		status      TEXT    NOT NULL DEFAULT 'open',
		created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS audit_logs (
		id        INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id   INTEGER REFERENCES users(id),
		action    TEXT    NOT NULL,
		details   TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS scheduled_scans (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		target_id   INTEGER NOT NULL REFERENCES targets(id),
		cron_expr   TEXT    NOT NULL,
		scan_profile TEXT   NOT NULL DEFAULT 'full',
		enabled     INTEGER NOT NULL DEFAULT 1,
		last_run    DATETIME,
		created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	log.Println("[DB] Schema migration complete")
	return nil
}
