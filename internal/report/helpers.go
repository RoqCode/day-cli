package report

import (
	"time"

	"github.com/roqcode/day/internal/db"
)

func scopeLabel(scope string) string {
	if scope == "" {
		return "no scope"
	}

	return scope
}

// pingTimeRange returns the earliest and latest timestamp in pings.
// pings must not be empty.
func pingTimeRange(pings []db.Ping) (time.Time, time.Time) {
	start := pings[0].TS
	end := pings[0].TS

	for _, ping := range pings[1:] {
		if ping.TS.Before(start) {
			start = ping.TS
		}
		if ping.TS.After(end) {
			end = ping.TS
		}
	}

	return start, end
}
