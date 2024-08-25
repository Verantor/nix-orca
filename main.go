package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

func output(build string, push string) {
	fmt.Println("==================")
	fmt.Printf("  Build: %v --> %v  \n", build, push)
	fmt.Println("==================")
}

func build() {

	cmd := exec.Command("nixos-rebuild", "build", "--flake", ".#main")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("rebuild failed: ", err)
	} else {
		cmd := exec.Command("nvd", "diff", "/run/current-system", "result")
		if err := cmd.Run(); err != nil {
			fmt.Println("nvd failed ", err)
		} else {

			cmd := exec.Command("sudo", "./result/activate")
			if err := cmd.Run(); err != nil {
				fmt.Println("failed to activate system ", err)
			}
		}
	}
}

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
						Action: func(cCtx *cli.Context) error {
							output("l", "l")
							build()
							return nil
						},
					},
					{
						Name:    "remote",
						Usage:   "build on remote machine",
						Aliases: []string{"r"}, //add to other functions
						Action: func(cCtx *cli.Context) error {
							build()
							return nil
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
