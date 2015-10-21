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


func LoadAction(context *cli.Context) {
	if len(context.Args()) == 0 {
		fmt.Println("no version given to load")
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
	nameTarGz := GetTarballName(version)

	err = DownloadAndSave(url, location, nameTarGz)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = UnGZip(location, nameTarGz)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	nameTar := nameTarGz[0:len(nameTarGz) - len(".gz")]

	destLocation := GetVersionsPath()
	EnsureDirExists(destLocation, PERM)

	err = UnTar(location, nameTar, destLocation)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = os.Remove(path.Join(location, nameTarGz))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = os.Remove(path.Join(location, nameTar))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	name := nameTar[0:len(nameTar) - len(".tar")]

	err = os.Rename(path.Join(destLocation, name), path.Join(destLocation, version))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("load complete")
}
