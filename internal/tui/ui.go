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
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/purpleclay/dns53/internal/tui/components"
	"github.com/purpleclay/dns53/internal/tui/components/footer"
	"github.com/purpleclay/dns53/internal/tui/components/header"
	"github.com/purpleclay/dns53/internal/tui/keymap"
	"github.com/purpleclay/dns53/internal/tui/message"
	"github.com/purpleclay/dns53/internal/tui/pages"
	"github.com/purpleclay/dns53/internal/tui/pages/dashboard"
	"github.com/purpleclay/dns53/internal/tui/pages/wizard"
	"github.com/purpleclay/dns53/internal/tui/theme"
)

type page int

const (
	wizardPage page = iota
	dashboardPage
)

type UI struct {
	header      components.Model
	pages       []pages.Model
	currentPage page
	footer      components.Model
}

func New(opts Options) UI {
	pages := []pages.Model{
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
	}

	currentPage := wizardPage

	return UI{
		header:      header.New(opts.About.Name, opts.About.Version, opts.About.ShortDescription),
		pages:       pages,
		currentPage: currentPage,
		footer:      footer.New(pages[currentPage]),
	}
}

func (u UI) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	cmds = append(cmds, u.header.Init(), u.footer.Init())

	for i := range u.pages {
		cmds = append(cmds, u.pages[i].Init())
	}

	return tea.Batch(cmds...)
}

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
	case message.R53ZoneSelectedMsg:
		u.currentPage += dashboardPage

		u = u.refreshFooterKeyMap()
	case message.RefreshKeymapMsg:
		u = u.refreshFooterKeyMap()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymap.Quit):
			u.pages[u.currentPage].Update(msg)
			return u, tea.Quit
		}
	}

	var page tea.Model
	page, cmd = u.pages[u.currentPage].Update(msg)
	u.pages[u.currentPage] = page.(pages.Model)
	cmds = append(cmds, cmd)

	return u, tea.Batch(cmds...)
}

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

func (u UI) refreshFooterKeyMap() UI {
	footer := u.footer.(footer.Model)
	u.footer = footer.SetKeyMap(u.pages[u.currentPage])
	return u
}
