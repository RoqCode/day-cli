package report

import (
	"fmt"
	"slices"
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
	var times []time.Time

	for _, ping := range pings {
		times = append(times, ping.TS)
		pingsByScope[ping.Scope]++
	}

	// sort time slice
	slices.SortFunc(times, func(a, b time.Time) int { return a.Compare(b) })

	// sort ping map
	keys := make([]string, 0, len(pingsByScope))

	longestKey := ""
	for key := range pingsByScope {
		if len(key) > len(longestKey) {
			longestKey = key
		}
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return pingsByScope[keys[i]] > pingsByScope[keys[j]]
	})

	fmt.Printf("\nPings today: %v (%v - %v)\n---\n", len(pings), times[0].Format("15:04"), times[len(times)-1].Format("15:04"))

	minPadding := len(longestKey) + 3
	for _, k := range keys {
		if k == "" {
			padding := minPadding - len("no scope")
			out("no scope", padding, pingsByScope[k])
		} else {
			padding := minPadding - len(k)
			out(k, padding, pingsByScope[k])
		}
	}

	return nil
}

func out(keyName string, padding, count int) {
	fmt.Println(fmt.Sprint(keyName, strings.Repeat(" ", padding), count, " ", strings.Repeat("█", count)))
}
