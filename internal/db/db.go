// Package db - entry into db handler
package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

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
		return nil, fmt.Errorf("could not resolve UserHomeDir: %v", err)
	}

	defaultPath := filepath.Join(homeDir, dbPath)

	if _, err := os.Stat(defaultPath); os.IsNotExist(err) {
		err = os.MkdirAll(defaultPath, 0o755)
		if err != nil {
			return nil, fmt.Errorf("could not create db path: %v", err)
		}
	}

	d := DB{}
	d.conn, err = sql.Open("sqlite", filepath.Join(defaultPath, "day.db"))
	if err != nil {
		return nil, fmt.Errorf("could not open db: %v", err)
	}

	if err = d.conn.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping db: %v", err)
	}

	return &d, nil
}
