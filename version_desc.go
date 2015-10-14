package main

import (
	"strings"
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
