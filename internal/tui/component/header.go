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

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/purpleclay/lipgloss-theme"
)

var borderBottom = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder(), false, false, true, false).
	BorderForeground(theme.S600).
	MarginBottom(1)

type Header struct {
	name        string
	version     string
	description string
	width       int
}

func NewHeader(name, version, description string) *Header {
	return &Header{
		name:        name,
		description: description,
		version:     version,
	}
}

func (*Header) Init() tea.Cmd {
	return nil
}

func (m *Header) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Header) View() string {
	var b strings.Builder

	nameVersion := lipgloss.JoinHorizontal(
		lipgloss.Left,
		theme.H2.Render(m.name),
		theme.H4.Render(m.version),
	)

	banner := lipgloss.JoinVertical(
		lipgloss.Top,
		nameVersion,
		"",
		faint.Render(m.description),
	)

	b.WriteString(borderBottom.Width(m.width).Render(banner))
	return b.String()
}

func (m *Header) Resize(width, _ int) Model {
	m.width = width
	return m
}

func (m *Header) Width() int {
	return m.width
}

func (m *Header) Height() int {
	return lipgloss.Height(m.View())
}
