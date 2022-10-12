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

package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	PrimaryColour        = lipgloss.Color("#3a1577")
	SecondaryColour      = lipgloss.Color("#5b1689")
	BorderColour         = lipgloss.Color("#2b0940")
	FeintColour          = lipgloss.Color("#807d8a")
	TextColour           = lipgloss.Color("#f6f5fc")
	HighlightColour      = lipgloss.Color("#bf31f7")
	HighlightFeintColour = lipgloss.Color("#b769d6")

	HelpStyle = lipgloss.NewStyle().Foreground(FeintColour).MarginLeft(1)

	// Bubbles

	SpinnerStyle = lipgloss.NewStyle().
			Foreground(HighlightColour).
			Bold(true)

	// Dashboard

	DashboardLabel = lipgloss.NewStyle().
			Background(SecondaryColour).
			Foreground(TextColour).
			Bold(true).
			MarginLeft(1).
			Width(11)

	TextStyle = lipgloss.NewStyle().Foreground(TextColour)
)
