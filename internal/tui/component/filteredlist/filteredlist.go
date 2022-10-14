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

package filteredlist

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/purpleclay/dns53/internal/tui/styles"
)

// New ...
func New(tems []list.Item, width, height int) list.Model {
	delegate := list.NewDefaultDelegate()

	// Override the colours within the existing styles
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		BorderForeground(styles.HighlightColour).
		Foreground(styles.HighlightColour).
		Bold(true)

	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(styles.HighlightFeintColour)

	delegate.Styles.DimmedDesc = delegate.Styles.DimmedDesc.
		Foreground(styles.FeintColour)

	delegate.Styles.FilterMatch = lipgloss.NewStyle().
		Underline(true).
		Bold(true)

	filteredList := list.New([]list.Item{}, delegate, width, height)

	// TODO: look into setting a spinner, removes the need for the first message ("Retrieving PHZs...")
	// TODO: provide a custom key map, this will be used when displaying the footer

	// Override the colours within the existing styles
	filteredList.Styles.FilterPrompt = filteredList.Styles.FilterPrompt.
		Foreground(styles.HighlightColour)

	filteredList.Styles.FilterCursor = filteredList.Styles.FilterCursor.
		Foreground(styles.HighlightColour)

	filteredList.Styles.StatusBarFilterCount = filteredList.Styles.StatusBarFilterCount.
		Foreground(styles.FeintColour)

	filteredList.SetShowTitle(false)
	filteredList.SetShowHelp(false)
	filteredList.DisableQuitKeybindings()

	return filteredList
}
