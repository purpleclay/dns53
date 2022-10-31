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

package header

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/purpleclay/dns53/internal/tui/components"
)

// Model ...
type Model struct {
	name        string
	version     string
	description string
	width       int
	Styles      *Styles
}

// New ...
func New(name, version, description string) Model {
	return Model{
		name:        name,
		description: description,
		version:     version,
		Styles:      DefaultStyles(),
	}
}

// Init ...
func (m Model) Init() tea.Cmd {
	return nil
}

// Update ...
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

// View ...
func (m Model) View() string {
	var b strings.Builder

	nameVersion := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.Styles.Name.Padding(0, 2).Render(m.name),
		m.Styles.Version.Padding(0, 2).Render(m.version),
	)

	desc := m.Styles.Description.Render(m.description)

	banner := lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().MarginBottom(1).Render(nameVersion),
		desc,
	)

	b.WriteString(m.Styles.Border.
		MarginBottom(1).
		Width(m.width).
		Render(banner))

	return b.String()
}

// Resize ...
func (m Model) Resize(width, height int) components.Model {
	m.width = width
	return m
}

// Width ...
func (m Model) Width() int {
	return m.width
}

// Height ...
func (m Model) Height() int {
	return lipgloss.Height(m.View())
}
