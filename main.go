package main

import (
	"fmt"
	"github.com/flanksource/gitlab-multirepo-deployer/cmd"
	"github.com/spf13/cobra"
	"log"
)

var (
	version = "dev"
)

func main() {
	root := &cobra.Command{
		Use: "gitlab-multirepo-deployer",
	}
	root.AddCommand(
		cmd.Trigger,
		cmd.Scan,
	)
	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version info",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	})

	root.PersistentFlags().StringP("config", "c", "projects.yaml", "Path to config file")
	root.PersistentFlags().StringP("branch", "b", "main", "branch to trigger against")
	root.PersistentFlags().StringP("token", "t", "", "ci token")
	root.PersistentFlags().StringP("pat", "p", "", "personal access token")

	if err := root.Execute(); err != nil {
		log.Fatalf("error running application: %s", err)
	}
}
