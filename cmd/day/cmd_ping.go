package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/roqcode/day/internal/ping"
	"github.com/roqcode/day/internal/scope"
	"github.com/spf13/cobra"
)

var (
	ago      int8
	at       time.Time
	silent   bool
	scopeArg string
	source   string
)

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "ping an activity",
	Long:  "ping an activity. the activity is recorded in a db.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("ago") && (ago < 1 || ago > 60) {
			fmt.Print("'--ago' only accepts values between 1 and 60. use --at for pings with a specific time\n\n")
			cmd.Help()
			os.Exit(1)
		}

		if cmd.Flags().Changed("ago") && cmd.Flags().Changed("at") {
			fmt.Print("'--at' and '--ago' can not be used together\n\n")
			cmd.Help()
			os.Exit(1)
		}

		if err := ping.Ping(database, args[0], ago, at, silent, scopeArg, source, c.Scopes.Predefined); err != nil {

			if errors.Is(err, scope.ErrAborted) {
				fmt.Println("canceled")
				os.Exit(0)
			}

			log.Fatal(err)
		}
	},
}

func init() {
	pingCmd.Flags().StringVar(&scopeArg, "scope", "", "the scope of the activity e.g. ticket name, scrum event, meeting")
	pingCmd.Flags().StringVar(&source, "source", "manual", "the source of the ping. used in automated pings")
	pingCmd.Flags().BoolVarP(&silent, "silent", "s", false, "suppress output. used in automated pings")
	pingCmd.Flags().Int8Var(&ago, "ago", 0, "used to make a retroactive ping in minutes.  e.g. 'day ping \"internal meeting\" --ago 30' for a ping 30 minutes ago. ago should only be used for values between 1 and 60. for pings further back use '--at' with a specific time")
	pingCmd.Flags().TimeVar(&at, "at", time.Time{}, []string{"15:04"}, "schedule ping at given time (e.g. 14:30)")
	rootCmd.AddCommand(pingCmd)
}
