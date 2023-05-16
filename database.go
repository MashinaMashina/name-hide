package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const dbFile = "names.db"

type Database struct {
	file string
	db   *sql.DB
}

func NewDatabase(folder string) *Database {
	return &Database{
		file: fmt.Sprintf("%s/%s", folder, dbFile),
	}
}

func (d *Database) FileName() string {
	return d.file
}

func (d *Database) Exists() bool {
	info, err := os.Stat(d.file)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (d *Database) Init() error {
	var err error
	d.db, err = sql.Open("sqlite3", d.file)
	if err != nil {
		return fmt.Errorf("open '%s': %w", d.file, err)
	}

	if err = d.createSchema(); err != nil {
		return fmt.Errorf("create schema: %w", err)
	}

	return nil
}

func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}

	return nil
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

func (d *Database) FreeNum(num int) error {
	if _, err := d.db.Exec("delete from names where space_count=$1", num); err != nil {
		return err
	}

	return nil
}

func (d *Database) List() (List, error) {
	rows, err := d.db.Query("select original, space_count from names")
	if err != nil {
		return nil, err
	}

	res := make(List)
	for rows.Next() {
		var name string
		var spaces int

		if err = rows.Scan(&name, &spaces); err != nil {
			return nil, fmt.Errorf("scan result: %w", err)
		}

		res[spaces] = name
	}

	return res, err
}
