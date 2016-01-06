package cli

import (
	"errors"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/ch3lo/inspector/util"
	"github.com/ch3lo/inspector/version"
	"github.com/codegangsta/cli"
)

var logFile *os.File = nil

func globalFlags() []cli.Flag {
	flags := []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Debug de la app",
		},
		cli.StringFlag{
			Name:   "log-level",
			Value:  "info",
			Usage:  "Nivel de verbosidad de log",
			EnvVar: "OVERLORD_LOG_LEVEL",
		},
		cli.StringFlag{
			Name:   "log-formatter",
			Value:  "text",
			Usage:  "Formato de log",
			EnvVar: "OVERLORD_LOG_FORMATTER",
		},
		cli.BoolFlag{
			Name:   "log-colored",
			Usage:  "Coloreo de log :D",
			EnvVar: "OVERLORD_LOG_COLORED",
		},
		cli.StringFlag{
			Name:   "log-output",
			Value:  "console",
			Usage:  "Output de los logs. console | file",
			EnvVar: "OVERLORD_LOG_OUTPUT",
		},
	}

	return flags
}

type logConfig struct {
	level     string
	Formatter string
	colored   bool
	output    string
	debug     bool
}

func setupLogger(config logConfig) error {
	var err error

	if util.Log.Level, err = log.ParseLevel(config.level); err != nil {
		return err
	}

	if config.debug {
		util.Log.Level = log.DebugLevel
	}

	switch config.Formatter {
	case "text":
		formatter := new(log.TextFormatter)
		formatter.ForceColors = config.colored
		formatter.FullTimestamp = true
		util.Log.Formatter = formatter
		break
	case "json":
		formatter := new(log.JSONFormatter)
		util.Log.Formatter = formatter
		break
	default:
		return errors.New("Formato de lo log desconocido")
	}

	switch config.output {
	case "console":
		util.Log.Out = os.Stdout
		break
	case "file":
		util.Log.Out = logFile
		break
	default:
		return errors.New("Output de logs desconocido")
	}

	return nil
}

func setupApplication(c *cli.Context) error {
	logConfig := logConfig{}
	logConfig.level = c.String("log-level")
	logConfig.Formatter = c.String("log-formatter")
	logConfig.colored = c.Bool("log-colored")
	logConfig.output = c.String("log-output")
	logConfig.debug = c.Bool("debug")

	err := setupLogger(logConfig)
	if err != nil {
		return err
	}
	return nil
}

func RunApp() {
	app := cli.NewApp()
	app.Name = "overlord"
	app.Usage = "Monitor de contenedores"
	app.Version = version.VERSION + " (" + version.GITCOMMIT + ")"

	app.Flags = globalFlags()

	app.Before = func(c *cli.Context) error {
		return setupApplication(c)
	}

	app.Commands = commands

	var err error
	logFile, err = os.OpenFile("overlord.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		util.Log.Warnln("Error al abrir el archivo")
	} else {
		defer logFile.Close()
	}

	err = app.Run(os.Args)
	if err != nil {
		util.Log.Fatalln(err)
	}
}
