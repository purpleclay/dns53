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
	"strings"
	"testing"

	"github.com/purpleclay/dns53/internal/imds"
	"github.com/purpleclay/dns53/internal/imds/imdsstub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagsCommand(t *testing.T) {
	var buf bytes.Buffer
	ctx := &globalContext{
		out:        &buf,
		imdsClient: imds.NewFromAPI(imdsstub.New(t)),
	}

	cmd := newTagsCommand()
	err := cmd.ExecuteContext(ctx)

	require.NoError(t, err)

	// Can't guarantee order of tags so break up the testing of the table
	table := buf.String()

	assert.True(t, strings.HasPrefix(table, `+-------------+----------+-----------------------+-------------------------------+
|     TAG     |  VALUE   |   PROPERTY CHAINING   |            INDEXED            |
+-------------+----------+-----------------------+-------------------------------+
`))
	assert.Contains(t, table, "| Name        | stub-ec2 | {{.Tags.Name}}        | {{index .Tags \"Name\"}}        |\n")
	assert.Contains(t, table, "| Environment | dev      | {{.Tags.Environment}} | {{index .Tags \"Environment\"}} |\n")
	assert.True(t, strings.HasSuffix(table, "+-------------+----------+-----------------------+-------------------------------+\n"))
}
