package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"strings"
)


func ListAction(c *cli.Context) {
	versions, err := GetLocalVersions()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	filter := c.Args().First()

	if filter != "" {
		versions = FilterStringSlice(versions, func(element *string) bool {
			return strings.Contains(*element, filter)
		})
	}

	for _, v := range *versions {
		fmt.Println(v)
	}
}
