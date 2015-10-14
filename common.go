package main

import (
	"os"
	"path"
	"container/list"
	"encoding/json"
	"net/http"
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

	if (result != nil) {
		result = os.Symlink(path.Join(versionPath, BIN), path.Join(getDataPath(), BIN))
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
	resp, err := http.Get(DIST_URL + "/index.json")
	defer (func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	})()

	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(resp.Body)
	var index VersionIndex

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
