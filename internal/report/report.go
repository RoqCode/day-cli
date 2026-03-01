package report

import (
	"fmt"
	"sort"
	"time"

	"github.com/roqcode/day/internal/db"
)

func Report(d *db.DB, day time.Time) error {
	pings, err := d.GetPingsForDay(day)
	if err != nil {
		return err
	}

	if len(pings) == 0 {
		fmt.Printf("\nNo pings found for %s. Use 'day ping' to record your first ping for that day!\n\n", day.Format("Monday, 2 January 2006"))
		return nil
	}

	pingsByScope := make(map[string][]db.Ping)

	for _, ping := range pings {
		pingsByScope[ping.Scope] = append(pingsByScope[ping.Scope], ping)
	}

	start, end := pingTimeRange(pings)

	fmt.Printf("\n%v - %v pings (%v - %v)\n---\n", day.Format("Monday, 2 January 2006"), len(pings), start.Format("15:04"), end.Format("15:04"))

	keys := make([]string, 0, len(pingsByScope))
	for k := range pingsByScope {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		scope := pingsByScope[k]
		fmt.Println(scopeLabel(k))
		for _, ping := range scope {
			fmt.Printf("  %v  %s\n", ping.TS.Format("15:04"), ping.Activity)
		}
		fmt.Print("\n")
	}

	return nil
}
