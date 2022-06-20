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

import "github.com/charmbracelet/lipgloss"

var (
	Purple = lipgloss.Color("#3a0ca3")
	Header = lipgloss.Color("#b5179e")
	Error  = lipgloss.Color("#d00000")
	White  = lipgloss.Color("#ffffff")
	Feint  = lipgloss.Color("#888888")

	Spinner = lipgloss.NewStyle().Foreground(lipgloss.Color("#7209b7"))

	Title = lipgloss.NewStyle().
		Background(Purple).
		Foreground(White).
		Bold(true).
		Padding(0, 1)

	Heading = lipgloss.NewStyle().Background(Header).Foreground(White)

	ErrorLabel = lipgloss.NewStyle().
			Background(Error).
			Foreground(White).
			Bold(true).
			Padding(0, 1).
			Render("Error:")

	MenuText = lipgloss.NewStyle().Foreground(Feint)
)
