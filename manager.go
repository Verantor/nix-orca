package main

import (
	"errors"
	"fmt"
)

func buildOS(remote bool) error {
	if remote {
		err := buildToRemote()
		if err != nil {
			return err
		}
		err = git()

		if err != nil {
			fmt.Println("!!! CANT COMMIT GIT CHANGES")
		}
		return nil
	} else {

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
			fmt.Println("!!! CANT COMMIT GIT CHANGES")
		}
		return nil
	}
}
func updateOS(remote bool) error {
	if remote {

		err := backupFlakeLock()
		if err != nil {
			fmt.Println("!!! CANT BACKUP flake.lock FILE")
		}
		err = update()
		if err != nil {
			return err
		}
		err = buildOS(remote)
		if err != nil {
			return err
		}
		return nil
	} else {

		err := backupFlakeLock()
		if err != nil {
			fmt.Println("!!! CANT BACKUP flake.lock FILE")
		}
		err = update()
		if err != nil {
			return err
		}
		err = buildOS(remote)
		if err != nil {
			return err
		}
		return nil
	}
}
func addPackageToOS(pName string, relFile string) error {
	if pName == "" {
		return errors.New("empty Package Name")
	}

	line, err := addPackage(pName, relFile)
	if err != nil {
		return err
	}
	err = printAddedPackageRes(relFile, line)
	if err != nil {
		return err
	}
	return nil

}
