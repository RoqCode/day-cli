// Seed script – fills the day.db with random pings for the current year
// so you can see what the heatmap looks like when populated.
//
// Usage: go run ./cmd/seed
package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/roqcode/day/internal/config"
	"github.com/roqcode/day/internal/db"
)

var activities = []string{
	"Code review",
	"Standup",
	"Feature implementation",
	"Bug fix",
	"Refactoring",
	"Documentation",
	"Pair programming",
	"Architecture discussion",
	"PR review",
	"Testing",
	"Debugging",
	"Meeting",
	"Planning",
	"Retrospective",
	"Refinement",
	"On-call support",
	"Deployment",
	"Config changes",
	"Research spike",
	"Writing ADR",
}

var scopes = []string{
	"PROJ-42",
	"PROJ-123",
	"PROJ-77",
	"PROJ-201",
	"PROJ-88",
	"daily",
	"meeting",
	"NOTICKET",
	"HOTFIX",
	"BUGFIX",
	"",
}

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	c, err := config.GetConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	d, err := db.InitDB(c.Day.DataDir)
	if err != nil {
		log.Fatalf("could not open db: %v", err)
	}
	defer func() {
		if err := d.Close(); err != nil {
			log.Printf("could not close db: %v", err)
		}
	}()

	now := time.Now()
	firstDay := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, now.Location())

	totalInserted := 0

	for day := firstDay; !day.After(now); day = day.AddDate(0, 0, 1) {
		isWeekend := day.Weekday() == time.Saturday || day.Weekday() == time.Sunday

		var (
			activationChance   float64
			minPings, maxPings int
		)
		if isWeekend {
			activationChance = 0.15
			minPings, maxPings = 1, 3
		} else {
			activationChance = 0.70
			minPings, maxPings = 3, 8
		}

		if rng.Float64() > activationChance {
			continue
		}

		count := minPings + rng.Intn(maxPings-minPings+1)
		for range count {
			// random time between 08:00 and 18:00
			hour := 8 + rng.Intn(10)
			minute := rng.Intn(60)
			ts := time.Date(day.Year(), day.Month(), day.Day(), hour, minute, 0, 0, day.Location())

			activity := activities[rng.Intn(len(activities))]
			scope := scopes[rng.Intn(len(scopes))]

			if _, err := d.InsertPing(db.Ping{
				TS:       ts,
				Activity: activity,
				Scope:    scope,
				Source:   "seed",
			}); err != nil {
				log.Printf("failed to insert ping for %s: %v", day.Format(time.DateOnly), err)
				continue
			}
			totalInserted++
		}
	}

	fmt.Printf("done – inserted %d pings into your day.db\n", totalInserted)
	fmt.Println("run 'day heatmap' to see the result")
}
