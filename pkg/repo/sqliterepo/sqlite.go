package sqliterepo

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // Import the SQLite driver
)

func Connect(filePath string) (*sql.DB, error) {
	// Open the SQLite database (or create it if it doesn't exist)
	db, err := sql.Open("sqlite", filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(1) // Only 1 connection at a time.
	db.SetMaxIdleConns(1) // Keep one around.
	db.SetConnMaxLifetime(0)

	return db, nil
}

func CreateTables(db *sql.DB) error { //nolint:funlen
	// Create users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		);
	`)
	if err != nil {
		return err
	}

	// Create tokens table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tokens (
			token TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return err
	}

	// Create templates table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS templates (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			image_name TEXT NOT NULL,
			dockerfile TEXT NOT NULL
		);
	`)
	if err != nil {
		return err
	}

	// Create containers table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS containers (
			id TEXT PRIMARY KEY,
			dockerid TEXT NOT NULL,
			imagename TEXT NOT NULL,
			containername TEXT NOT NULL,
			gitrepo TEXT,
			user_id TEXT NOT NULL,
			environment TEXT,
			ports TEXT,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return err
	}

	// Create containers table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS ports (
		port INTEGER PRIMARY KEY,
		in_use BOOLEAN NOT NULL,
		container_id TEXT,
		FOREIGN KEY(container_id) REFERENCES containers(id)
	);
	`)
	if err != nil {
		return err
	}

	return nil
}

func InitData(db *sql.DB, minPort, maxPort int) error {
	if err := fillPorts(db, minPort, maxPort); err != nil {
		return err
	}

	return nil
}

func fillPorts(db *sql.DB, minPort, maxPort int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO ports (port, in_use,container_id) VALUES (?, 0,NULL)")
	if err != nil {
		return err
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			fmt.Printf("failed to close rows: %s", err.Error()) //nolint:forbidigo
		}
	}()

	for p := minPort; p <= maxPort; p++ {
		_, err = stmt.Exec(p)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
