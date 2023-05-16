package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const dbFile = "names.db"

type Database struct {
	db *sql.DB
}

func InitDatabase(folder string) (*Database, error) {
	file := fmt.Sprintf("%s/%s", folder, dbFile)

	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, fmt.Errorf("open '%s': %w", file, err)
	}

	d := &Database{
		db: db,
	}

	if err = d.createSchema(); err != nil {
		return nil, fmt.Errorf("create schema: %w", err)
	}

	return d, err
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) createSchema() error {
	_, err := d.db.Exec(`
		CREATE TABLE IF NOT EXISTS names (
		    original TEXT,
		    space_count INTEGER,
		    created_at INTEGER
		)
	`)

	return err
}

func (d *Database) GetUnusedNum() (int, error) {
	var cnt int

	if err := d.db.QueryRow("select ifnull(max(space_count), 0) from names").Scan(&cnt); err != nil {
		return 0, err
	}

	return cnt + 1, nil
}

func (d *Database) SaveName(originalName string, num int) error {
	q := "insert into names(original, space_count, created_at) values($1, $2, $3)"
	args := []any{
		originalName, num, time.Now().Unix(),
	}

	if _, err := d.db.Exec(q, args...); err != nil {
		return err
	}

	return nil
}
