package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	circleci "github.com/jszwedko/go-circleci"
)

const (
	PollFrequency = 1 * time.Minute
)

var SameBuild = errors.New("same build")

var project = &Project{
	Client: &circleci.Client{},
}

func usage(w io.Writer) {
	fmt.Fprintf(w, `Usage: %s <options>

Environment Variables:
  CIRCLECI_TOKEN      - API token to use
  CIRCLECI_ACCOUNT    - account that project belongs to
  CIRCLECI_REPOSITORY - repository builds will be requested from
  CIRCLECI_BRANCH     - limit builds to a specific branch
  CIRCLECI_ARTIFACT   - path of artifact to run (defaults to first build artifact)
  CIRCLECI_DEBUG      - enable debug log events
`, os.Args[0])
}

func init() {
	project.Token = os.Getenv("CIRCLECI_TOKEN")
	project.Account = os.Getenv("CIRCLECI_ACCOUNT")
	project.Repository = os.Getenv("CIRCLECI_REPOSITORY")
	project.Branch = os.Getenv("CIRCLECI_BRANCH")
	project.Artifact = os.Getenv("CIRCLECI_ARTIFACT")
}

func main() {
	log.SetHandler(text.New(os.Stderr))

	if os.Getenv("CIRCLECI_DEBUG") != "" {
		log.SetLevel(log.DebugLevel)
	}

	if err := project.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		usage(os.Stderr)
		os.Exit(1)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, os.Kill)

	artifacts := make(chan *Artifact, 1)

	go poll(artifacts)
	run(done, artifacts)
}

func run(done <-chan os.Signal, artifacts <-chan *Artifact) {
	var curr *Artifact
	for {
		select {
		case <-done:
			if curr != nil {
				curr.Stop()
			}
			return
		case next := <-artifacts:
			if curr != nil {
				curr.Stop()
			}
			next.Start(os.Args[1:])
			curr = next
		}
	}
}

func pollOne(curr *Build) (*Build, *Artifact, error) {
	var err error
	defer log.Trace("checking for new builds").Stop(&err)

	next, err := project.LatestBuild()
	if err != nil {
		return nil, nil, err
	}
	if next.Equal(curr) {
		return next, nil, SameBuild
	}

	artifact, err := next.Artifact()
	return next, artifact, err
}

func poll(artifacts chan<- *Artifact) {
	log.Infof("polling every %v", PollFrequency)

	curr, artifact, err := pollOne(nil)
	if err == nil {
		artifacts <- artifact
	}

	for _ = range time.Tick(PollFrequency) {
		next, artifact, err := pollOne(curr)
		if err == nil {
			curr = next
			artifacts <- artifact
		}
	}
}