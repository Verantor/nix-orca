package main

import (
	"errors"
	"fmt"
)

func buildOS(remote bool, directory string) error {
	if remote {
		err := buildToRemote(directory)
		if err != nil {
			return err
		}
		err = git(directory)

		if err != nil {
			fmt.Println("!!! CANT COMMIT GIT CHANGES")
		}
		return nil
	} else {

		err := build(directory)
		if err != nil {
			return err
		}
		err = diff(directory)
		if err != nil {
			return err
		}
		err = activate(directory)
		if err != nil {
			return err
		}
		err = remove(directory)
		if err != nil {
			return err
		}
		err = git(directory)

		if err != nil {
			fmt.Println("!!! CANT COMMIT GIT CHANGES")
		}
		return nil
	}
}
func updateOS(remote bool, directory string) error {
	if remote {

		err := backupFlakeLock(directory)
		if err != nil {
			fmt.Println("!!! CANT BACKUP flake.lock FILE")
		}
		err = update(directory)
		if err != nil {
			return err
		}
		err = buildOS(remote, directory)
		if err != nil {
			return err
		}
		return nil
	} else {

		err := backupFlakeLock(directory)
		if err != nil {
			fmt.Println("!!! CANT BACKUP flake.lock FILE")
		}
		err = update(directory)
		if err != nil {
			return err
		}
		err = buildOS(remote, directory)
		if err != nil {
			return err
		}
		return nil
	}
}
func addPackageToOS(pName string, relFile string, directory string) error {
	if pName == "" {
		return errors.New("empty Package Name")
	}

	line, err := addPackage(pName, directory+relFile)
	if err != nil {
		return err
	}
	err = printAddedPackageRes(directory+relFile, line)
	if err != nil {
		return err
	}
	return nil

}
