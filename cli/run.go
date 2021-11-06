package cli

func Run(args []string) int {
	app := &App{ExitValue: exitSuccess}
	configFile, isDefault := app.ParseFlags(args)
	wanted := app.Unmarshal(configFile, isDefault)
	app.MirrorRepos(wanted)
	app.DisplayInfo()
	app.DisplayError()
	return app.ExitValue
}
