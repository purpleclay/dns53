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

package dashboard

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/purpleclay/dns53/internal/tui/theme"
)

// Styles ...
type Styles struct {
	PrimaryLabel   lipgloss.Style
	SecondaryLabel lipgloss.Style
	Spacing        string
	PendingStatus  lipgloss.Style
	ActiveStatus   lipgloss.Style
	Highlight      lipgloss.Style
}

// DefaultStyles ...
func DefaultStyles() *Styles {
	s := &Styles{}

	s.PrimaryLabel = lipgloss.NewStyle().
		Background(theme.PrimaryColour).
		Foreground(theme.TextColour).
		Bold(true).
		Padding(0, 3).
		MarginRight(8)

	s.SecondaryLabel = lipgloss.NewStyle().
		Background(theme.SecondaryColour).
		Foreground(theme.TextColour).
		PaddingLeft(2).
		Width(11).
		Bold(true)

	s.Spacing = lipgloss.NewStyle().
		Foreground(theme.VeryFeintColour).
		Render(strings.Repeat(".", 12))

	s.PendingStatus = lipgloss.NewStyle().
		Foreground(theme.AmberColour).
		Bold(true)

	s.ActiveStatus = lipgloss.NewStyle().
		Foreground(theme.GreenColour).
		Bold(true)

	s.Highlight = lipgloss.NewStyle().Foreground(theme.HighlightColour).Bold(true)

	return s
}
