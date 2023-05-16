package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "hide",
				Usage:  "Hide all file names",
				Action: hide,
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:        "path",
						DefaultText: ".",
					},
				},
			},
			{
				Name:   "show",
				Usage:  "Show all file names",
				Action: show,
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:        "path",
						DefaultText: ".",
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
