package main

import (
	"fmt"
	"log"
	"time"

	"github.com/roqcode/day/internal/db"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	pings := []db.Ping{
		{
			TS:       time.Now(),
			Activity: "commit auf branch: MP-1111",
			Scope:    "MP-1111",
			Source:   "gc",
		},
		{
			TS:       time.Now(),
			Activity: "Refinement",
			Scope:    "Scrum Events",
			Source:   "manual",
		},
		{
			TS:       time.Now(),
			Activity: "branch wechsel: MP-1111 -> origin/develop",
			Scope:    "MP-1111",
			Source:   "gs",
		},
	}

	for _, ping := range pings {
		if err := database.InsertPing(ping); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Ping geschrieben.")
	}

	pingsRes, err := database.GetPingsForDay(time.Now())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nPings heute (%d):\n", len(pingsRes))
	for _, p := range pingsRes {
		fmt.Printf("  [%s] %s (scope: %s, source: %s)\n",
			p.TS.Format("15:04:05"), p.Activity, p.Scope, p.Source)
	}
}
