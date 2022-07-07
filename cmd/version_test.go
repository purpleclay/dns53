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
	"bytes"
	"encoding/json"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersion_Long(t *testing.T) {
	version = "v0.1.0"
	gitBranch = "main"
	gitCommit = "d4b3bd00406444561e646607d7f941097dbd1b40"
	buildDate = "2022-06-29T20:05:51Z"

	var buf bytes.Buffer
	cmd := newVersionCmd(&buf)

	err := cmd.Execute()
	require.NoError(t, err)

	assert.NotEmpty(t, buf.String())

	var bi BuildInfo
	json.Unmarshal(buf.Bytes(), &bi)
	assert.Equal(t, version, bi.Version)
	assert.Equal(t, gitBranch, bi.GitBranch)
	assert.Equal(t, gitCommit, bi.GitCommit)
	assert.Equal(t, buildDate, bi.BuildDate)
	assert.Equal(t, runtime.Version(), bi.Go.Version)
	assert.Equal(t, runtime.GOOS, bi.Go.OS)
	assert.Equal(t, runtime.GOARCH, bi.Go.Arch)
}

func TestVersion_Short(t *testing.T) {
	version = "v0.1.0"

	var buf bytes.Buffer
	cmd := newVersionCmd(&buf)
	cmd.SetArgs([]string{"--short"})

	err := cmd.Execute()
	require.NoError(t, err)

	assert.NotEmpty(t, buf.String())
	assert.Contains(t, buf.String(), version)
}
