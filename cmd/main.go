package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/woojiahao/go-http-server/internal/server"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	var port int
	var path, config, serverName string
	currentDir, _ := os.Getwd()

	app := &cli.App{
		Name:  "go-http",
		Usage: "starts a simple HTTP server",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "port",
				Value:       8000,
				Usage:       "port of the HTTP server",
				Destination: &port,
			},
			// TODO Ensure that this is an absolute path
			&cli.StringFlag{
				Name:        "path",
				Value:       currentDir,
				Usage:       "root path of HTTP server to serve documents from",
				Destination: &path,
			},
			&cli.StringFlag{
				Name:        "config",
				Value:       "",
				Usage:       "configuration filename (in YAML)",
				Destination: &config,
			},
			&cli.StringFlag{
				Name:        "server_name",
				Value:       "HTTP server",
				Usage:       "server name",
				Destination: &serverName,
			},
		},
		Action: func(c *cli.Context) error {
			if strings.TrimSpace(config) != "" {
				var c *server.Config
				fmt.Println("Using configuration")
				yamlFile, err := ioutil.ReadFile(config)
				if err != nil {
					log.Fatal(err)
				}
				err = yaml.Unmarshal(yamlFile, &c)
				if err != nil {
					log.Fatal(err)
				}

				// TODO Handle if blank, use default value instead
				port = c.Port
				path = c.Path
				serverName = c.ServerName
			}

			s := server.Create(port, path, serverName)
			s.Start()

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
