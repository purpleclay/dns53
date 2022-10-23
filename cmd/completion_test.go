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

// func TestCompletionBash(t *testing.T) {
// 	var buf bytes.Buffer
// 	cmd := newCompletionCmd(&buf)
// 	cmd.SetArgs([]string{"bash"})

// 	err := cmd.Execute()
// 	require.NoError(t, err)

// 	assert.NotEmpty(t, buf.String())
// 	assert.Contains(t, buf.String(), "bash completion V2 for completion ")
// }

// func TestCompletionZsh(t *testing.T) {
// 	var buf bytes.Buffer
// 	cmd := newCompletionCmd(&buf)
// 	cmd.SetArgs([]string{"zsh"})

// 	err := cmd.Execute()
// 	require.NoError(t, err)

// 	assert.NotEmpty(t, buf.String())
// 	assert.Contains(t, buf.String(), "zsh completion for completion")
// }

// func TestCompletionZshNoDescriptions(t *testing.T) {
// 	var buf bytes.Buffer
// 	cmd := newCompletionCmd(&buf)
// 	cmd.SetArgs([]string{"zsh", "--no-descriptions"})

// 	err := cmd.Execute()
// 	require.NoError(t, err)

// 	assert.NotEmpty(t, buf.String())
// 	assert.Contains(t, buf.String(), "zsh completion for completion")
// 	assert.Contains(t, buf.String(), "__completeNoDesc")
// }

// func TestCompletionFish(t *testing.T) {
// 	var buf bytes.Buffer
// 	cmd := newCompletionCmd(&buf)
// 	cmd.SetArgs([]string{"fish"})

// 	err := cmd.Execute()
// 	require.NoError(t, err)

// 	assert.NotEmpty(t, buf.String())
// 	assert.Contains(t, buf.String(), "fish completion for completion")
// }
