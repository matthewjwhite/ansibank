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

// PathTime is a minimal struct that associates a playbook path with a start time.
type PathTime struct {
	Path      string
	StartTime time.Time
}

// New returns DB struct, given a valid path to the local database file.
func New(path string) (DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return DB{}, err
	}

	return DB{db}, nil
}

// GetPathTimes returns the paths and start times of all previously executed runs.
func (d DB) GetPathTimes() ([]PathTime, error) {
	row, err := d.Query(
		fmt.Sprintf("SELECT %s, %s FROM %s",
			startTimeCol, pathCol, resultTable))
	if err != nil {
		return nil, err
	}
	defer row.Close()

	pathTimes := make([]PathTime, 0)

	for row.Next() {
		var startTime time.Time
		var path string

		err = row.Scan(&startTime, &path)
		if err != nil {
			return nil, err
		}

		pathTimes = append(pathTimes, PathTime{path, startTime})
	}

	return pathTimes, nil
}

// GetOutput returns the playbook output for a particular path and time.
// Will error if more than one that output is returned.
func (d DB) GetOutput(p PathTime) (string, error) {
	row, err := d.Query(
		fmt.Sprintf("SELECT %s FROM %s WHERE %s=? AND %s=?",
			outputCol, resultTable, pathCol, startTimeCol),
		p.Path, p.StartTime)
	if err != nil {
		return "", err
	}
	// It's possible for there to be more than 2.
	defer row.Close()

	if !row.Next() {
		return "", fmt.Errorf("no results for path %s, time %s", p.Path, p.StartTime)
	}

	var output string
	if err = row.Scan(&output); err != nil {
		return "", err
	}

	// Fail if more than one result, https://stackoverflow.com/a/37629440.
	if row.Next() {
		return "", fmt.Errorf("more than one result for path %s, time %s", p.Path, p.StartTime)
	}

	return output, nil
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
