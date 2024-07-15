package pkg

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/xanzy/go-gitlab"
)

type GitLabProject struct {
	Name          string   `json:"name" yaml:"name"`
	ID            int      `json:"id" yaml:"id"`
	GroupID       int      `json:"groupID" yaml:"groupID"`
	ServiceName   string   `json:"serviceName" yaml:"serviceName"`
	ExtraServices []string `json:"extraServices" yaml:"extraServices"`
	Token         string   `json:"token", yaml:"token"`
	branch        GitLabBranch
	pipeline      *gitlab.Pipeline
	client        *gitlab.Client
}

type GitLabBranch struct {
	Name string
	Ref  string
}

func (g *GitLabProject) SetBranch(name string) error {
	id, err := g.GetID()
	if err != nil {
		return err
	}
	branch, _, err := g.client.Branches.GetBranch(id, name)
	if err != nil {
		return err
	}

	g.branch.Name = branch.Name
	g.branch.Ref = branch.Commit.ID
	return nil
}

func (g *GitLabProject) GetBranch() (string, string) {
	return g.branch.Name, g.branch.Ref
}

func (g *GitLabProject) SetPipeline(pipeline *gitlab.Pipeline) {
	g.pipeline = pipeline
}

func (g *GitLabProject) GetPipeline() *gitlab.Pipeline {
	return g.pipeline
}

func (g *GitLabProject) SetClient(client *gitlab.Client) {
	g.client = client
}

func (g *GitLabProject) SetToken(deployToken string) error {
	if g.Token != "" {
		return nil
	}
	if deployToken == "" {
		return errors.New("no deploy token provided")
	}
	g.Token = deployToken
	return nil
}

func (g *GitLabProject) TriggerPipeline(variables map[string]string) error {
	token := g.Token
	id, err := g.GetID()
	if err != nil {
		return err
	}
	if g.branch.Name == "" {
		return errors.New("no branch selected for pipeline to trigger on")
	}
	pipeline, _, err := g.client.PipelineTriggers.RunPipelineTrigger(id, &gitlab.RunPipelineTriggerOptions{
		Ref:       gitlab.String(g.branch.Name),
		Token:     gitlab.String(token),
		Variables: variables,
	})
	if err != nil {
		return err
	}
	g.SetPipeline(pipeline)
	log.Println(fmt.Sprintf("Started pipeline ID %d in project %s: %s\n", pipeline.ID, g.Name, pipeline.WebURL))

	return nil
}

func (g *GitLabProject) GetPipeLineStatus() (string, error) {
	if g.pipeline == nil {
		return "", errors.New("no pipeline has been triggered")
	}
	status, _, err := g.client.Pipelines.GetPipeline(g.pipeline.ProjectID, g.pipeline.ID)
	if err != nil {
		return "", err
	}
	return status.Status, nil
}

func (g *GitLabProject) GetID() (int, error) {
	if g.ID == 0 {
		if g.Name == "" {
			return 0, errors.New("insufficient information for project - name or id must be provided")
		}
		projSearch, _, err := g.client.Groups.ListGroupProjects(g.GroupID, &gitlab.ListGroupProjectsOptions{Search: gitlab.String(g.Name), IncludeSubGroups: true})
		if err != nil {
			return 0, errors.New(fmt.Sprintf("could not determine group ID: %v", err))
		}
		if len(projSearch) > 1 {
			var found []string
			for _, project := range projSearch {
				found = append(found, project.Name)
				if strings.Compare(g.Name, project.Name) == 0 {
					return project.ID, nil
				}
			}
			foundString := strings.Join(found, ", ")
			return 0, errors.New(fmt.Sprintf("Could not lookup project: %s, found [%s]", g.Name, foundString))
		} else if len(projSearch) == 0 {
			return 0, errors.New(fmt.Sprintf("Could not lookup project: %s, found no matches", g.Name))
		}
		g.ID = projSearch[0].ID
	}
	return g.ID, nil
}
