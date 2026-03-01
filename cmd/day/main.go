package main

import (
	"fmt"
	"log"
	"os"

	"github.com/roqcode/day/internal/config"
	"github.com/roqcode/day/internal/db"
	"github.com/spf13/cobra"
)

var (
	c        *config.Config
	database *db.DB
)

func main() {
	Execute()

	if err := database.Close(); err != nil {
		log.Printf("could not close db connection: %v", err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "day",
	Short: "day is a time tracking tool",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error

		c, err = config.GetConfig()
		if err != nil {
			log.Fatal(err)
		}

		database, err = db.InitDB(c.Day.DataDir)
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
