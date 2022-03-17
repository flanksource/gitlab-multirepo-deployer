package cmd

import (
	"fmt"
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
		timeOut, _ := cmd.Flags().GetInt("timeout")

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
			errName := project.Name
			if errName == "" {
				errName = string(project.ID)
			}
			for _, lookup := range searchBranches {
				err = project.SetBranch(lookup)
				if err != nil {
					if !strings.Contains(err.Error(), "{message: 404 Branch Not Found}") {
						log.Fatalf("unexpected error looking up branch in porject %s: %s", errName, err.Error())
					}
				} else {
					foundBranch = true
					break
				}
			}
			if !foundBranch {
				log.Fatalf("Could not find deployable branch in project %s", errName)
			}
			err = project.TriggerPipeline(cfg.Variables)
			if err != nil {
				log.Fatalf("could not trigger pipeline in project %s: %v", errName, err)
			}

		}
		startTime := time.Now()
		parsedTimeout, err := time.ParseDuration(fmt.Sprintf("%sm", timeOut))
		if err != nil {
			log.Fatalf("Could not parse timeout duration: %s", err)
		}
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
				errName := project.Name
				if errName == "" {
					errName = string(project.ID)
				}
				status, err := project.GetPipeLineStatus()
				if err != nil {
					log.Printf("Error retrieving pipeline state for project %s: %v", errName, err)
					continue
				}
				count[status]++
			}
			if count["success"] == len(cfg.Projects) {
				break
			}
			if time.Now().After(startTime.Add(parsedTimeout)) {
				log.Fatalf("Timed out waiting for deployments")
			}
			time.Sleep(30 * time.Second)
		}
	},
}
