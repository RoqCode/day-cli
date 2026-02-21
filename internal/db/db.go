// Package db - entry into db handler
package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
}

const dbPath = ".local/share/day"

func InitDB() (*DB, error) {
	// check if db path exists, create if not
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not resolve UserHomeDir: %w", err)
	}

	defaultPath := filepath.Join(homeDir, dbPath)

	if _, err := os.Stat(defaultPath); os.IsNotExist(err) {
		err = os.MkdirAll(defaultPath, 0o755)
		if err != nil {
			return nil, fmt.Errorf("could not create db path: %w", err)
		}
	}

	d := DB{}
	d.conn, err = sql.Open("sqlite", filepath.Join(defaultPath, "day.db"))
	if err != nil {
		return nil, fmt.Errorf("could not open db: %w", err)
	}

	if err = d.conn.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping db: %w", err)
	}

	if err = migrate(&d); err != nil {
		return nil, err
	}

	return &d, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func migrate(db *DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS pings (
			id       INTEGER PRIMARY KEY AUTOINCREMENT,
			ts       DATETIME NOT NULL,
			activity TEXT NOT NULL,
			scope    TEXT,
			source   TEXT NOT NULL DEFAULT 'manual'
		);`,
		`CREATE INDEX IF NOT EXISTS idx_pings_ts    ON pings(ts);`,
		`CREATE INDEX IF NOT EXISTS idx_pings_scope ON pings(scope);`,
	}

	for _, s := range stmts {
		_, err := db.conn.Exec(s)
		if err != nil {
			return fmt.Errorf("could not migrate db: %w", err)
		}
	}

	return nil
}

func (db *DB) InsertPing(ts time.Time, activity, scope, source string) error {
	_, err := db.conn.Exec(
		`INSERT INTO pings (ts, activity, scope, source) VALUES (?, ?, ?, ?)`,
		ts, activity, scope, source,
	)
	if err != nil {
		return fmt.Errorf("insert ping failed: %w", err)
	}

	return nil
}
