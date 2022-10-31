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
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/purpleclay/dns53/internal/tui/components"
	"github.com/purpleclay/dns53/internal/tui/components/footer"
	"github.com/purpleclay/dns53/internal/tui/components/header"
	"github.com/purpleclay/dns53/internal/tui/message"
	"github.com/purpleclay/dns53/internal/tui/pages"
	"github.com/purpleclay/dns53/internal/tui/pages/dashboard"
	"github.com/purpleclay/dns53/internal/tui/pages/wizard"
	"github.com/purpleclay/dns53/internal/tui/theme"
)

// Basic page management
type page int

const (
	wizardPage page = iota
	dashboardPage
)

// UI ...
type UI struct {
	header      components.Model
	pages       []pages.Model
	currentPage page
	footer      components.Model
}

// New ...
func New(opts Options) UI {
	// TODO: figure out how to bind keymaps to footer
	return UI{
		header: header.New(opts.About.Name, opts.About.Version, opts.About.ShortDescription),
		pages: []pages.Model{
			wizard.New(wizard.Options{
				Client:       opts.R53Client,
				Metadata:     opts.EC2Metadata,
				HostedZoneID: opts.HostedZoneID,
				DomainName:   opts.DomainName,
			}),
			dashboard.New(dashboard.Options{
				Client:     opts.R53Client,
				Metadata:   opts.EC2Metadata,
				DomainName: opts.DomainName,
			}),
		},
		currentPage: wizardPage,
		footer:      footer.New(),
	}
}

// Init ...
func (u UI) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	cmds = append(cmds, u.header.Init(), u.footer.Init())

	for i := range u.pages {
		cmds = append(cmds, u.pages[i].Init())
	}

	return tea.Batch(cmds...)
}

// Update ...
func (u UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		x, y := u.margins()

		// Only need to resize the width of both the header and footer
		u.header = u.header.Resize(msg.Width-x, u.header.Height())
		u.footer = u.footer.Resize(msg.Width-x, u.footer.Height())

		// Pages need resizing in both axis within the available space
		pageX := msg.Width - x
		pageY := msg.Height - (y + u.header.Height() + u.footer.Height())

		// TODO: manage a central viewport here, rather than on each page
		for i := range u.pages {
			u.pages[i] = u.pages[i].Resize(pageX, pageY)
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return u, tea.Quit
		}
	case message.R53ZoneSelectedMsg:
		u.currentPage += dashboardPage
	}

	var page tea.Model
	page, cmd = u.pages[u.currentPage].Update(msg)
	u.pages[u.currentPage] = page.(pages.Model)
	cmds = append(cmds, cmd)

	// TODO: if a Quit message has been sent (tea.Quit as the page will have handled the key msg)
	return u, tea.Batch(cmds...)
}

// View ...
func (u UI) View() string {
	view := lipgloss.JoinVertical(
		lipgloss.Left,
		u.header.View(),
		u.pages[u.currentPage].View(),
		u.footer.View(),
	)

	return theme.AppStyle.Render(view)
}

func (u UI) margins() (int, int) {
	s := theme.AppStyle.Copy()
	return s.GetHorizontalFrameSize(), s.GetVerticalFrameSize()
}
