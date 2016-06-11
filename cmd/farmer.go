package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

func main() {

	app := cli.NewApp()
	app.Name = "farmer"
	app.Usage = "Flexible Application Runtime Management w/ Elastic Resources"
	app.Commands = []cli.Command{
		{
			Name:    "application",
			Aliases: []string{"app"},
			Usage:   "Deploy and manage farmer applications",
			Subcommands: []cli.Command{
				{
					Name:  "ls",
					Usage: "List your applications",
					Action: func(c *cli.Context) error {
						fmt.Println("list applications: ", c.Args().First())
						return nil
					},
				},
			},
		},
		{
			Name:    "network",
			Aliases: []string{"t"},
			Usage:   "Deploy and manage application networks",
			Subcommands: []cli.Command{
				{
					Name:  "ls",
					Usage: "List application networks",
					Action: func(c *cli.Context) error {
						fmt.Println("list app networks: ", c.Args().First())
						return nil
					},
				},
			},
		},
		{
			Name:    "server",
			Aliases: []string{"t"},
			Usage:   "Deploy and manage farmer server",
			Subcommands: []cli.Command{
				{
					Name:  "ls",
					Usage: "List farmer server",
					Action: func(c *cli.Context) error {
						fmt.Println("list app networks: ", c.Args().First())
						return nil
					},
				},
				{
					Name:   "start",
					Usage:  "Start a farm worker on this machine",
					Action: server,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "addr, a",
							Value: "0.0.0.0:0",
							Usage: "Address of circuit server.",
						},
						cli.StringFlag{
							Name:  "if",
							Value: "",
							Usage: "Bind any available port on the specified interface.",
						},
						cli.StringFlag{
							Name:  "var",
							Value: "",
							Usage: "Lock and log directory for the circuit server.",
						},
						cli.StringFlag{
							Name:  "join, j",
							Value: "",
							Usage: "Join a circuit through a current member by address.",
						},
						cli.StringFlag{
							Name:  "hmac",
							Value: "",
							Usage: "File with HMAC credentials for HMAC/RC4 transport security.",
						},
						cli.StringFlag{
							Name:  "discover",
							Value: "228.8.8.8:8822",
							Usage: "Multicast address for peer server discovery",
						},
						cli.BoolFlag{
							Name:  "docker",
							Usage: "Enable docker elements; docker command must be executable",
						},
					},
				},
			},
		},
	}
	app.Run(os.Args)
}
