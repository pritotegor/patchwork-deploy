package schedule

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// fileSchedule is the JSON shape for a schedule config file.
type fileSchedule struct {
	Windows []struct {
		Weekdays []string `json:"weekdays"`
		Start    string   `json:"start"`
		End      string   `json:"end"`
	} `json:"windows"`
}

var weekdayNames = map[string]time.Weekday{
	"sunday": time.Sunday, "monday": time.Monday, "tuesday": time.Tuesday,
	"wednesday": time.Wednesday, "thursday": time.Thursday,
	"friday": time.Friday, "saturday": time.Saturday,
}

// LoadFile reads a JSON schedule file and returns a Schedule.
// out is passed to New; nil defaults to os.Stdout.
func LoadFile(path string, out interface{ Write([]byte) (int, error) }) (*Schedule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("schedule: read %s: %w", path, err)
	}
	var fs fileSchedule
	if err := json.Unmarshal(data, &fs); err != nil {
		return nil, fmt.Errorf("schedule: parse %s: %w", path, err)
	}
	var windows []Window
	for i, fw := range fs.Windows {
		var days []time.Weekday
		for _, name := range fw.Weekdays {
			d, ok := weekdayNames[name]
			if !ok {
				return nil, fmt.Errorf("schedule: window %d: unknown weekday %q", i, name)
			}
			days = append(days, d)
		}
		windows = append(windows, Window{
			Weekdays: days,
			Start:    fw.Start,
			End:      fw.End,
		})
	}
	return New(windows, out), nil
}
