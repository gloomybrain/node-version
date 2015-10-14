package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"container/list"
)

func ListRemoteAction(context *cli.Context) {

	versions, err := GetRemoteVersions()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(context.Args()) > 0 {
		filterBy := context.Args().First()

		versions = FilterList(versions, func(e *list.Element) bool {
			desc := e.Value.(VersionDesc)
			return desc.BelongsTo(filterBy)
		})
	}

	for e := versions.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value.(VersionDesc).Version)
	}
}
