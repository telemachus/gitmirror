# TODO

+ Switch configuration from json to toml.
+ Currently, I add `GIT_TERMINAL_PROMPT=0` to the environment of the git command
  that I call with `os.Exec`.  This works for me, but others may want git to
  prompt them for their username and password.  I should make this configurable,
  maybe in the configuration file, maybe on the command line, and maybe both.
  (This shows up in clone.go and in update.go.)
+ The clone.go and update.go files share a lot of structure and code.  Maybe
  I can DRY up these two?
+ Investigate linting more.  I probably don't need all the linting options
  I have now, and I should simplify the Makefile.  I should also study the
  options for each linter more so that I make sure to use them as well as
  possible.
