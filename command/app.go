package command

import (
	"github.com/kkoch986/ai-skeletons-output-generation/command/internal/server"
	cli "github.com/urfave/cli/v2"
)

// App presents the command line interface to start the different
// processes.
func App() *cli.App {
	app := cli.NewApp()
	app.Name = "output-generation"
	app.Description = "a service for converting the response generation into sequences of physical actions"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "log-level",
			EnvVars: []string{"LOG_LEVEL"},
			Usage:   "Set the log level",
			Value:   "WARN",
		},
	}
	app.Commands = []*cli.Command{
		server.Command,
	}
	return app
}
