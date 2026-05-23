// Package reporter provides structured, human-readable output for the results
// of patch execution across one or more remote hosts.
//
// Usage:
//
//	r := reporter.New(os.Stdout)
//
//	r.Print(reporter.Result{
//		Host:     "web-01",
//		Patch:    "001_init.sh",
//		Output:   "initialized",
//		Duration: 200 * time.Millisecond,
//	})
//
//	r.Summary(results)
//
// Each Result captures the host, patch filename, combined output, any error,
// and the wall-clock duration of the operation. Summary prints an aggregate
// success/failure count after all patches have been applied.
package reporter
