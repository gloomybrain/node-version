package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"path"
	"path/filepath"
	"strings"
)


func DefaultAction(c *cli.Context) {
	versionName := c.Args().First()
	if versionName != "" {
		setDefault(versionName)
	}else {
		printDefault()
	}
}

func setDefault(versionName string) {
	localVersionNames, err := GetLocalVersions()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, name := range *localVersionNames {
		if strings.Contains(name, versionName) {
			pathToLink := path.Join(GetDataPath(), DEFAULT)

			err := os.Remove(pathToLink)
			if err != nil && !os.IsNotExist(err) {
				fmt.Println(err.Error())
				return
			}

			err = os.Symlink(path.Join(GetVersionsPath(), name), pathToLink)
			if err != nil {
				fmt.Println(err.Error())
			}

			return
		}
	}

	fmt.Println("looks like there is no such version loaded: " + versionName)
}

func printDefault() {
	pathToLink := path.Join(GetDataPath(), DEFAULT)
	info, err := os.Lstat(pathToLink)

	if os.IsNotExist(err) {
		return
	}

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if info.Mode() & os.ModeSymlink == 0 {
		fmt.Println(pathToLink + " should be a symbolic link")
		return
	}

	pathToCurrent, err := filepath.EvalSymlinks(pathToLink)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if !strings.HasPrefix(pathToCurrent, GetVersionsPath()) {
		fmt.Println("current points to somewhere out of node-version scope")
		return
	}

	versionName := pathToCurrent[len(GetVersionsPath()) + 1:]
	fmt.Println(versionName)
}
