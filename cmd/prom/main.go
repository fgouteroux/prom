package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"

	"github.com/fgouteroux/prom/internal/version"
)

func BuildApp() *cli.App {
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}- {{ . }}
   {{end}}{{end}}{{if .Commands}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}{{if .Version}}
VERSION:
   {{.Version}}
   {{end}}
`
	app := cli.NewApp()
	app.Name = "prom"
	app.Usage = "Prometheus tooling"
	app.Copyright = "(c) 2023 François Gouteroux"
	app.Version = version.Version
	app.UseShortOptionHandling = true
	app.Authors = []*cli.Author{
		{
			Name:  "François Gouteroux",
			Email: "francois.gouteroux@gmail.com",
		},
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"D"},
			Usage:   "show debug output",
			EnvVars: []string{"PROM_DEBUG", "DEBUG"},
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	}
	app.Suggest = true

	app.Commands = []*cli.Command{
		metricsCommand(),
	}

	return app
}

func main() {
	app := BuildApp()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Print(err.Error())
	}
}
