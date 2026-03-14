package main

import (
	"log"

	"github.com/roqcode/day/internal/report"
	"github.com/spf13/cobra"
)

var heatmapCmd = &cobra.Command{
	Use:   "heatmap",
	Short: "shows a github activities style heatmap graph",
	Run: func(cmd *cobra.Command, args []string) {
		if err := report.Heatmap(database); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(heatmapCmd)
}
