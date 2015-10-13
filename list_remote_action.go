package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"net/http"
	"strings"
	"encoding/json"
	"container/list"
)


type VersionIndex []VersionDesc

type VersionDesc struct {
	Version string `json:"version"`
	Date string `json:"date"`
	Files []string `json:"files"`
	Npm string `json:"npm"`
	V8 string `json:"v8"`
	Uv string `json:"uv"`
	Zlib string `json:"zlib"`
	Openssl string `json:"openssl"`
	Modules string `json:"modules"`
}

func (this *VersionDesc) BelongsTo(versionName string) bool {

	if strings.HasPrefix(versionName, "v") {
		buf := []rune(versionName)
		v := string(buf[1:])

		versionName = v
	}

	if strings.HasPrefix(this.Version, "v") {
		buf := []rune(this.Version)
		v := string(buf[1:])

		return strings.HasPrefix(v, versionName)
	}

	return strings.HasPrefix(this.Version, versionName)
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

func ListRemoteAction(context *cli.Context) {

	resp, err := http.Get(DIST_URL + "/index.json")
	defer (func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	})()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	decoder := json.NewDecoder(resp.Body)
	var index VersionIndex

	err = decoder.Decode(&index)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	result := list.New()

	for _, versionDesc := range index {
		result.PushBack(versionDesc)
	}

	if len(context.Args()) > 0 {
		filterBy := context.Args().First()

		result = FilterList(result, func(e *list.Element) bool {
			desc := e.Value.(VersionDesc)
			return desc.BelongsTo(filterBy)
		})
	}

	for e := result.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value.(VersionDesc).Version)
	}
}
