package cli

import (
	"os"

	"github.com/ch3lo/inspector/logger"
	"github.com/ch3lo/inspector/version"
	"github.com/codegangsta/cli"
)

func globalFlags() []cli.Flag {
	flags := []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Debug de la app",
		},
		cli.StringFlag{
			Name:   "listen",
			Value:  ":8080",
			Usage:  "Direccion y puerto donde escuchara el servidor",
			EnvVar: "INSPECTOR_LISTEN",
		},
		cli.StringFlag{
			Name:   "log-level",
			Value:  "info",
			Usage:  "Nivel de verbosidad de log",
			EnvVar: "INSPECTOR_LOG_LEVEL",
		},
		cli.StringFlag{
			Name:   "log-formatter",
			Value:  "text",
			Usage:  "Formato de log",
			EnvVar: "INSPECTOR_LOG_FORMATTER",
		},
		cli.BoolFlag{
			Name:   "log-colored",
			Usage:  "Coloreo de log :D",
			EnvVar: "INSPECTOR_LOG_COLORED",
		},
		cli.StringFlag{
			Name:   "log-output",
			Value:  "console",
			Usage:  "Output de los logs. console | file",
			EnvVar: "INSPECTOR_LOG_OUTPUT",
		},
	}

	return flags
}

func setupApplication(c *cli.Context) error {
	logConfig := logger.Config{
		Level:     c.String("log-level"),
		Formatter: c.String("log-formatter"),
		Colored:   c.Bool("log-colored"),
		Output:    c.String("log-output"),
		Debug:     c.Bool("debug"),
	}

	return logger.Configure(logConfig)
}

// RunApp punto de entrada para la aplicacion
// Procesa todos los comandos y argumentos de la APP
func RunApp() {
	app := cli.NewApp()
	app.Name = "inspector"
	app.Usage = "API para visualizar informacion de contenedores."
	app.Version = version.VERSION + " (" + version.GITCOMMIT + ")"

	app.Flags = globalFlags()

	app.Before = func(c *cli.Context) error {
		return setupApplication(c)
	}

	app.Commands = commands

	if err := app.Run(os.Args); err != nil {
		logger.Instance().Fatalln(err)
	}
}
