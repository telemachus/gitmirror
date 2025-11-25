# gitmirror

Use git commands to clone and update backups of git repositories.

## Configuration

Gitmirror uses JSON for configuration. By default, gitmirror assumes
a configuration file at ~/.gitmirror.json. In order to use a different
configuration file, use the --config option.

```json
{
    "repos": [
        {
            "name": "whatever.git",
            "url": "https://github.com/username/whatever.git"
        },
        {
            "name": "otherstuff.git",
            "url": "git@github.com:otherusername/otherstuff.git"
        }
    ]
}
```

## Usage

Gitmirror has no subcommands. On each run, it will do the following.

1. Read the configuration file.
2. Clone any repo listed in the configuration file that is not already cloned
   using the git command `git clone --mirror`.
3. Update any repo listed in the configuration file that is already present
   using the git command `git remote update`.

Note that gitmirror does not delete repositories that are not listed in the
configuration file. If you want to delete something, you must do it manually.
Also note that gitmirror never prunes branches that have been deleted on the
remote. If you want to prune branches, you must do it manually. These two
choices are deliberate. Since gitmirror aims at preservation, only the user can
decide when to delete or prune a repository.

(c) 2024 Peter Aronoff. BSD 3-Clause license; see [LICENSE.txt][license] for
details.

[license]: /LICENSE.txt
