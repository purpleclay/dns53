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

package errorpanel

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model ...
type Model struct {
	reason string
	cause  string
	raised bool
	width  int
	Styles *Styles
}

// New ...
func New() Model {
	return Model{
		raised: false,
		Styles: DefaultStyles(),
	}
}

// Init ...
func (m Model) Init() tea.Cmd {
	return nil
}

// Update ...
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = int(float32(msg.Width) * 0.75)
	}
	return m, nil
}

// View ...
func (m Model) View() string {
	if m.raised {
		var b strings.Builder

		reason := lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.Styles.Label.Padding(0, 2).Render("Error"),
			m.Styles.Reason.Width(m.width).Render(": "+m.reason),
		)

		desc := m.Styles.Cause.Width(m.width).MarginLeft(1).Render(m.cause)

		panel := lipgloss.JoinVertical(
			lipgloss.Top,
			lipgloss.NewStyle().Width(m.width).Margin(0, 0, 1, 1).Render(reason),
			desc,
		)

		b.WriteString(m.Styles.Border.Render(panel))
		return b.String()
	}

	return ""
}

// RaiseError ...
func (m *Model) RaiseError(reason string, cause error) {
	m.raised = true
	m.reason = reason
	if cause != nil {
		m.cause = cause.Error()
	}
}
