package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/apex/log"
	circleci "github.com/jszwedko/go-circleci"
)

const (
	StopTimeout = 5 * time.Second
)

type Artifact struct {
	*circleci.Artifact

	Build *circleci.Build

	project *Project
	cmd     *exec.Cmd
	file    string
}

func (a *Artifact) download() error {
	var err error
	defer log.Trace("downloading artifact").Stop(&err)

	f, err := ioutil.TempFile("", "circleci-artifact")
	if err != nil {
		return err
	}
	defer f.Close()

	a.file = f.Name()

	if err = f.Chmod(0755); err != nil {
		defer a.clean()
		return err
	}

	res, err := http.Get(a.URL + "?circle-token=" + project.Token)
	defer res.Body.Close()

	if _, err = io.Copy(f, res.Body); err != nil {
		defer a.clean()
		return err
	}
	return err
}

func (a *Artifact) clean() {
	os.Remove(a.file)
}

func (a *Artifact) Start(args []string) error {
	var err error
	defer log.WithFields(log.Fields{
		"args": args,
	}).Trace("starting build artifact").Stop(&err)

	if err = a.download(); err != nil {
		return err
	}

	a.cmd = exec.Command(a.file, args...)
	a.cmd.Stdout = os.Stdout
	a.cmd.Stderr = os.Stderr

	err = a.cmd.Start()
	return err
}

func (a *Artifact) Stop() error {
	var err error
	defer log.Trace("stopping build artifact").Stop(&err)

	defer a.clean()

	if a.cmd.Process == nil {
		return nil
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- a.cmd.Wait()
	}()

	select {
	case err = <-errChan:
		// noop
	case <-time.After(StopTimeout):
		err = a.cmd.Process.Kill()
	}
	return err
}
