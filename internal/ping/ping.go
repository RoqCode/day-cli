package ping

import (
	"fmt"
	"log"
	"time"

	"github.com/roqcode/day/internal/db"
)

func Ping(d *db.DB, activity string, ago int8, at time.Time, silent bool, scope, source string) {
	fmt.Println("recieved Ping")

	if at.IsZero() {
		at = time.Now().Add(time.Duration(ago) * time.Minute)
	} else {
		now := time.Now()
		at = time.Date(now.Year(), now.Month(), now.Day(), at.Hour(), at.Minute(), 0, 0, now.Location())
	}

	ping := db.Ping{
		Activity: activity,
		TS:       at,
		Scope:    scope,
		Source:   source,
	}

	err := d.InsertPing(ping)
	if err != nil {
		log.Fatal(err)
	}
}
