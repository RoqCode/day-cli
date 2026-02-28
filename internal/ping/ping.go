package ping

import (
	"fmt"
	"time"

	"github.com/roqcode/day/internal/db"
	"github.com/roqcode/day/internal/scope"
)

func Ping(d *db.DB, activity string, ago int8, at time.Time, silent bool, scopeArg, source string, predefinedScopes []string) error {
	if at.IsZero() {
		at = time.Now().Add(time.Duration(ago) * time.Minute * -1)
	} else {
		now := time.Now()
		at = time.Date(now.Year(), now.Month(), now.Day(), at.Hour(), at.Minute(), 0, 0, now.Location())
	}

	var chosenScope string
	if len(scopeArg) == 0 && !silent {
		scopes, err := d.GetRecentScopes(5)
		if err != nil {
			return err
		}

		scopes = append(scopes, "CUSTOM")

		seen := make(map[string]struct{}, len(scopes))
		for _, s := range scopes {
			seen[s] = struct{}{}
		}
		for _, s := range predefinedScopes {
			if _, ok := seen[s]; !ok {
				scopes = append(scopes, s)
				seen[s] = struct{}{}
			}
		}

		chosenScope, err = scope.GetScope(scopes)
		if err != nil {
			return err
		}
	} else {
		chosenScope = scopeArg
	}

	ping := db.Ping{
		Activity: activity,
		TS:       at,
		Scope:    chosenScope,
		Source:   source,
	}

	insertedPing, err := d.InsertPing(ping)
	if err != nil {
		return err
	}

	if !silent {
		fmt.Printf("recorded \"%s\" at %s (scope: %s, source: %s)\n", insertedPing.Activity, insertedPing.TS.Format("15:04"), insertedPing.Scope, insertedPing.Source)
	}

	return nil
}
