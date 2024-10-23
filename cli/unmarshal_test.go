package cli_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/telemachus/gitmirror/cli"
)

const (
	exitFailure = 1
	exitSuccess = 0
)

func newApp(config string) *cli.App {
	app := cli.NewApp()
	app.ConfigFile = config
	app.CustomConfig = true
	return app
}

func makeRepos() []cli.Repo {
	return []cli.Repo{
		{URL: "https://github.com/foo/foo.git", Name: "foo.git"},
		{URL: "https://github.com/bar/bar.git", Name: "bar.git"},
		{URL: "https://example.com/buzz/fizz.git", Name: "random.git"},
	}
}

func TestUnmarshalSuccess(t *testing.T) {
	expected := makeRepos()
	app := newApp("testdata/backups.json")
	actual := app.Unmarshal()
	if app.ExitValue != exitSuccess {
		t.Fatal("app.ExitValue != exitSuccess")
	}
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("app.Unmarshal() failure (-want +got)\n%s", diff)
	}
}

func TestUnmarshalFailure(t *testing.T) {
	app := newApp("testdata/nope.json")
	app.Unmarshal()
	if app.ExitValue != exitFailure {
		t.Errorf("app.Unmarshal() exit value: %d; expected %d", app.ExitValue, exitFailure)
	}
}

func TestUnmarshalRepoChecks(t *testing.T) {
	app := newApp("testdata/repo-checks.json")
	actual := app.Unmarshal()
	if len(actual) != 0 {
		t.Errorf("app.Unmarshal(\"testdata/repo-checks.json\") expected len(repos) = 0; actual: %d", len(actual))
	}
}
