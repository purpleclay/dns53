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
	"github.com/charmbracelet/lipgloss"
	"github.com/purpleclay/dns53/internal/tui/theme"
)

type Styles struct {
	Label  lipgloss.Style
	Reason lipgloss.Style
	Cause  lipgloss.Style
	Border lipgloss.Style
}

func DefaultStyles() *Styles {
	s := &Styles{}

	s.Label = lipgloss.NewStyle().
		Background(theme.RedColour).
		Foreground(theme.TextColour).
		Bold(true)

	s.Reason = lipgloss.NewStyle().
		Foreground(theme.TextColour).
		Bold(true)

	s.Cause = lipgloss.NewStyle().
		Foreground(theme.TextColour)

	s.Border = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(theme.BorderColour)

	return s
}
