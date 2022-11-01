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

package footer

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/purpleclay/dns53/internal/tui/components"
)

type Model struct {
	text   string
	width  int
	Styles *Styles
}

func New() Model {
	return Model{
		text:   "Made with 💜 at Purple Clay",
		Styles: DefaultStyles(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	// Border
	// Margin 1
	// Help Keys
	// Margin 1
	// Text (full block of colour)

	b.WriteString(m.Styles.Border.
		Width(m.width).
		Render(m.text))

	return b.String()
}

func (m Model) Resize(width, height int) components.Model {
	m.width = width
	return m
}

func (m Model) Width() int {
	return m.width
}

func (m Model) Height() int {
	return lipgloss.Height(m.View())
}
