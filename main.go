package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matthewjwhite/ansibank/db"
	"github.com/matthewjwhite/ansibank/playbook"
	"github.com/yookoala/realpath"
)

const (
	dbError = 1 << iota
	playbookError
	pathError
	tuiError
)

func listTUI(db db.DB) error {
	// Get path times for the list.
	pathTimes, err := db.GetPathTimes()
	if err != nil {
		return err
	}

	// Initialize with set of PathTimes. Another option is to Init
	// with GetPathTimes and pass an error message to Update, easier
	// to do this for now and about as clean.
	p := tea.NewProgram(listModel{choices: pathTimes, db: db})

	// Start will block until Tea completes, ex. via tea.Quit.
	if err := p.Start(); err != nil {
		return err
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
		if err = listTUI(db); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(tuiError)
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
