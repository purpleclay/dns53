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

package component

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/purpleclay/dns53/internal/tui/keymap"
	theme "github.com/purpleclay/lipgloss-theme"
)

func NewFilteredList(items []list.Item, width, height int) list.Model {
	delegate := list.NewDefaultDelegate()

	// Override the colours within the existing styles
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		BorderForeground(theme.S100).
		Foreground(theme.S100).
		Bold(true)

	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		BorderForeground(theme.S100).
		Foreground(theme.S50).
		Faint(true)

	delegate.Styles.FilterMatch = lipgloss.NewStyle().
		Underline(true).
		Bold(true)

	filteredList := list.New(items, delegate, width, height)

	// Override the colours within the existing styles
	filteredList.Styles.FilterPrompt = filteredList.Styles.FilterPrompt.
		Foreground(theme.S100)

	filteredList.Styles.FilterCursor = filteredList.Styles.FilterCursor.
		Foreground(theme.S100)

	filteredList.SetShowTitle(false)
	filteredList.SetShowHelp(false)
	filteredList.DisableQuitKeybindings()

	// Override key bindings to force expected behaviour
	filteredList.KeyMap.GoToEnd.SetEnabled(false)
	filteredList.KeyMap.GoToStart.SetEnabled(false)

	filteredList.KeyMap.CursorUp = keymap.Up
	filteredList.KeyMap.CursorDown = keymap.Down
	filteredList.KeyMap.NextPage = keymap.Right
	filteredList.KeyMap.PrevPage = keymap.Left

	return filteredList
}
