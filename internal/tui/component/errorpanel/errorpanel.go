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

package errorpanel

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

type Model struct {
	reason string
	cause  string
	width  int
	Styles *Styles
}

func New() Model {
	return Model{
		Styles: DefaultStyles(),
	}
}

func (Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	reason := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.Styles.Label.Padding(0, 2).Render("Error"),
		m.Styles.Reason.Render(": "+m.reason),
	)

	desc := m.Styles.Cause.MarginLeft(1).Render(wordwrap.String(m.cause, m.width))

	panel := lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().Margin(0, 0, 1, 1).Render(wordwrap.String(reason, m.width)),
		desc,
	)

	b.WriteString(m.Styles.Border.Render(panel))
	return b.String()
}

func (m Model) RaiseError(reason string, cause error) Model {
	m.reason = reason
	if cause != nil {
		m.cause = cause.Error()
	}

	return m
}

func (m Model) Resize(width, _ int) Model {
	// Restrict the error panel to be 3/4 the width of the containing component
	m.width = int(float32(width) * 0.75)
	return m
}

func (m Model) Width() int {
	return m.width
}

func (m Model) Height() int {
	return lipgloss.Height(m.View())
}
