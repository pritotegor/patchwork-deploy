// Package schedule implements deployment time-window enforcement for
// patchwork-deploy.
//
// A Schedule holds one or more Windows. Each Window specifies:
//   - An optional list of allowed weekdays (empty = all days).
//   - A start and end clock time in 24-hour "HH:MM" format (UTC).
//
// Deployments are only permitted when the current UTC time falls inside at
// least one window. When no windows are configured every time is allowed,
// making the schedule opt-in.
//
// Example JSON schedule file:
//
//	{
//	  "windows": [
//	    { "weekdays": ["monday","tuesday","wednesday","thursday","friday"],
//	      "start": "08:00", "end": "20:00" }
//	  ]
//	}
//
// Use LoadFile to parse a schedule from disk, or construct a Schedule
// directly with New for programmatic use.
package schedule
