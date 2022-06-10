package main

import (
	"fmt"
	"log"

	"github.com/flanksource/gitlab-multirepo-deployer/cmd"
	"github.com/spf13/cobra"
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
	root.PersistentFlags().StringP("token", "t", "", "ci trigger token")
	root.PersistentFlags().StringP("pat", "p", "", "personal access token for api calls")
	root.PersistentFlags().StringP("token-file", "f", "", "file with key-value pairs of project specific tokens")
	root.PersistentFlags().IntP("timeout", "o", 5, "timeout")

	if err := root.Execute(); err != nil {
		log.Fatalf("error running application: %s", err)
	}
}
