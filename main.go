package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{

			{
				Name:    "build",
				Aliases: []string{"b"}, //add to other functions
				Usage:   "build config",
				Subcommands: []*cli.Command{
					{
						Name:    "local",
						Aliases: []string{"l"}, //add to other functions
						Usage:   "build on local machine",
						Subcommands: []*cli.Command{
							{
								Name:    "remote",
								Aliases: []string{"r"}, //add to other functions
								Usage:   "pushes build to remote machine",
								Action: func(cCtx *cli.Context) error {
									fmt.Println("b on l and p2 r ", cCtx.Args().First())
									cmd := exec.Command("nixbr")
									cmd.Stdout = os.Stdout
									if err := cmd.Run(); err != nil {
										fmt.Println("could not run command: ", err)
									}
									return nil
								},
							},
							{
								Name:    "local",
								Usage:   "pushes build to local machine",
								Aliases: []string{"l"}, //add to other functions
								Action: func(cCtx *cli.Context) error {
									fmt.Println("b on l and p2 l ", cCtx.Args().First())
									cmd := exec.Command("nixrb")
									cmd.Stdout = os.Stdout
									if err := cmd.Run(); err != nil {
										fmt.Println("could not run command: ", err)
									}
									return nil
									return nil
								},
							},
						},
					},
					{
						Name:    "remote",
						Usage:   "build on remote machine",
						Aliases: []string{"r"}, //add to other functions
						Subcommands: []*cli.Command{
							{
								Name:    "local",
								Aliases: []string{"l"}, //add to other functions
								Usage:   "pushes build to local machine",
								Action: func(cCtx *cli.Context) error {
									fmt.Println("b on r and p2 l", cCtx.Args().First())
									return nil
								},
							},
							{
								Name:    "remote",
								Aliases: []string{"r"}, //add to other functions
								Usage:   "pushes build to remote machine",
								Action: func(cCtx *cli.Context) error {
									fmt.Println("b on r and p2 r idk if it exists", cCtx.Args().First())
									return nil
								},
							},
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
