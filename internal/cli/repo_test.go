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

func fakeAppEnv(config string) *appEnv {
	return &appEnv{
		cmd:     "test",
		subCmd:  "testing",
		config:  config,
		exitVal: exitSuccess,
	}
}

func TestGetReposSuccess(t *testing.T) {
	expected := makeRepos()
	config := "testdata/backups.json"
	app := fakeAppEnv(config)
	actual := app.repos()

	if app.exitVal != exitSuccess {
		t.Fatal("test cannot finish since app.repos() failed")
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("app.repos(%q) failure (-want +got)\n%s", config, diff)
	}
}

func TestGetReposFailure(t *testing.T) {
	app := fakeAppEnv("testdata/nope.json")
	app.repos()

	if app.exitVal != exitFailure {
		t.Error("app.exitVal expected exitFailure; actual exitSuccess")
	}
}

func TestRepoChecks(t *testing.T) {
	config := "testdata/repo-checks.json"
	app := fakeAppEnv(config)
	actual := app.repos()

	if app.exitVal != exitSuccess {
		t.Fatal("test cannot finish since app.repos() failed")
	}

	if len(actual) != 0 {
		t.Errorf("app.repos(%q) expected len(repos) = 0; actual: %d", config, len(actual))
	}
}
