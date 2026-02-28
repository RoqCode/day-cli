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

type Ping struct {
	TS       time.Time
	Activity string
	Scope    string
	Source   string
}

func InitDB(dataDir string) (*DB, error) {
	// check if db path exists, create if not
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not resolve UserHomeDir: %w", err)
	}

	defaultPath := filepath.Join(homeDir, dataDir)

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

func (db *DB) InsertPing(p Ping) error {
	_, err := db.conn.Exec(
		`INSERT INTO pings (ts, activity, scope, source) VALUES (?, ?, ?, ?)`,
		p.TS, p.Activity, p.Scope, p.Source,
	)
	if err != nil {
		return fmt.Errorf("insert ping failed: %w", err)
	}

	return nil
}

func (db *DB) GetPingsForDay(date time.Time) ([]Ping, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.AddDate(0, 0, 1)

	rows, err := db.conn.Query(`SELECT ts, activity, scope, source FROM pings WHERE ts >= ? AND ts < ? ORDER BY scope, ts`, start, end)
	if err != nil {
		return nil, fmt.Errorf("pings for day query failed: %w", err)
	}

	defer rows.Close()

	var pings []Ping
	for rows.Next() {
		var ping Ping
		if err := rows.Scan(&ping.TS, &ping.Activity, &ping.Scope, &ping.Source); err != nil {
			return nil, fmt.Errorf("error during rows scan: %w", err)
		}
		pings = append(pings, ping)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	return pings, nil
}

func (db *DB) GetRecentScopes(n int) ([]string, error) {
	rows, err := db.conn.Query(
		`SELECT scope
		FROM pings
		WHERE scope IS NOT NULL AND scope != ''
		GROUP BY scope
		ORDER BY MAX(ts) DESC
		LIMIT ?`,
		n,
	)
	if err != nil {
		return nil, fmt.Errorf("recent scopes query failed: %w", err)
	}

	defer rows.Close()

	var scopes []string
	for rows.Next() {
		var scope string
		if err := rows.Scan(&scope); err != nil {
			return nil, fmt.Errorf("error during rows scan: %w", err)
		}
		scopes = append(scopes, scope)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	return scopes, nil
}
