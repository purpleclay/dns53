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

package tui

import (
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/purpleclay/dns53/internal/imds"
	"github.com/purpleclay/dns53/internal/r53"
	"github.com/purpleclay/dns53/internal/tui/component"
	"github.com/purpleclay/dns53/internal/tui/keymap"
	"github.com/purpleclay/dns53/internal/tui/message"
	"github.com/purpleclay/dns53/internal/tui/page"
)

var framed = lipgloss.NewStyle().Margin(1)

type pageIndex int

const (
	wizardPage pageIndex = iota
	dashboardPage
)

type Options struct {
	About        About
	R53Client    *r53.Client
	EC2Metadata  imds.Metadata
	DomainName   string
	HostedZoneID string
}

type About struct {
	Name             string
	Version          string
	ShortDescription string
}

type UI struct {
	header    component.Model
	pages     []page.Model
	pageIndex pageIndex
	footer    component.Model
}

func New(opts Options) UI {
	output := termenv.NewOutput(os.Stderr)

	pages := []page.Model{
		page.NewWizard(page.WizardOptions{
			Client:       opts.R53Client,
			Metadata:     opts.EC2Metadata,
			HostedZoneID: opts.HostedZoneID,
			DomainName:   opts.DomainName,
		}),
		page.NewDashboard(page.DashboardOptions{
			Client:     opts.R53Client,
			Metadata:   opts.EC2Metadata,
			DomainName: opts.DomainName,
			Output:     output,
		}),
	}

	index := wizardPage

	return UI{
		header:    component.NewHeader(opts.About.Name, opts.About.Version, opts.About.ShortDescription),
		pages:     pages,
		pageIndex: index,
		footer:    component.NewFooter(pages[index]),
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

		// Viewport need resizing in both axis within the available space
		pageX := msg.Width - x
		pageY := msg.Height - (y + u.header.Height() + u.footer.Height())

		for i := range u.pages {
			u.pages[i] = u.pages[i].Resize(pageX, pageY)
		}
	case message.R53ZoneSelectedMsg:
		u.pageIndex += dashboardPage

		u = u.refreshFooterKeyMap()
	case message.RefreshKeymapMsg:
		u = u.refreshFooterKeyMap()
	case tea.KeyMsg:
		if key.Matches(msg, keymap.Quit) {
			u.pages[u.pageIndex].Update(msg)
			return u, tea.Quit
		}
	}

	var currentPage tea.Model
	currentPage, cmd = u.pages[u.pageIndex].Update(msg)
	u.pages[u.pageIndex] = currentPage.(page.Model)
	cmds = append(cmds, cmd)

	return u, tea.Batch(cmds...)
}

func (u UI) View() string {
	view := lipgloss.JoinVertical(
		lipgloss.Left,
		u.header.View(),
		u.pages[u.pageIndex].View(),
		u.footer.View(),
	)

	return framed.Render(view)
}

func (u UI) margins() (int, int) {
	s := framed.Copy()
	return s.GetHorizontalFrameSize(), s.GetVerticalFrameSize()
}

func (u UI) refreshFooterKeyMap() UI {
	footer := u.footer.(*component.Footer)
	u.footer = footer.SetKeyMap(u.pages[u.pageIndex])
	return u
}
