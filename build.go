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

	circleci "github.com/jszwedko/go-circleci"
)

var (
	ErrNoArtifacts       = errors.New("no artifacts for build")
	ErrNoMatchedArtifact = errors.New("no matching artifact for build")
)

type Build struct {
	*circleci.Build

	project *Project
}

func (b *Build) Equal(other *Build) bool {
	return other != nil && b.BuildNum == other.BuildNum
}

func (b *Build) Artifact() (*Artifact, error) {
	artifacts, err := b.project.artifacts(b.BuildNum)
	if err != nil {
		return nil, err
	}
	if len(artifacts) == 0 {
		return nil, ErrNoArtifacts
	}

	if b.project.Artifact == "" {
		return &Artifact{
			Artifact: artifacts[0],
			Build:    b.Build,
			project:  project,
		}, nil
	}

	for _, a := range artifacts {
		if a.Path == b.project.Artifact {
			return &Artifact{
				Artifact: a,
				Build:    b.Build,
			}, nil
		}
	}
	return nil, ErrNoMatchedArtifact
}
