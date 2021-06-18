package db

import (
	"database/sql"
	"strings"

	"time"

	"github.com/matthewjwhite/ansibank/ansible"
	_ "github.com/mattn/go-sqlite3"
)

const runTable = "run"

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

func (d DB) GetResults() ([]*ansible.PlaybookResult, error) {
	row, err := d.Query("SELECT start_time, playbook, args, output FROM " + runTable)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	results := make([]*ansible.PlaybookResult, 0)

	for row.Next() {
		var startTime time.Time
		var playbookPath string
		var args string
		var output string

		row.Scan(&startTime, &playbookPath, &args, &output)

		results = append(results, &ansible.PlaybookResult{
			Invocation: &ansible.PlaybookInvocation{
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
		"CREATE TABLE IF NOT EXISTS " + runTable + " " +
			"(id INTEGER PRIMARY KEY, " +
			"start_time DATETIME DEFAULT CURRENT_TIMESTAMP, " +
			"playbook TEXT, " +
			"args TEXT, " +
			"output TEXT)")
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
func (d DB) Insert(r *ansible.PlaybookResult) error {
	statement, err := d.Prepare("INSERT INTO " + runTable +
		"(start_time, playbook, args, output) VALUES (?, ?, ?, ?)")
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
