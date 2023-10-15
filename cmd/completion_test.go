/*
Copyright (c) 2022 - 2023 Purple Clay

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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompletionCommandBash(t *testing.T) {
	var buf bytes.Buffer
	ctx := &globalContext{
		out: &buf,
	}

	cmd := newCompletionCmd()
	cmd.SetArgs([]string{"bash"})

	err := cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	assert.NotEmpty(t, buf.String())
	assert.Contains(t, buf.String(), "bash completion V2 for completion ")
}

func TestCompletionCommandZsh(t *testing.T) {
	var buf bytes.Buffer
	ctx := &globalContext{
		out: &buf,
	}

	cmd := newCompletionCmd()
	cmd.SetArgs([]string{"zsh"})

	err := cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	assert.NotEmpty(t, buf.String())
	assert.Contains(t, buf.String(), "zsh completion for completion")
}

func TestCompletionCommandZshNoDescriptions(t *testing.T) {
	var buf bytes.Buffer
	ctx := &globalContext{
		out: &buf,
	}

	cmd := newCompletionCmd()
	cmd.SetArgs([]string{"zsh", "--no-descriptions"})

	err := cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	assert.NotEmpty(t, buf.String())
	assert.Contains(t, buf.String(), "zsh completion for completion")
	assert.Contains(t, buf.String(), "__completeNoDesc")
}

func TestCompletionCommandFish(t *testing.T) {
	var buf bytes.Buffer
	ctx := &globalContext{
		out: &buf,
	}

	cmd := newCompletionCmd()
	cmd.SetArgs([]string{"fish"})

	err := cmd.ExecuteContext(ctx)
	require.NoError(t, err)

	assert.NotEmpty(t, buf.String())
	assert.Contains(t, buf.String(), "fish completion for completion")
}
