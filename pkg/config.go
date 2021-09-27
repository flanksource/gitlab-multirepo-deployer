package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/xanzy/go-gitlab"
	"io/ioutil"
	yaml "gopkg.in/flanksource/yaml.v3"
	"os"
)

func NewConfig(file string, accessToken string, jobToken string) (Config, error){
	cfg := Config{}
	data, err := ioutil.ReadFile("projects.yaml")
	reader := bytes.NewReader(data)
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)

	if err := decoder.Decode(&cfg); err != nil {
		return cfg, errors.New(fmt.Sprintf("Failed to parse project file: %v", err))
	}

	git, err := gitlab.NewClient(accessToken)
	if err != nil {
		return cfg, errors.New(fmt.Sprintf("failed to create client: %v", err))
	}
	if jobToken == "" {
		jobToken = os.Getenv("CI_JOB_TOKEN")
	}
	for i := range cfg.Projects {
		cfg.Projects[i].SetClient(git)
		if err := cfg.Projects[i].SetToken(jobToken); err != nil {
			return cfg, errors.New(fmt.Sprintf("No deployment token for %s: %s", cfg.Projects[i].Name, err))
		}
		if cfg.Projects[i].GroupID == 0 {
			cfg.Projects[i].GroupID = cfg.GroupID
		}
	}
	return cfg, nil
}

type Config struct {
	GroupID int           `yaml:"groupID" json:"groupID"`
	Projects []GitLabProject `yaml:"projects" json:"projects"`
	Variables map[string]string `yaml:"variables" json:"variables"`
}