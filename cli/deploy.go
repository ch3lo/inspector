package cli

import (
	"errors"
	"fmt"

	"github.com/ch3lo/inspector/api"
	"github.com/codegangsta/cli"
)

func deployFlags() []cli.Flag {
	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "host-ip",
			Usage: "Ip del host",
		},
		cli.StringFlag{
			Name:  "address",
			Usage: "Endpoint de Docker Engine. Formato ip:puerto",
		},
		cli.BoolFlag{
			Name:  "tlsverify",
			Usage: "Usa certificados con tlsverify",
		},
		cli.StringFlag{
			Name:  "tlscacert",
			Value: "ca.pem",
			Usage: "Ruta del archivo de configuración",
		},
		cli.StringFlag{
			Name:  "tlscert",
			Value: "cert.pem",
			Usage: "Ruta del archivo de configuración",
		},
		cli.StringFlag{
			Name:  "tlskey",
			Value: "key.pem",
			Usage: "Ruta del archivo de configuración",
		},
	}

	return flags
}

func deployBefore(c *cli.Context) error {
	if c.String("host-ip") == "" {
		return fmt.Errorf("Debe existir la ip de host")
	}

	if c.String("address") == "" {
		return fmt.Errorf("Debe existir el parametro address")
	}

	if c.Bool("tlsverify") {
		if c.String("tlscacert") == "" {
			return errors.New("Parametro tlscacert no existe")
		}

		if c.String("tlscert") == "" {
			return errors.New("Parametro tlscert no existe")
		}

		if c.String("tlskey") == "" {
			return errors.New("Parametro tlskey no existe")
		}
	}

	return nil
}

func deployCmd(c *cli.Context) {
	appConfig := api.Configuration{
		HostIP:    c.String("host-ip"),
		Address:   c.String("address"),
		TLSVerify: c.Bool("tlsverify"),
		TLSCacert: c.String("tlscacert"),
		TLSCert:   c.String("tlscert"),
		TLSKey:    c.String("tlskey"),
	}

	api.Server(appConfig)
}
