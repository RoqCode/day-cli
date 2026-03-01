package report

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/roqcode/day/internal/db"
)

func Status(d *db.DB) error {
	pings, err := d.GetPingsForDay(time.Now())
	if err != nil {
		return err
	}

	if len(pings) == 0 {
		fmt.Printf("\nNo pings found for today. Use 'day ping' to record your first ping for today!\n\n")
		return nil
	}

	pingsByScope := make(map[string]int)

	for _, ping := range pings {
		pingsByScope[ping.Scope]++
	}

	start, end := pingTimeRange(pings)

	// sort ping map
	keys := make([]string, 0, len(pingsByScope))

	longestLabel := ""
	for key := range pingsByScope {
		label := scopeLabel(key)
		if len(label) > len(longestLabel) {
			longestLabel = label
		}
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return pingsByScope[keys[i]] > pingsByScope[keys[j]]
	})

	fmt.Printf("\nPings today: %v (%v - %v)\n---\n", len(pings), start.Format("15:04"), end.Format("15:04"))

	minPadding := len(longestLabel) + 3
	for _, k := range keys {
		label := scopeLabel(k)
		padding := minPadding - len(label)
		out(label, padding, pingsByScope[k])
	}

	return nil
}

func out(keyName string, padding, count int) {
	fmt.Println(fmt.Sprint(keyName, strings.Repeat(" ", padding), count, " ", strings.Repeat("█", count)))
}
