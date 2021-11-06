package cli

// TODO: add other constants or data
const (
	appName    = "git-backup"
	appVersion = "v0.0.1"
	appUsage   = `usage: git-backup [-config FILENAME]

options:
    -config FILENAME    Specify a configuration file (default = ~/.gitmirror.toml)
    -help               Show this message
    -version            Show version`
	defaultConfig = ".gitmirror.toml"
	exitSuccess   = 0
	exitFailure   = 1
)
