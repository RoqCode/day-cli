package main

import (
	"fmt"
	"log"
	"os"

	"github.com/roqcode/day/internal/db"
	"github.com/spf13/cobra"
)

var database *db.DB

func main() {
	Execute()

	// pingsRes, err := database.GetPingsForDay(time.Now())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// fmt.Printf("\nPings heute (%d):\n", len(pingsRes))
	// for _, p := range pingsRes {
	// 	fmt.Printf("  [%s] %s (scope: %s, source: %s)\n",
	// 		p.TS.Format("15:04:05"), p.Activity, p.Scope, p.Source)
	// }

	database.Close()
}

var rootCmd = &cobra.Command{
	Use:   "day",
	Short: "day is a time tracking tool",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		database, err = db.InitDB()
		if err != nil {
			log.Fatal(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing day '%s'\n", err)
		os.Exit(1)
	}
}
