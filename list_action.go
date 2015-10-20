package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"strings"
)


func ListAction(c *cli.Context) {
	infoList, err := ioutil.ReadDir(GetVersionsPath())

	if os.IsNotExist(err) {
		fmt.Println("nothing here")
		return
	}

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	filter := c.Args().First()

	if filter != nil {
		fmt.Println("currently installed versions (filtered):")
	}else {
		fmt.Println("currently installed versions:")
	}

	for _, info := range infoList {

		if filter != "" {
			if strings.Contains(info.Name(), filter) {
				fmt.Println(info.Name())
			}
		}else {
			fmt.Println(info.Name())
		}
	}
}
