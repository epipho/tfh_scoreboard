package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// SQLIte is an interface for storing scroes in a standard
// sqlite db
type SQLite struct {
	db *sql.DB
}

// NewSQLIteDB returns a sqlite-backed scores database and creates the schema if needed
func NewSQLiteDB(filename string) (*SQLite, error) {
	s := &SQLite{}
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	s.db = db

	err = s.create()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (db *SQLite) create() error {
	// users table
	_, err := db.db.Exec(`
CREATE TABLE IF NOT EXISTS users (
     name TEXT NOT NULL,
     email TEXT,
     created_at TEXT NOT NULL DEFAULT current_timestamp,
     PRIMARY KEY (name)
)`)

	if err != nil {
		return err
	}

	// scores table
	_, err = db.db.Exec(`
CREATE TABLE IF NOT EXISTS scores (
    name TEXT NOT NULL,
    class TEXT NOT NULL,
    score REAL NOT NULL,
    updated_at TEXT NOT NULL DEFAULT current_timestamp,
    PRIMARY KEY (name, class)
    FOREIGN KEY (name) REFERENCES users (name)
)`)
	if err != nil {
		return err
	}

	// attempts table
	_, err = db.db.Exec(`
CREATE TABLE IF NOT EXISTS attempts (
    name TEXT NOT NULL,
    class TEXT NOT NULL,
    score REAL NOPT NULL,
    created_at TEXT NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (name) REFERENCES users (name)
)`)
	if err != nil {
		return err
	}

	return nil
}

func (db *SQLite) CreateOrUpdateUser(name string, email *string) error {
	upd := "UPDATE users SET email = COALESCE(?, email) WHERE name=?"
	ins := "INSERT INTO users (name, email) VALUES (?, ?)"

	// technically an upsert race here but in practice it won't matter
	res, err := db.db.Exec(upd, email, name)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		_, err = db.db.Exec(ins, name, email)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *SQLite) UpdateScore(name string, class string, score float32, incremental bool, replace bool) error {
	ins := "INSERT INTO scores (name, class, score) VALUES (?, ?, ?)"
	att := "INSERT INTO attempts (name, class, score) VALUES (?, ?, ?)"
	upd_max := "UPDATE scores SET score = max(score, ?), updated_at = datetime('now') WHERE name = ? AND class = ?"
	upd_replace := "UPDATE scores SET score = ?, updated_at = datetime('now') WHERE name = ? AND class = ?"
	upd_inc := "UPDATE scores SET score = score + ?, updated_at = datetime('now') WHERE name = ? and class = ?"

	var res sql.Result
	var err error
	if incremental {
		res, err = db.db.Exec(upd_inc, score, name, class)
		if err != nil {
			return err
		}
	} else if replace {
		res, err = db.db.Exec(upd_replace, score, name, class)
		if err != nil {
			return err
		}
	} else {
		res, err = db.db.Exec(upd_max, score, name, class)
		if err != nil {
			return err
		}
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		_, err = db.db.Exec(ins, name, class, score)
		if err != nil {
			return err
		}
	}

	// insert attempt, ignore the error, no reason to mess up a call if the attempt cant be written
	db.db.Exec(att, name, class, score)
	return nil
}

func (db *SQLite) GetAllScores(class string) error {
	return nil
}
