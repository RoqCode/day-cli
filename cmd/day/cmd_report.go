package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/roqcode/day/internal/report"
	"github.com/spf13/cobra"
)

var (
	back    int
	dateRaw string
)

func strictDate(year int, month time.Month, day int, loc *time.Location) (time.Time, error) {
	t := time.Date(year, month, day, 0, 0, 0, 0, loc)
	if t.Year() != year || t.Month() != month || t.Day() != day {
		return time.Time{}, fmt.Errorf("invalid date")
	}

	return t, nil
}

func parseDateInput(raw string, now time.Time) (time.Time, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return time.Time{}, fmt.Errorf("empty date")
	}

	loc := now.Location()

	if t, err := time.ParseInLocation("2.1.2006", s, loc); err == nil {
		day, err := strictDate(t.Year(), t.Month(), t.Day(), loc)
		if err == nil {
			return day, nil
		}
	}

	if t, err := time.ParseInLocation("2.1.", s, loc); err == nil {
		day, err := strictDate(now.Year(), t.Month(), t.Day(), loc)
		if err == nil {
			return day, nil
		}
	}

	if t, err := time.ParseInLocation("2.", s, loc); err == nil {
		day, err := strictDate(now.Year(), now.Month(), t.Day(), loc)
		if err == nil {
			return day, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid --date %q (use DD., DD.MM. or DD.MM.YYYY)", raw)
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "generate a report for the given day",
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("back") && cmd.Flags().Changed("date") {
			fmt.Print("'--date' and '--back' can not be used together\n\n")
			cmd.Help()
			os.Exit(1)
		}

		if cmd.Flags().Changed("back") && back < 0 {
			fmt.Print("'--back' only accepts values greater than or equal to 0\n\n")
			cmd.Help()
			os.Exit(1)
		}

		now := time.Now()
		day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

		if cmd.Flags().Changed("date") {
			parsed, err := parseDateInput(dateRaw, now)
			if err != nil {
				fmt.Printf("\n%v\n\n", err)
				os.Exit(1)
			}

			day = parsed
		} else if cmd.Flags().Changed("back") {
			day = day.AddDate(0, 0, -back)
		}

		if err := report.Report(database, day); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	reportCmd.Flags().IntVar(&back, "back", 0, "get reports for previous days relative to today. '--back 1' for yesterday etc.")
	reportCmd.Flags().StringVar(&dateRaw, "date", "", "use a date string to get a report for the specified day. use 'DD.', 'DD.MM.' or 'DD.MM.YYYY' as format. if month or year is omitted, the current values are used.")
	rootCmd.AddCommand(reportCmd)
}
