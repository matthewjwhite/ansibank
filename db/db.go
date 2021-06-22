// Package db provides helper functionality for interacting with the database.
package db

import (
	"database/sql"
	"fmt"
	"strings"

	"time"

	"github.com/matthewjwhite/ansibank/playbook"
	_ "github.com/mattn/go-sqlite3"
)

const (
	resultTable  = "results"
	startTimeCol = "start_time"
	pathCol      = "playbook"
	argsCol      = "args"
	outputCol    = "output"
)

// DB is a wrapper for the standard library DB object (via embedding), with some added helpers.
// This effectively emulates inheritance. Note that since this simply contains a pointer, we
// can recreate freely, no need to worry passing this by value.
type DB struct {
	*sql.DB
}

// New returns DB struct, given a valid path to the local database file.
func New(path string) (DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return DB{}, err
	}

	return DB{db}, nil
}

// GetResults returns all previously executed playbook runs.
func (d DB) GetResults() ([]*playbook.Result, error) {
	row, err := d.Query(
		fmt.Sprintf("SELECT %s, %s, %s, %s FROM %s",
			startTimeCol, pathCol, argsCol, outputCol, resultTable))
	if err != nil {
		return nil, err
	}
	defer row.Close()

	results := make([]*playbook.Result, 0)

	for row.Next() {
		var startTime time.Time
		var playbookPath string
		var args string
		var output string

		err = row.Scan(&startTime, &playbookPath, &args, &output)
		if err != nil {
			return nil, err
		}

		results = append(results, &playbook.Result{
			Invocation: &playbook.Invocation{
				Path:      playbookPath,
				Arguments: strings.Split(args, " "),
			},
			StartTime: startTime,
			Output:    output,
		})
	}

	return results, nil
}

// Init creates the required table for storing runs, if it doesn't exist already.
func (d DB) Init() error {
	statement, err := d.Prepare(
		fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
			"(id INTEGER PRIMARY KEY, "+
			"%s DATETIME DEFAULT CURRENT_TIMESTAMP, "+
			"%s TEXT, %s TEXT, %s TEXT)",
			resultTable, startTimeCol, pathCol, argsCol, outputCol))
	if err != nil {
		return err
	}

	_, err = statement.Exec()
	if err != nil {
		return err
	}

	return nil
}

// Insert adds the results of a playbook's execution to the database.
func (d DB) Insert(r *playbook.Result) error {
	statement, err := d.Prepare(
		fmt.Sprintf("INSERT INTO %s (%s, %s, %s, %s) VALUES (?, ?, ?, ?)",
			resultTable, startTimeCol, pathCol, argsCol, outputCol))
	if err != nil {
		return err
	}

	_, err = statement.Exec(r.StartTime, r.Invocation.Path,
		strings.Join(r.Invocation.Arguments, " "), r.Output)
	if err != nil {
		return err
	}

	return nil
}
