package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"container/list"
	"strings"
	"errors"
	"os"
	"path"
)


func InstallAction(context *cli.Context) {
	if len(context.Args()) == 0 {
		fmt.Println("no version given to install")
		return
	}

	requestedVersionName := context.Args().First()
	availableVersions, err := GetRemoteVersions()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	suitableVersions := FilterList(availableVersions, func(el *list.Element) bool {
		vd, ok := el.Value.(VersionDesc)
		return ok && vd.BelongsTo(requestedVersionName)
	})

	ve := suitableVersions.Back()
	if ve == nil {
		fmt.Println("no suitable version found")
		return
	}

	vd, ok := ve.Value.(VersionDesc)
	if !ok {
		fmt.Println(errors.New("unable to typecast list element").Error())
		return
	}

	version := vd.Version

	url := strings.Join([]string{DIST_URL, version, GetTarballName(version)}, "/")
	location := GetDataPath()
	name := GetTarballName(version)

	err = DownloadAndSave(url, location, name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = UnGZip(location, name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	name = name[0:len(name) - 3]

	err = UnTar(location, name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = os.Remove(path.Join(location, name))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = os.Remove(path.Join(location, name + ".gz"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("installation complete")
}
