package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

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
