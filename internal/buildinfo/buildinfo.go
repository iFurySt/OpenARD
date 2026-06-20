package buildinfo

import (
	"encoding/json"
	"fmt"
)

var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

type Info struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"buildDate"`
}

func Current() Info {
	return Info{
		Version: Version,
		Commit:  Commit,
		Date:    Date,
	}
}

func (info Info) String() string {
	return fmt.Sprintf("version=%s commit=%s buildDate=%s", info.Version, info.Commit, info.Date)
}

func (info Info) JSON() ([]byte, error) {
	return json.MarshalIndent(info, "", "  ")
}
