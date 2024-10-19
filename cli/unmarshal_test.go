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
		{Remote: "backup", Dir: "/home/username/foo"},
		{Remote: "backblaze", Dir: "/Users/foo/bar"},
		{Remote: "backup", Dir: "/home/username/bar"},
		{Remote: "rsync", Dir: "foo/bar"},
		{Remote: "backup", Dir: "foo/bar/buzz"},
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
		t.Errorf("app.Unmarshal() failure (-want +got)\n%s", diff)
	}
}

func TestUnmarshalFailure(t *testing.T) {
	app := cli.NewApp()
	app.Unmarshal("testdata/nope.json", false)
	if app.ExitValue != exitFailure {
		t.Errorf("expected exit status: %d; actual exit status: %d", exitFailure, app.ExitValue)
	}
}

func TestUnmarshalRepoChecks(t *testing.T) {
	app := cli.NewApp()
	actual := app.Unmarshal("testdata/repo-checks.json", false)
	if len(actual) != 0 {
		t.Error("expected no repos from testdata/repo-checks.json")
	}
}
