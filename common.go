package main

import (
	"os"
	"path"
	"container/list"
	"encoding/json"
	"net/http"
	"strings"
	"errors"
)


const DIST_URL string = "https://nodejs.org/dist"
const DATA_DIR string = ".node-version"
const BIN string = "bin"
const VERSIONS string = "versions"


func GetDataPath() string {
	return path.Join(os.Getenv("HOME"), DATA_DIR)
}

func GetTarballName(versionName string) string {
	return "node-" + versionName + "-darwin-x64.tar.gz"
}

func MakeVersionDir(versionName string) error {
	p := path.Join(GetDataPath(), VERSIONS, versionName)
	return os.MkdirAll(p, 0755)
}

func LinkDefaultVersion(versionName string) error {
	var result error = nil

	versionPath := path.Join(GetDataPath(), VERSIONS, versionName)
	_, result = os.Stat(versionPath)

	if (result != nil) {
		result = os.Symlink(path.Join(versionPath, BIN), path.Join(GetDataPath(), BIN))
	}

	return result
}

func FilterList(source *list.List, filter func (e *list.Element) bool) *list.List {
	result := list.New()

	for e := source.Front(); e != nil; e = e.Next() {
		if filter(e) {
			result.PushBack(e.Value)
		}
	}

	return result
}

func GetRemoteVersions() (*list.List, error) {
	url := strings.Join([]string{DIST_URL, "index.json" }, "/")
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if resp.Body == nil {
		return nil, errors.New("some network problem had occured")
	} else {
		defer resp.Body.Close()
	}

	decoder := json.NewDecoder(resp.Body)

	index := VersionIndex{}

	err = decoder.Decode(&index)

	if err != nil {
		return nil, err
	}

	result := list.New()

	for _, versionDesc := range index {
		result.PushBack(versionDesc)
	}

	return result, nil
}
