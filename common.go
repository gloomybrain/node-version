package main

import (
	"os"
	"path"
)


const DATA_DIR_NAME string = ".gnvm"


func GetUserPath() string {
	return os.Getenv("HOME")
}

func GetDataPath() string {
	return path.Join(GetUserPath(), DATA_DIR_NAME)
}

func MakeVersionDir(versionName string) error {
	path := path.Join(GetDataPath(), versionName)
	return os.MkdirAll(path, 0755)
}