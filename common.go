package main

import (
	"os"
	"path"
)


const DIST_URL string = "https://nodejs.org/dist/"
const DATA_DIR string = ".gnvm"
const BIN string = "bin"
const VERSIONS string = "versions"


func getDataPath() string {
	return path.Join(os.Getenv("HOME"), DATA_DIR)
}

func MakeVersionDir(versionName string) error {
	p := path.Join(getDataPath(), VERSIONS, versionName)
	return os.MkdirAll(p, 0755)
}

func LinkDefaultVersion(versionName string) error {
	var result error = nil

	versionPath := path.Join(getDataPath(), VERSIONS, versionName)
	_, result = os.Stat(versionPath)

	if (result == nil) {
		result = os.Symlink(path.Join(versionPath, BIN), path.Join(getDataPath(), BIN))
	}

	return result
}
