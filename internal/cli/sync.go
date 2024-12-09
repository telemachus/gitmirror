package cli

func subCmdSync(app *appEnv) {
	rs := app.repos()
	app.update(rs)
	app.clone(rs)
}
