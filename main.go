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
				Name:  "hide",
				Usage: "Hide all file names",
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:        "path",
						DefaultText: ".",
					},
				},
				Action: hide,
			},
			{
				Name:  "show",
				Usage: "Show all file names",
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:        "path",
						DefaultText: ".",
					},
				},
				Action: show,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
