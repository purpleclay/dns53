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
	PrimaryColour   = lipgloss.Color("#3A0CA3")
	SecondaryColour = lipgloss.Color("#8333B7")
	FeintColour     = lipgloss.Color("#A2A0A0")
	RedColour       = lipgloss.Color("#d00000")
	TextColour      = lipgloss.Color("#FFFFFF")

	// Common

	//br = "\n"

	// Header

	AppNameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(PrimaryColour).
			MarginLeft(1)

	HelpStyle = lipgloss.NewStyle().Foreground(FeintColour).MarginLeft(1)

	// Bubbles

	SpinnerStyle = lipgloss.NewStyle().
			Foreground(SecondaryColour).
			Bold(true).
			MarginLeft(1)

	// Error

	ErrorLabelStyle = lipgloss.NewStyle().
			Background(RedColour).
			Foreground(TextColour).
			Bold(true).
			Padding(0, 1).
			Render("Error:")

	// Dashboard

	DashboardLabel = lipgloss.NewStyle().
			Background(SecondaryColour).
			Foreground(TextColour).
			Bold(true).
			MarginLeft(1).
			Width(11)
)
