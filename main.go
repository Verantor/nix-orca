package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

func output(build string, push string) {
	fmt.Println("==================")
	fmt.Printf("  Build: %v --> %v  \n", build, push)
	fmt.Println("==================")
}
func update() error {
	cmd := exec.Command("nix", "flake", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = "/home/ver/.dotfiles/"
	if err := cmd.Run(); err != nil {
		fmt.Println("failed to update flake", err)
		return err
	} else {
		return nil
	}
}
func build() error {

	cmd := exec.Command("nixos-rebuild", "build", "--flake", ".#main", "--log-format", "internal-json", "-v")
	cmdNom := exec.Command("nom", "--json")
	pipeReader, pipeWriter := io.Pipe()
	cmd.Stdout = pipeWriter
	cmd.Stderr = pipeWriter

	cmdNom.Stdin = pipeReader
	cmd.Dir = "/home/ver/.dotfiles/"
	cmdNom.Stdout = os.Stdout
	cmdNom.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("nixos-rebuild failed: %v", err)

		return err
	}
	if err := cmdNom.Start(); err != nil {
		log.Fatalf("nom failed: %v", err)
		return err
	}
	go func() {
		cmd.Wait()
		pipeWriter.Close()
	}()
	if err := cmdNom.Wait(); err != nil {
		log.Fatalf("nom failed to finish: %v", err)
		return err
	}
	return nil

}
func buildToRemote() error {

	cmd := exec.Command("nixos-rebuild", "build", "--flake", ".#main", "--log-format", "internal-json", "-v", "--target-host", "ver@192.168.178.190", "--use-remote-sudo")
	cmdNom := exec.Command("nom", "--json")
	pipeReader, pipeWriter := io.Pipe()
	cmd.Stdout = pipeWriter
	cmd.Stderr = pipeWriter

	cmdNom.Stdin = pipeReader
	cmd.Dir = "/home/ver/.dotfiles/"
	cmdNom.Stdout = os.Stdout
	cmdNom.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("nixos-rebuild failed: %v", err)

		return err
	}
	if err := cmdNom.Start(); err != nil {
		log.Fatalf("nom failed: %v", err)
		return err
	}
	go func() {
		cmd.Wait()
		pipeWriter.Close()
	}()
	if err := cmdNom.Wait(); err != nil {
		log.Fatalf("nom failed to finish: %v", err)
		return err
	}
	return nil

}
func diff() error {
	cmd := exec.Command("nvd", "diff", "/run/current-system", "result")
	cmd.Dir = "/home/ver/.dotfiles/"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("nvd failed ", err)
		return err
	} else {
		return nil
	}
}
func activate() error {
	cmd := exec.Command("sudo", "./result/activate")
	cmd.Dir = "/home/ver/.dotfiles/"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("failed to activate system ", err)
		return err
	} else {
		return nil
	}

}
func remove() error {
	cmd := exec.Command("rm", "./result")
	cmd.Dir = "/home/ver/.dotfiles/"
	if err := cmd.Run(); err != nil {
		fmt.Println("failed to remove result ", err)
		return err
	} else {
		return nil
	}
}
func git() error {
	cmd := exec.Command("git", "commit", "-am", "'NixOS Rebuilt'")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = "/home/ver/.dotfiles/"
	if err := cmd.Run(); err != nil {
		fmt.Println("!!! CANT COMMIT GIT CHANGES !!!")
		return err
	}
	return nil
}
func backupFlakeLock() error {
	if askForConfirmation("Do you want to backup flake.lock") {
		cmd := exec.Command("cp", "flake.lock", "flakeBackup.lock")
		cmd.Dir = "/home/ver/.dotfiles/"
		if err := cmd.Run(); err != nil {
			fmt.Println("failed to backup flake.lock", err)
			return err
		}
		return nil
	} else {
		return nil
	}
}

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
func main() {
	app := &cli.App{
		EnableBashCompletion: true,
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
							err := build()
							if err != nil {
								return err
							}
							err = diff()
							if err != nil {
								return err
							}
							err = activate()
							if err != nil {
								return err
							}
							err = remove()
							if err != nil {
								return err
							}
							err = git()
							if err != nil {
								fmt.Println("!!! CANT COMMIT GIT CHANGES !!!")
							}

							return nil
						},
					},
					{
						Name:    "remote",
						Usage:   "build on remote machine",
						Aliases: []string{"r"}, //add to other functions
						Action: func(cCtx *cli.Context) error {

							output("l", "r")
							err := buildToRemote()
							if err != nil {
								return err
							}
							err = git()
							if err != nil {
								fmt.Println("!!! CANT COMMIT GIT CHANGES !!!")
							}

							return nil
						},
					},
				},
			},

			{
				Name:    "update",
				Aliases: []string{"u"}, //add to other functions
				Usage:   "update config",
				Subcommands: []*cli.Command{
					{
						Name:    "local",
						Aliases: []string{"l"}, //add to other functions
						Usage:   "build on local machine",
						Action: func(cCtx *cli.Context) error {
							err := update()
							if err != nil {
								return err
							}
							err = build()
							if err != nil {
								return err
							}
							err = diff()
							if err != nil {
								return err
							}
							err = activate()
							if err != nil {
								return err
							}
							err = remove()
							if err != nil {
								return err
							}
							err = git()
							if err != nil {
								fmt.Println("!!! CANT COMMIT GIT CHANGES")
							}
							err = backupFlakeLock()
							if err != nil {
								fmt.Println("!!! CANT BACKUP flake.lock FILE")
							}
							return nil
						},
					},
					{
						Name:    "remote",
						Usage:   "build on remote machine",
						Aliases: []string{"r"}, //add to other functions
						Action: func(cCtx *cli.Context) error {
							return nil
						},
					},
				},
			},

			{
				Name:    "add",
				Aliases: []string{"a"}, //add to other functions
				Usage:   "add package to config",
				Subcommands: []*cli.Command{
					{
						Name:    "homeManager",
						Aliases: []string{"h"}, //add to other functions
						Usage:   "build on local machine",
						Action: func(cCtx *cli.Context) error {

							return nil
						},
					},
					{
						Name:    "system",
						Usage:   "build on remote machine",
						Aliases: []string{"s"}, //add to other functions
						Action: func(cCtx *cli.Context) error {

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
