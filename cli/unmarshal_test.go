package cli_test

import (
	"reflect"
	"testing"

	"git.sr.ht/~telemachus/gitmirror/cli"
)

const exitFailure = 1

func makeRepos() []*cli.Repo {
	backups := make([]*cli.Repo, 0, 5)
	backups = append(backups, &cli.Repo{Remote: "backup", Dir: "/home/username/foo"})
	backups = append(backups, &cli.Repo{Remote: "backblaze", Dir: "/Users/foo/bar"})
	backups = append(backups, &cli.Repo{Remote: "backup", Dir: "/home/username/bar"})
	backups = append(backups, &cli.Repo{Remote: "rsync", Dir: "foo/bar"})
	backups = append(backups, &cli.Repo{Remote: "backup", Dir: "foo/bar/buzz"})
	return backups
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
	expected := makeRepos()
	app := &cli.App{}
	actual := app.Unmarshal("testdata/backups.toml", false)
	if app.Err != nil {
		t.Fatal(app.Err)
	}
	if !reflect.DeepEqual(expected, actual.Repos) {
		t.Errorf("expected %#v; actual %#v", expected, actual.Repos)
	}
}

func TestUnmarshalFailure(t *testing.T) {
	app := &cli.App{}
	app.Unmarshal("testdata/nope.toml", false)
	if app.Err == nil {
		t.Errorf("error is nil but we tried to unmarshal a nonexistent file")
	}
	if app.ExitValue != exitFailure {
		t.Errorf("expected %d; actual %v", exitFailure, app.ExitValue)
	}
}

func TestUnmarshalReplace(t *testing.T) {
	expected := makeRepos()
	app := &cli.App{}
	actual := app.Unmarshal("testdata/backups.toml", false)
	if app.Err != nil {
		t.Fatal(app.Err)
	}
	if !reflect.DeepEqual(expected, actual.Repos) {
		t.Errorf("expected %#v; actual %#v", expected, actual.Repos)
	}
}
