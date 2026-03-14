package report

import (
	"fmt"
	"strings"
	"time"

	"github.com/roqcode/day/internal/db"
)

func Heatmap(d *db.DB) error {
	start := time.Date(time.Now().Year(), 0, 0, 0, 0, 0, 0, time.Now().Location())
	end := start.AddDate(1, 0, 0)

	days, highestCount, err := d.GetAllPingsInRange(start, end)
	if err != nil {
		return err
	}

	if len(days) == 0 {
		fmt.Printf("\nNo pings found. Use 'day ping' to record your first ping for today!\n\n")
		return nil
	}

	renderYear(days, highestCount)

	return nil
}

func renderYear(days db.DayMap, highestCount int) {
	now := time.Now()
	firstDay := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, now.Location())
	lastDay := time.Date(now.Year()+1, time.January, 1, 0, 0, 0, 0, now.Location()).AddDate(0, 0, -1)

	weekDayFirst := int(firstDay.Weekday())
	weekDayLast := int(lastDay.Weekday())

	weeks := 53

	fmt.Printf("\nHeatmap for %v\n\n", now.Year())

	for i := range 7 {
		for j := range weeks {
			if (j == 0 && i < weekDayFirst) || (j == weeks-1 && i > weekDayLast) {
				fmt.Print(" ")
			} else {
				cellIndex := j*7 + i - weekDayFirst
				cellDay := firstDay.AddDate(0, 0, cellIndex)

				cellDayDate, _ := time.Parse(time.DateOnly, strings.SplitN(cellDay.String(), " ", 2)[0])

				if count, ok := days[cellDayDate]; !ok {
					printSquare(0, highestCount)
				} else {
					printSquare(count, highestCount)
				}

			}
			fmt.Print(" ")
		}
		fmt.Print("\n")
	}
}

func printSquare(count, highestCount int) {
	if count == 0 {
		fmt.Printf("\033[38;5;%dm■\033[0m", 236)
		return
	}

	relativeCount := float64(count) / float64(highestCount)

	color := 0
	if relativeCount > 0 && relativeCount < .26 {
		color = 22
	} else if relativeCount < .51 {
		color = 28
	} else if relativeCount < .76 {
		color = 34
	} else {
		color = 46
	}
	fmt.Printf("\033[38;5;%dm■\033[0m", color)
}
