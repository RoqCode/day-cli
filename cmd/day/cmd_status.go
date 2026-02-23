package main

import (
	"log"

	"github.com/roqcode/day/internal/report"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "request a short overview of the pings today",
	Run: func(cmd *cobra.Command, args []string) {
		if err := report.Status(database); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
