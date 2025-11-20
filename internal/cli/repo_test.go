package cli

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func makeRepos() []Repo {
	return []Repo{
		{URL: "https://github.com/foo/foo.git", Name: "foo.git"},
		{URL: "https://github.com/bar/bar.git", Name: "bar.git"},
		{URL: "https://example.com/buzz/fizz.git", Name: "random.git"},
	}
}

func fakeCmdEnv(confFile string) *cmdEnv {
	return &cmdEnv{
		name:     "test",
		confFile: confFile,
	}
}

func TestGetReposSuccess(t *testing.T) {
	t.Parallel()
	expected := makeRepos()
	confFile := "testdata/backups.json"
	cmd := fakeCmdEnv(confFile)
	actual, err := cmd.repos()
	if err != nil {
		t.Fatalf("test cannot finish since cmd.repos() failed: %v", err)
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("cmd.repos(%q) failure (-want +got)\n%s", confFile, diff)
	}
}

func TestGetReposFailure(t *testing.T) {
	t.Parallel()
	cmd := fakeCmdEnv("testdata/nope.json")
	_, err := cmd.repos()

	if err == nil {
		t.Error("expected error from cmd.repos(); got nil")
	}
}

func TestRepoChecks(t *testing.T) {
	t.Parallel()
	confFile := "testdata/repo-checks.json"
	cmd := fakeCmdEnv(confFile)
	actual, err := cmd.repos()
	if err != nil {
		t.Fatalf("test cannot finish since cmd.repos() failed: %v", err)
	}

	if len(actual) != 0 {
		t.Errorf("cmd.repos(%q) expected len(repos) = 0; actual: %d", confFile, len(actual))
	}
}
