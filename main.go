package main

import (
	"log"
	"os"

	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "Cherry"
	app.Usage = "Git utility"
	app.EnableBashCompletion = true
	app.HideHelp = true
	app.HideVersion = true
	app.Commands = []cli.Command{
		{
			Name:  "diff",
			Usage: "diff source_branch target_branch",
			Action: func(c *cli.Context) error {
				diff(c)
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "author",
				},
				cli.StringFlag{
					Name:  "issue",
					Usage: "\"1000,1001\"",
				},
			},
		},
		{
			Name:  "pick",
			Usage: "pick source_branch target_branch",
			Action: func(c *cli.Context) error {
				pick(c)
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "author",
				},
				cli.StringFlag{
					Name:  "issue",
					Usage: "\"1000,1001\"",
				},
				cli.BoolFlag{
					Name: "step",
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
