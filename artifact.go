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
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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

func (a *Artifact) Name() string {
	return fmt.Sprintf("%s-%s-%d-g%s_%s", a.Build.Username, a.Build.Reponame,
		a.Build.BuildNum, a.Build.VcsRevision[0:7], filepath.Base(a.Path))
}

func (a *Artifact) download() error {
	f, err := ioutil.TempFile("", a.Name())
	if err != nil {
		return err
	}
	defer f.Close()
	a.file = f.Name()

	if err = f.Chmod(0755); err != nil {
		defer a.clean()
		return err
	}

	defer log.WithField("path", a.file).Trace("downloading artifact").Stop(&err)

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
