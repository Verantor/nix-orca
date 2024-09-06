package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

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

func fishCompletionAction(ctx context.Context, c *cli.Command) error {
	s, err := c.ToFishCompletion()
	if err != nil {
		return fmt.Errorf("generate fish completion: %w", err)
	}
	fmt.Fprintln(os.Stdout, s)
	return nil
}
