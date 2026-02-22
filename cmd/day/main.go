package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	// database, err := db.InitDB()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer database.Close()

	Execute()
}

var rootCmd = &cobra.Command{
	Use:   "day",
	Short: "day is a time tracking tool",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello day")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing Zero '%s'\n", err)
		os.Exit(1)
	}
}
