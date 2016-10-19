package main

import (
	"errors"
	"fmt"

	circleci "github.com/jszwedko/go-circleci"
)

var (
	ErrNoBuilds = errors.New("no recent, successful builds for project")
)

type Project struct {
	*circleci.Client
	Account    string
	Repository string
	Branch     string
	Artifact   string
}

func (p *Project) Validate() error {
	if p.Token == "" {
		return fmt.Errorf("API token required")
	}
	if p.Account == "" {
		return fmt.Errorf("account required")
	}
	if p.Repository == "" {
		return fmt.Errorf("repository required")
	}
	return nil
}

func (p *Project) LatestBuild() (*Build, error) {
	builds, err := p.ListRecentBuildsForProject(
		p.Account, p.Repository, p.Branch, "successful", 1, 0,
	)
	if err != nil {
		return nil, err
	}

	if len(builds) == 0 {
		return nil, ErrNoBuilds
	}
	return &Build{builds[0], p}, err
}

func (p *Project) artifacts(buildNum int) ([]*circleci.Artifact, error) {
	return p.ListBuildArtifacts(
		p.Account, p.Repository, buildNum,
	)
}
