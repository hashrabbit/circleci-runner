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
