package main

import (
	"marvsman/dodobackend/logging"
	"marvsman/dodobackend/pbcustom"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	logging.SetupLogging(log)
	run()
}

func run() {
	app := initPocketbase()

	// add a handler to create new default userSettings whenever a user registers
	pbcustom.AddNewSettingsHandler(app, log)

	// start pocketbase
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func initPocketbase() *pocketbase.PocketBase {
	pocketbase.Version = "v0.0.1"
	app := pocketbase.New()
	migratecmd.MustRegister(app, app.RootCmd, &migratecmd.Options{
		Automigrate: true,
	})

	return app
}
