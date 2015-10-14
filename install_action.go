package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"container/list"
	"strings"
	"net/http"
	"os"
	"path"
	"io"
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
		vd := el.Value.(VersionDesc)
		return vd.BelongsTo(requestedVersionName)
	})

	versionEl := suitableVersions.Back()
	if versionEl == nil {
		fmt.Println("no suitable version found")
		return
	}

	version := versionEl.Value.(VersionDesc).Version
	url := strings.Join([]string{DIST_URL, version, GetTarballName(version)}, "/")

	req, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if req.Body == nil {
		fmt.Println("some network problem has occured")
		return
	}

	defer req.Body.Close()


	fileName := path.Join(GetDataPath(), GetTarballName(version))
	fileFlag := os.O_WRONLY
	fileMode := os.FileMode(0744)

	err = os.Mkdir(GetDataPath(), fileMode)
	if err != nil {
		fmt.Println(err.Error())
		return
	}


	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	file.Close()

	err = os.Chmod(fileName, fileMode)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	file, err = os.OpenFile(fileName, fileFlag, fileMode)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer file.Close()

	reading := true
	buf := make([]byte, 4096)

	for reading {
		n_r, err_r := req.Body.Read(buf)

		if err_r == nil || err_r == io.EOF {

			for sum_w := 0; sum_w < n_r; {
				n_w, err_w := file.Write(buf[0:n_r])

				if err_w != nil {
					fmt.Println(err_w.Error())
					return
				}

				sum_w += n_w
			}

			if err_r == io.EOF {
				reading = false
				break
			}

		} else {
			fmt.Println(err_r.Error())
			return
		}
	}

	fmt.Println("file succesfully downloaded")
}
