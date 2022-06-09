package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/flanksource/gitlab-multirepo-deployer/pkg"
	"github.com/spf13/cobra"
)

var Scan = &cobra.Command{
	Use:   "scan",
	Short: "Scan all projects for specified branch",
	Run: func(cmd *cobra.Command, args []string) {
		branch, _ := cmd.Flags().GetString("branch")
		configFile, _ := cmd.Flags().GetString("config")

		accessToken, _ := cmd.Flags().GetString("pat")
		if accessToken == "" {
			accessToken = os.Getenv("PERSONAL_ACCESS_TOKEN")
			if accessToken == "" {
				accessToken, _ = cmd.Flags().GetString("token")
				if accessToken == "" {
					accessToken = os.Getenv("CI_JOB_TOKEN")
					if accessToken == "" {
						log.Fatalf("No access token provided")
					}
				}
			}
		}
		deployToken, _ := cmd.Flags().GetString("token")

		cfg, err := pkg.NewConfig(configFile, accessToken, deployToken)
		if err != nil {
			log.Fatalf("Could not create config: %v", err)
		}

		for i := range cfg.Projects {
			project := &cfg.Projects[i]
			if err = project.SetBranch(branch); err != nil {
				if !strings.Contains(err.Error(), "{message: 404 Branch Not Found}") {
					log.Fatalf("unexpected error looking up branch: %s", err.Error())
				}
			}
		}
		fmt.Printf("Branch %s present in the following projects:\n", branch)

		count := 0
		for i := range cfg.Projects {
			project := &cfg.Projects[i]
			if branch, _ = project.GetBranch(); branch != "" {
				fmt.Println(project.Name)
				count++
			}
		}
		fmt.Printf("Total: %d\n", count)
	},
}
