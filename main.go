package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
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

	scanner := bufio.NewScanner(f)

	line := 1
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), text) {
			return line, nil
		}

		line++
	}

	if err := scanner.Err(); err != nil {
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
		CommandNotFound: func(ctx context.Context, cmd *cli.Command, command string) {
			fmt.Fprintf(cmd.Root().Writer, "Thar be no %q here.\n", command)
		},
		OnUsageError: func(ctx context.Context, cmd *cli.Command, err error, isSubcommand bool) error {
			if isSubcommand {
				return err
			}

			fmt.Fprintf(cmd.Root().Writer, "WRONG: %#v\n", err)
			return nil
		},

		Commands: []*cli.Command{
			{
				Name:  "completion",
				Usage: "Output shell completion script for bash, zsh, or fish",
				Description: `Output shell completion script for bash, zsh, or fish.
		Source the output to enable completion.
		nix-orca completion fish > ~/.config/fish/completions/nix-orca.fish
		`,
				Commands: []*cli.Command{
					// {
					// 	Name:   "bash",
					// 	Usage:  "Output shell completion script for bash",
					// 	Action: r.bashCompletionAction,
					// },
					// {
					// 	Name:   "zsh",
					// 	Usage:  "Output shell completion script for zsh",
					// 	Action: r.zshCompletionAction,
					// },
					{
						Name:   "fish",
						Usage:  "Output shell completion script for fish",
						Action: fishCompletionAction,
					},
				},
			},
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
							err := buildOS(false)
							return err
						},
					},
					{
						Name:    "remote",
						Usage:   "build on remote machine",
						Aliases: []string{"r"}, //add to other functions
						Action: func(_ context.Context, cmd *cli.Command) error {
							err := buildOS(true)
							return err
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
						err := updateOS(false)
						return err

					},
				},
					{
						Name:    "remote",
						Usage:   "build on remote machine",
						Aliases: []string{"r"}, //add to other functions
						Action: func(_ context.Context, cmd *cli.Command) error {
							err := updateOS(true)
							return err
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
						relFile := "./home/packages.nix"
						err := addPackageToOS(cmd.Args().Get(0), relFile)
						return err
					},
				},
					{
						Name:    "system",
						Aliases: []string{"s"}, //add to other functions
						Usage:   "add system package",
						Action: func(_ context.Context, cmd *cli.Command) error {
							relFile := "./hosts/main/system-packages.nix"
							err := addPackageToOS(cmd.Args().Get(0), relFile)
							return err
						},
					},
				},
			},
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
