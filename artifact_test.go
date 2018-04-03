package main_test

import (
	"testing"

	. "github.com/hashrabbit/circleci-runner"
	circleci "github.com/jszwedko/go-circleci"
)

func TestArtifactHasInformativeName(t *testing.T) {
	a := &Artifact{
		Artifact: &circleci.Artifact{
			Path: "foo/bar/baz/main",
		},
		Build: &circleci.Build{
			BuildNum:    123,
			Username:    "acme",
			Reponame:    "widgets",
			VcsRevision: "abc123487915bf84c7c67a653c55aad148dfa508",
		},
	}

	expected := "acme-widgets-123-gabc1234_main"
	actual := a.Name()
	if actual != expected {
		t.Errorf("expected artifact name to be %q, got %q", expected, actual)
	}
}
