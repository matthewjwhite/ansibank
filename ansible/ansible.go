package ansible

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

const playbookBinary = "ansible-playbook"

// PlaybookInvocation is a specific invocation of a playbook.
type PlaybookInvocation struct {
	Path      string
	Arguments []string
}

// PlaybookResult is the result of executing a playbook. Pointer not too necessary
// for invocation, just allows us to point at the original invocation.
type PlaybookResult struct {
	Invocation *PlaybookInvocation
	StartTime  time.Time
	Output     string
}

// Tee executes a playbook while simultaneously writing to both stdout and a buffer.
// A stringified version of the buffer along with other information is returned.
// Based on https://stackoverflow.com/a/62630988.
func (pi PlaybookInvocation) Tee() (*PlaybookResult, error) {
	allArgs := append(pi.Arguments, pi.Path)
	cmd := exec.Command(playbookBinary, allArgs...)
	cmd.Env = append(os.Environ(), "ANSIBLE_FORCE_COLOR=true") // Preserve colors.

	// Obtain pipe. Note that StdoutPipe is a ReadCloser, so we
	// should make sure to defer closing.
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer stdout.Close()

	// Redirect stderr to stdout.
	cmd.Stderr = cmd.Stdout

	// Create TeeReader for printing out to stdout every time we read.
	reader := io.TeeReader(stdout, os.Stdout)

	// Start program.
	startTime := time.Now()
	if err = cmd.Start(); err != nil {
		return nil, err
	}

	outputBuilder := strings.Builder{}

	// Streaming in this way so we're progressively dumping stdout, to preserve the
	// live view of "ansible-playbook". All encodings used should be UTF-8 compatible
	// so not concerned.
	for {
		tmp := make([]byte, 1024)
		_, err := reader.Read(tmp)
		if err != nil {
			break
		}

		// Write to output builder, in preparation for DB INSERT.
		outputBuilder.Write(tmp)
	}

	// If it simply failed due to a failed Ansible run, don't fail this process!
	var e *exec.ExitError
	if err = cmd.Wait(); err != nil && !errors.As(err, &e) {
		return nil, err
	}

	// Output could be huge, return a pointer!
	return &PlaybookResult{
		Invocation: &pi,
		StartTime:  startTime,
		Output:     outputBuilder.String(),
	}, nil
}
