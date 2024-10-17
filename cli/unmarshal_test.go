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

func makeRepos() []*cli.Repo {
	return []*cli.Repo{
		{Remote: "backup", Dir: "/home/username/foo"},
		{Remote: "backblaze", Dir: "/Users/foo/bar"},
		{Remote: "backup", Dir: "/home/username/bar"},
		{Remote: "rsync", Dir: "foo/bar"},
		{Remote: "backup", Dir: "foo/bar/buzz"},
	}
}

// TODO: Write a test method for toml files with $HOME replaced by the user’s
// home directory.
// func makeBackupsForReplacement() []*cli.Repo {
// 	backups := make([]*cli.Repo, 0, 2)
// 	backups = append(backups, &cli.Repo{Remote: "backup"})
// 	backups = append(backups, &cli.Repo{Remote: "backblaze", Dir: ""})
// 	return backups
// }

func TestUnmarshalSuccess(t *testing.T) {
	t.Parallel()
	expected := makeRepos()
	app := &cli.App{}
	actual := app.Unmarshal("testdata/backups.toml", false)
	if app.ExitValue != exitSuccess {
		t.Fatal("app.ExitValue != exitSuccess")
	}
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("app.Unmarshal() failure (-want +got)\n%s", diff)
	}
}

func TestUnmarshalFailure(t *testing.T) {
	t.Parallel()
	app := &cli.App{}
	app.Unmarshal("testdata/nope.toml", false)
	if app.ExitValue != exitFailure {
		t.Errorf("expected exit status: %d; actual exit status: %d", exitFailure, app.ExitValue)
	}
}

func TestUnmarshalReplace(t *testing.T) {
	t.Parallel()
	expected := makeRepos()
	app := &cli.App{}
	actual := app.Unmarshal("testdata/backups.toml", false)
	if app.ExitValue != exitSuccess {
		t.Fatal("app.ExitValue != exitSuccess")
	}
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("app.Unmarshal() failure (-want +got)\n%s", diff)
	}
}
