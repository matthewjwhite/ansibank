package main

import (
	"fmt"
	"os"

	"github.com/matthewjwhite/ansibank/db"
	"github.com/matthewjwhite/ansibank/playbook"
	"github.com/yookoala/realpath"
)

const (
	dbError = 1 << iota
	playbookError
	pathError
)

func list(d db.DB) error {
	results, err := d.GetResults()
	if err != nil {
		return err
	}

	for _, result := range results {
		fmt.Println(result.Invocation.Path + ":" + result.StartTime.String())
	}

	return nil
}

func main() {
	// Load DB.
	db, err := db.New("ansibank-db")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(dbError)
	}
	defer db.Close()

	if err = db.Init(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(dbError)
	}

	if len(os.Args) == 2 && os.Args[1] == "list" {
		if err = list(db); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(dbError)
		}

		return
	}

	// For now, assume playbook path is the **last** argument.
	path := os.Args[len(os.Args)-1]
	realPath, err := realpath.Realpath(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(pathError)
	}

	invocation := playbook.Invocation{
		Path:      realPath,
		Arguments: os.Args[1 : len(os.Args)-1],
	}

	result, err := invocation.Tee()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(playbookError)
	}

	if err = db.Insert(result); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(dbError)
	}
}
