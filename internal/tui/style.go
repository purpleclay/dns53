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

package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colours

	primary   = lipgloss.Color("#3A0CA3")
	secondary = lipgloss.Color("#8333B7")
	feint     = lipgloss.Color("#A2A0A0")
	red       = lipgloss.Color("#d00000")
	text      = lipgloss.Color("#FFFFFF")

	// Common

	br = "\n"

	// Header

	appNameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(primary).
			MarginLeft(1)

	helpStyle = lipgloss.NewStyle().Foreground(feint).MarginLeft(1)

	// Bubbles

	spinnerStyle = lipgloss.NewStyle().
			Foreground(secondary).
			Bold(true).
			MarginLeft(1)

	// Error

	errorLabelStyle = lipgloss.NewStyle().
			Background(red).
			Foreground(text).
			Bold(true).
			Padding(0, 1).
			Render("Error:")

	// Dashboard

	dashboardLabel = lipgloss.NewStyle().
			Background(secondary).
			Foreground(text).
			Bold(true).
			MarginLeft(1).
			Width(11)
)
