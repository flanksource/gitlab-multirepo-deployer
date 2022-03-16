package cmd

import (
	"github.com/flanksource/gitlab-multirepo-deployer/pkg"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"time"
)

var Trigger = &cobra.Command{
	Use:   "trigger",
	Short: "trigger all pipelines across names branch, falling back to main or master",
	Run: func(cmd *cobra.Command, args []string) {
		deployBranch, _ := cmd.Flags().GetString("branch")
		searchBranches := []string{deployBranch, "main", "master"}
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
			foundBranch := false
			for _, lookup := range searchBranches {
				err = project.SetBranch(lookup)
				if err != nil {
					if !strings.Contains(err.Error(), "{message: 404 Branch Not Found}") {
						log.Fatalf("unexpected error looking up branch: %s", err.Error())
					}
				} else {
					foundBranch = true
					break
				}
			}
			if !foundBranch {
				log.Fatalf("Could not find deployable branch")
			}
			err = project.TriggerPipeline(cfg.Variables)
			if err != nil {
				log.Fatalf("could not trigger pipeline: %v", err)
			}

		}
		startTime := time.Now()
		for {
			count := map[string]int{
				"created":              0,
				"waiting_for_resource": 0,
				"preparing":            0,
				"pending":              0,
				"running":              0,
				"success":              0,
				"failed":               0,
				"canceled":             0,
				"skipped":              0,
				"manual":               0,
				"scheduled":            0,
			}
			for i := range cfg.Projects {
				project := &cfg.Projects[i]
				if project == nil {
					log.Fatalf("HOW??")
				}
				status, err := project.GetPipeLineStatus()
				if err != nil {
					log.Printf("Error retrieving pipeline state: %v", err)
					continue
				}
				count[status]++
			}
			if count["success"] == len(cfg.Projects) {
				break
			}
			if time.Now().After(startTime.Add(5 * time.Minute)) {
				log.Fatalf("Timed out waiting for deployments")
			}
			time.Sleep(30 * time.Second)
		}
	},
}
