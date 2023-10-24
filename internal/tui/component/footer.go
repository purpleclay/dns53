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

package component

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/purpleclay/lipgloss-theme"
)

var (
	faint     = lipgloss.NewStyle().Faint(true)
	borderTop = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(theme.S600)
)

type Footer struct {
	help   help.Model
	keymap help.KeyMap
	width  int
}

func NewFooter(keymap help.KeyMap) Footer {
	help := help.New()
	help.Styles.ShortSeparator = lipgloss.NewStyle().Foreground(theme.S100)
	help.Styles.ShortKey = lipgloss.NewStyle()
	help.Styles.ShortDesc = faint.Copy()

	return Footer{
		help:   help,
		keymap: keymap,
	}
}

func (Footer) Init() tea.Cmd {
	return nil
}

func (m Footer) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Footer) View() string {
	var b strings.Builder
	panel := lipgloss.JoinVertical(lipgloss.Top, m.help.View(m.keymap))

	b.WriteString(borderTop.Width(m.width).Render(panel))
	return b.String()
}

func (m Footer) Resize(width, _ int) Model {
	m.width = width
	return m
}

func (m Footer) Width() int {
	return m.width
}

func (m Footer) Height() int {
	return lipgloss.Height(m.View())
}

func (m Footer) SetKeyMap(keymap help.KeyMap) Footer {
	m.keymap = keymap
	return m
}
