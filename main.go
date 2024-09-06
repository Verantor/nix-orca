package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
)

var directory string = "/home/ver/.dotfiles/"
var remoteIp string = "192.168.178.190"
var remoteUser string = "ver"

func output(build string, push string) {
	fmt.Println("==================")
	fmt.Printf("  Build: %v --> %v  \n", build, push)
	fmt.Println("==================")
}
func update() error {
	cmd := exec.Command("nix", "flake", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = directory
	if err := cmd.Run(); err != nil {
		fmt.Println("failed to update flake", err)
		return err
	} else {
		return nil
	}
}
func build() error {
	output("LOCAL", "LOCAL")
	cmdBuild := exec.Command("nixos-rebuild", "build", "--flake", ".#main", "--log-format", "internal-json", "-v")
	cmdNom := exec.Command("nom", "--json")
	pipeReader, pipeWriter := io.Pipe()
	cmdBuild.Stdout = pipeWriter
	cmdBuild.Stderr = pipeWriter

	cmdNom.Stdin = pipeReader
	cmdBuild.Dir = directory
	cmdNom.Stdout = os.Stdout
	cmdNom.Stderr = os.Stderr

	if err := cmdBuild.Start(); err != nil {
		log.Fatalf("nixos-rebuild failed: %v", err)
		return err
	}
	if err := cmdNom.Start(); err != nil {
		log.Fatalf("nom failed: %v", err)
		return err
	}
	go func() {
		cmdBuild.Wait()
		pipeWriter.Close()
	}()
	if err := cmdNom.Wait(); err != nil {
		log.Fatalf("nom failed to finish: %v", err)
		return err
	}
	return nil

}
func buildToRemote() error {

	output("LOCAL", "REMOTE")
	cmd := exec.Command("nixos-rebuild", "build", "--flake", ".#main", "--log-format", "internal-json", "-v", "--target-host", remoteUser+"@"+remoteIp, "--use-remote-sudo")
	cmdNom := exec.Command("nom", "--json")
	pipeReader, pipeWriter := io.Pipe()
	cmd.Stdout = pipeWriter
	cmd.Stderr = pipeWriter

	cmdNom.Stdin = pipeReader
	cmd.Dir = directory
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

	cmd.Dir = directory
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

	cmd.Dir = directory

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

	cmd.Dir = directory
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

	cmd.Dir = directory
	if err := cmd.Run(); err != nil {
		fmt.Println("!!! CANT COMMIT GIT CHANGES !!!")
		fmt.Println("!!! CANT COMMIT GIT CHANGES !!!")
		fmt.Println("!!! CANT COMMIT GIT CHANGES !!!")

		return err
	}
	return nil
}
func backupFlakeLock() error {
	if askForConfirmation("Do you want to backup flake.lock") {
		cmd := exec.Command("cp", "flake.lock", "flakeBackup.lock")
		cmd.Dir = directory
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
func addPackage(packageName string, file string) (int, error) {
	line, err := findLineOfInsert(directory+file, "### Insert Point")
	if err != nil {
		return 0, err
	}
	err = InsertStringToFile(directory+file, packageName+"\n", line+1)
	if err != nil {
		return 0, err
	}
	return line + 1, nil
}
func File2lines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LinesFromReader(f)
}

func LinesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

/**
 * Insert sting to n-th line of file.
 * If you want to insert a line, append newline '\n' to the end of the string.
 */
func InsertStringToFile(path, str string, index int) error {
	lines, err := File2lines(path)
	if err != nil {
		return err
	}

	fileContent := ""
	for i, line := range lines {
		if i == index {
			fileContent += str
		}
		fileContent += line
		fileContent += "\n"
	}

	return os.WriteFile(path, []byte(fileContent), 0644)
}
func findLineOfInsert(path string, text string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	line := 1
	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), text) {
			return line, nil
		}

		line++
	}

	if err := scanner.Err(); err != nil {
		// Handle the error
		return 0, err
	}
	return line, nil
}

func printAddedPackageRes(file string, line int) error {
	lines, err := File2lines(directory + file)
	if err != nil {
		return err
	}
	for i := line - 3; i < line+3; i++ {
		if i >= len(lines) {
			break
		}
		if i == line {
			d := color.New(color.FgWhite, color.Bold, color.BgRed)
			d.Printf("+++")
			fmt.Print(" ")
			d = color.New(color.FgBlack, color.Bold, color.BgGreen)
			d.Printf(lines[i] + "\n")
		} else {
			fmt.Printf(lines[i] + "\n")
		}
	}
	return nil
}

func main() {
	cmd := &cli.Command{
		EnableShellCompletion: true,
		Suggest:               true,

		Commands: []*cli.Command{
			{
				Name:    "build",
				Aliases: []string{"b"}, //add to other functions
				Usage:   "build config",
				Commands: []*cli.Command{
					{
						Name:    "local",
						Aliases: []string{"l"}, //add to other functions
						Usage:   "build on local machine",
						Action: func(_ context.Context, cmd *cli.Command) error {
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
								return nil
							}

							return nil
						},
					},
					{
						Name:    "remote",
						Usage:   "build on remote machine",
						Aliases: []string{"r"}, //add to other functions
						Action: func(_ context.Context, cmd *cli.Command) error {
							err := buildToRemote()
							if err != nil {
								return err
							}
							err = git()
							if err != nil {
								return nil
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
				Commands: []*cli.Command{{
					Name:    "local",
					Aliases: []string{"l"}, //add to other functions
					Usage:   "build on local machine",
					Action: func(_ context.Context, cmd *cli.Command) error {
						err := backupFlakeLock()
						if err != nil {
							fmt.Println("!!! CANT BACKUP flake.lock FILE")
						}

						err = update()
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
						return nil
					},
				},
					{
						Name:    "remote",
						Usage:   "build on remote machine",
						Aliases: []string{"r"}, //add to other functions
						Action: func(_ context.Context, cmd *cli.Command) error {
							err := update()
							if err != nil {
								return err
							}
							err = buildToRemote()
							if err != nil {
								return err
							}
							err = git()
							if err != nil {
								return nil
							}

							return nil
						},
					},
				},
			},

			{
				Name:    "add",
				Aliases: []string{"a"}, //add to other functions
				Usage:   "add package to config",
				Commands: []*cli.Command{{
					Name:    "homeManager",
					Aliases: []string{"hm"}, //add to other functions
					Usage:   "add homeManager package",
					Action: func(_ context.Context, cmd *cli.Command) error {
						pName := cmd.Args().Get(0)
						if pName == "" {
							return nil
						}
						relFile := "./home/packages.nix"
						line, err := addPackage(pName, relFile)
						if err != nil {
							return err
						}
						err = printAddedPackageRes(relFile, line)
						if err != nil {
							fmt.Println(err)
						}
						return nil
					},
				},
					{
						Name:    "system",
						Aliases: []string{"s"}, //add to other functions
						Usage:   "add system package",
						Action: func(_ context.Context, cmd *cli.Command) error {
							pName := cmd.Args().Get(0)
							if pName == "" {
								return nil
							}
							relFile := "./hosts/main/system-packages.nix"
							line, err := addPackage(pName, relFile)
							if err != nil {
								return err
							}
							err = printAddedPackageRes(relFile, line)
							if err != nil {
								fmt.Println(err)
							}
							return nil
						},
					},
				},
			},
		},
	}
	cmd.Run(context.Background(), os.Args)
	// if err := app.Run(os.Args); err != nil {
	// 	log.Fatal(err)
	// }
}
