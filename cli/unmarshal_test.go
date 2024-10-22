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

func makeRepos() []cli.Repo {
	return []cli.Repo{
		{URL: "https://github.com/foo/foo.git", Name: "foo.git"},
		{URL: "https://github.com/bar/bar.git", Name: "bar.git"},
		{URL: "https://example.com/buzz/fizz.git", Name: "random.git"},
	}
}

func TestUnmarshalSuccess(t *testing.T) {
	expected := makeRepos()
	app := cli.NewApp()
	actual := app.Unmarshal("testdata/backups.json", false)
	if app.ExitValue != exitSuccess {
		t.Fatal("app.ExitValue != exitSuccess")
	}
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("app.Unmarshal(\"testdata/backups.json\") failure (-want +got)\n%s", diff)
	}
}

func TestUnmarshalFailure(t *testing.T) {
	app := cli.NewApp()
	app.Unmarshal("testdata/nope.json", false)
	if app.ExitValue != exitFailure {
		t.Errorf("app.Unmarshal(\"testdata/nope.json\") exit value: %d; expected %d", app.ExitValue, exitFailure)
	}
}

func TestUnmarshalRepoChecks(t *testing.T) {
	app := cli.NewApp()
	actual := app.Unmarshal("testdata/repo-checks.json", false)
	if len(actual) != 0 {
		t.Errorf("app.Unmarshal(\"testdata/repo-checks.json\") expected len(repos) = 0; actual: %d", len(actual))
	}
}
