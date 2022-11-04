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

package wizard

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	UpDown       key.Binding
	LeftRight    key.Binding
	Enter        key.Binding
	ForwardSlash key.Binding
	Escape       key.Binding
	Quit         key.Binding
}

func DefaultKeyMap() *KeyMap {
	km := &KeyMap{}

	km.UpDown = key.NewBinding(
		key.WithKeys("up", "down"), key.WithHelp("↑↓", "up/down"),
	)
	km.LeftRight = key.NewBinding(
		key.WithKeys("←", "→"), key.WithHelp("← →", "left/right"),
	)
	km.Enter = key.NewBinding(
		key.WithKeys("↲"), key.WithHelp("↲", "select"),
	)
	km.ForwardSlash = key.NewBinding(
		key.WithKeys("/"), key.WithHelp("/", "filter"),
	)
	km.Escape = key.NewBinding(
		key.WithKeys("esc"), key.WithHelp("esc", "back"),
	)
	km.Quit = key.NewBinding(
		key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit"),
	)

	return km
}
