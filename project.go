// Copyright 2016-2017 HashRabbit, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"fmt"

	circleci "github.com/jszwedko/go-circleci"
)

var (
	ErrNoBuilds = errors.New("no recent, successful builds")
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
