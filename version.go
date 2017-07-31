//go:generate go-build-info -p treksum -o build_info.go
package treksum

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	Version = "v0.0.4"

	VersionTag = "dev"

	sep = "-"
)

var AppInfo appInfo

type appInfo struct {
	App         string      `json:"app"`
	VersionInfo VersionInfo `json:"version_info"`
}

func (this *appInfo) String() string {
	return fmt.Sprintf("%s %s", this.App, this.VersionInfo.VersionString)
}

type VersionInfo struct {
	Version       string `json:"version"`
	VersionTag    string `json:"tag"`
	VersionString string `json:"version_string"`
	Revision      string `json:"revision"`
	Branch        string `json:"branch"`
	BuildUser     string `json:"build_user"`
	BuildDate     string `json:"build_date"`
	GoVersion     string `json:"go_version"`
}

func init() {
	AppInfo = appInfo{
		VersionInfo: VersionInfo{
			Version:       Version,
			VersionTag:    VersionTag,
			VersionString: GetVersionString(),
			Revision:      BuildRevision,
			Branch:        BuildBranch,
			BuildUser:     BuildUser,
			BuildDate:     BuildDate,
			GoVersion:     runtime.Version(),
		},
	}
}

func GetVersionString() string {
	return strings.TrimRight(strings.Join([]string{Version, VersionTag}, sep), " "+sep)
}
