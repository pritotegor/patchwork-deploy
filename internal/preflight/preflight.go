// Package preflight runs pre-deployment checks against remote hosts
// before any patches are applied. Checks verify connectivity, disk space,
// and required command availability.
package preflight

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Check represents a single preflight verification.
type Check struct {
	Name    string
	Command string
}

// Result holds the outcome of a single preflight check.
type Result struct {
	Host  string
	Check string
	OK    bool
	Error error
}

// Runner executes preflight checks against a list of hosts.
type Runner struct {
	out    io.Writer
	checks []Check
}

// Executor is satisfied by ssh.Client and test doubles.
type Executor interface {
	Run(cmd string) (string, error)
}

// DefaultChecks returns the standard set of preflight checks.
var DefaultChecks = []Check{
	{Name: "disk-space", Command: `df -h / | awk 'NR==2{print $5}'`},
	{Name: "bash-available", Command: "bash --version | head -1"},
	{Name: "whoami", Command: "whoami"},
}

// New creates a Runner with the given checks. If checks is nil, DefaultChecks
// are used. Output defaults to os.Stdout.
func New(out io.Writer, checks []Check) *Runner {
	if out == nil {
		out = os.Stdout
	}
	if checks == nil {
		checks = DefaultChecks
	}
	return &Runner{out: out, checks: checks}
}

// RunAll executes all checks against exec for the named host.
// It returns a slice of Results and a non-nil error if any check failed.
func (r *Runner) RunAll(host string, exec Executor) ([]Result, error) {
	results := make([]Result, 0, len(r.checks))
	var firstErr error

	for _, c := range r.checks {
		out, err := exec.Run(c.Command)
		res := Result{
			Host:  host,
			Check: c.Name,
			OK:    err == nil,
			Error: err,
		}
		results = append(results, res)

		status := "OK"
		detail := strings.TrimSpace(out)
		if err != nil {
			status = "FAIL"
			detail = err.Error()
			if firstErr == nil {
				firstErr = fmt.Errorf("preflight check %q failed on %s: %w", c.Name, host, err)
			}
		}
		fmt.Fprintf(r.out, "[preflight] host=%s check=%s status=%s detail=%q\n",
			host, c.Name, status, detail)
	}
	return results, firstErr
}
