/*
Copyright (c) 2022 Purple Clay

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// The current built version
	version = ""
	// The git branch associated with the current built version
	gitBranch = ""
	// The git SHA1 of the commit
	gitCommit = ""
	// The date associated with the current built version
	buildDate = ""
)

// BuildInfo contains build time information about the application
type BuildInfo struct {
	Version   string    `json:"version,omitempty"`
	GitBranch string    `json:"gitBranch,omitempty"`
	GitCommit string    `json:"gitCommit,omitempty"`
	BuildDate string    `json:"buildDate,omitempty"`
	Go        GoDetails `json:"go,omitempty"`
}

// GoDetails contains Go compile time information
type GoDetails struct {
	Version string `json:"version,omitempty"`
	OS      string `json:"os,omitempty"`
	Arch    string `json:"arch,omitempty"`
}

// Short returns the semantic version of the application
func Short() string {
	return version
}

// Long returns the build time version information of the application
func Long() BuildInfo {
	return BuildInfo{
		Version:   version,
		GitBranch: gitBranch,
		GitCommit: gitCommit,
		BuildDate: buildDate,
		Go: GoDetails{
			Version: runtime.Version(),
			OS:      runtime.GOOS,
			Arch:    runtime.GOARCH,
		},
	}
}

type versionOptions struct {
	short bool
}

func newVersionCmd(out io.Writer) *cobra.Command {
	opts := versionOptions{}

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Prints the build time version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.run(out)
		},
	}

	f := cmd.Flags()
	f.BoolVar(&opts.short, "short", false, "only print the semantic version number")

	return cmd
}

func (o versionOptions) run(out io.Writer) error {
	fmt.Fprintln(out, formatVersion(o.short))
	return nil
}

func formatVersion(short bool) string {
	if short {
		return Short()
	}

	out, _ := json.Marshal(Long())
	return string(out)
}
