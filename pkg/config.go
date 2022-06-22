package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/xanzy/go-gitlab"
	yaml "gopkg.in/flanksource/yaml.v3"
)

func NewConfig(file string, tokenFile string, accessToken string, jobToken string) (Config, error) {
	cfg := Config{}
	data, err := ioutil.ReadFile(file)
	reader := bytes.NewReader(data)
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)

	if err := decoder.Decode(&cfg); err != nil {
		return cfg, errors.New(fmt.Sprintf("Failed to parse project file: %v", err))
	}

	tokens := map[int]string{}
	data, err = ioutil.ReadFile(tokenFile)
	reader = bytes.NewReader(data)
	decoder = yaml.NewDecoder(reader)

	if err := decoder.Decode(&tokens); err != nil {
		return cfg, errors.New(fmt.Sprintf("Failed to parse token file: %v", err))
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
		if cfg.Projects[i].GroupID == 0 {
			cfg.Projects[i].GroupID = cfg.GroupID
		}
		id, err := cfg.Projects[i].GetID()
		if err != nil {
			return cfg, errors.New(fmt.Sprintf("Failed to lookup project ID: %v", err))
		}
		projectToken, ok := tokens[id]
		if !ok {
			projectToken = jobToken
		}
		if err := cfg.Projects[i].SetToken(projectToken); err != nil {
			return cfg, errors.New(fmt.Sprintf("No deployment token for %s: %s", cfg.Projects[i].Name, err))
		}
	}
	return cfg, nil
}

type Config struct {
	GroupID   int               `yaml:"groupID" json:"groupID"`
	Projects  []GitLabProject   `yaml:"projects" json:"projects"`
	Variables map[string]string `yaml:"variables" json:"variables"`
}
