package main

import (
	"os"
	"path"
	"container/list"
	"encoding/json"
	"net/http"
	"strings"
	"errors"
	"io"
	"compress/gzip"
	"fmt"
	"archive/tar"
	"io/ioutil"
)


const DIST_URL string = "https://nodejs.org/dist"
const DATA_DIR string = ".node-version"
const VERSIONS = "versions"
const BIN string = "bin"
const PERM uint32 = 0744
const DEFAULT string = "default"


func GetDataPath() string {
	return path.Join(os.Getenv("HOME"), DATA_DIR)
}

func GetVersionsPath() string {
	return path.Join(GetDataPath(), VERSIONS)
}

func EnsureDirExists(p string, perm uint32) error {
	info, err := os.Stat(p)
	fileMode := os.ModeDir | os.FileMode(perm)

	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("creating directory " + p)
			return os.MkdirAll(p, fileMode)
		} else {
			return err
		}
	} else if !info.IsDir() {
		return errors.New(p + " is not a directory")
	} else if info.Mode() != fileMode {
		return errors.New(p + " has wrong permissions")
	}

	return nil
}

func GetTarballName(versionName string) string {
	return "node-" + versionName + "-darwin-x64.tar.gz"
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

func FilterStringSlice(source *[]string, filter func(element *string) bool) *[]string {
	var result []string

	for _, v := range *source {
		if filter(&v) {
			result = append(result, v)
		}
	}

	return &result
}

func GetRemoteVersions() (*list.List, error) {
	url := strings.Join([]string{DIST_URL, "index.json" }, "/")

	fmt.Println("loading versions... ")
	fmt.Println(url)

	resp, err := http.Get(url)

	if err != nil { return nil, err }

	if resp.Body == nil {
		return nil, errors.New("some network problem had occured")
	}else {
		defer resp.Body.Close()
	}

	decoder := json.NewDecoder(resp.Body)

	index := VersionIndex{}

	err = decoder.Decode(&index)
	if err != nil { return nil, err }

	result := list.New()

	for _, e := range index {
		result.PushBack(e)
	}

	fmt.Println("done")
	return result, nil
}

func GetLocalVersions() (*[]string, error) {
	infoList, err := ioutil.ReadDir(GetVersionsPath())

	if os.IsNotExist(err) {
		result := make([]string, 0)
		return &result, nil
	}

	if err != nil {
		return nil, err
	}

	result := make([]string, len(infoList))

	for index, info := range infoList {
		result[index] = info.Name()
	}

	return &result, nil
}

func DownloadAndSave(url, location, name string) error {
	fmt.Println("loading installation archive... ")
	fmt.Println(url)

	resp, err := http.Get(url)

	if err != nil { return err }
	if resp.Body == nil { return errors.New("some network problem has occured") }

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return errors.New("unexpected response status code " + resp.Status)
	}

	err = EnsureDirExists(location, PERM)
	if err != nil { return err }

	filePath := path.Join(location, name)

	file, err := os.OpenFile(filePath, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, os.FileMode(PERM))
	if err != nil { return err }

	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil { return err }

	fmt.Println("done")
	return nil
}

func UnGZip(location, fileName string) error {
	filePath := path.Join(location, fileName)

	fmt.Println("uncompressing file...")
	fmt.Println(filePath)

	_, err := os.Stat(filePath)
	if err != nil { return err }

	sourceFile, err := os.Open(filePath)
	if err != nil { return err }
	defer sourceFile.Close()

	gzipReader, err := gzip.NewReader(sourceFile)
	if err != nil { return err }
	defer gzipReader.Close()

	if !strings.HasSuffix(filePath, ".gz") {
		return errors.New("Unsupported tarball file extension")
	}

	destFilePath := filePath[0:len(filePath) - 3]
	destFile, err := os.OpenFile(destFilePath, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, os.FileMode(PERM))
	if err != nil { return err }
	defer destFile.Close()

	_, err = io.Copy(destFile, gzipReader)
	if err != nil { return err }

	fmt.Println("done")
	return nil
}

func UnTar(location, fileName, destLocation string) error {
	sourcePath := path.Join(location, fileName)

	fmt.Println("opening file...")
	fmt.Println(sourcePath)

	_, err := os.Stat(sourcePath)
	if err != nil { return err }

	sourceFile, err := os.Open(sourcePath)
	if err != nil { return err }
	defer sourceFile.Close()

	tarReadder := tar.NewReader(sourceFile)

	for {
		hdr, err := tarReadder.Next()

		if err == io.EOF { break }
		if err != nil { return err }

		destPath := path.Join(destLocation, hdr.Name)

		if strings.HasSuffix(hdr.Name, "/") && hdr.Size == 0 {
			err = os.MkdirAll(destPath, os.ModeDir | hdr.FileInfo().Mode())
			if err != nil { return err }
		}else {
			file, err := os.OpenFile(destPath, os.O_CREATE | os.O_WRONLY, hdr.FileInfo().Mode())
			if err != nil { return err }

			_, err = io.Copy(file, tarReadder)
			if err != nil { return err }

			err = file.Close()
			if err != nil { return err }
		}
	}

	return nil
}
