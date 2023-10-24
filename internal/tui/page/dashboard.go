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

package page

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/purpleclay/dns53/internal/imds"
	"github.com/purpleclay/dns53/internal/r53"
	"github.com/purpleclay/dns53/internal/tui/component"
	"github.com/purpleclay/dns53/internal/tui/keymap"
	"github.com/purpleclay/dns53/internal/tui/message"
	"github.com/purpleclay/lipgloss-theme"
)

var (
	spacing        = lipgloss.NewStyle().Faint(true).Render(strings.Repeat(".", 12))
	faint          = lipgloss.NewStyle().Faint(true)
	highlight      = lipgloss.NewStyle().Foreground(theme.S50).Bold(true)
	primaryLabel   = theme.H6.Copy().Padding(0, 3).MarginRight(8)
	secondaryLabel = theme.H2.Copy().PaddingLeft(2).Width(11)
	pending        = lipgloss.NewStyle().Foreground(lipgloss.Color("#e68a35")).Bold(true)
	active         = lipgloss.NewStyle().Foreground(lipgloss.Color("#26a621")).Bold(true)
)

// Common layout for the dashboard
const dashboardLine = "%s %s %s"

type r53AssociatedMsg struct{}

type DashboardOptions struct {
	Client     *r53.Client
	Metadata   imds.Metadata
	DomainName string
	Output     *termenv.Output
}

type Dashboard struct {
	viewport         viewport.Model
	options          DashboardOptions
	domainName       string
	clipboardStatus  string
	clipboardTimeout stopwatch.Model
	selected         r53.PrivateHostedZone
	connected        bool
	elapsed          stopwatch.Model
	errorPanel       component.ErrorPanel
	errorRaised      bool
}

func NewDashboard(opts DashboardOptions) *Dashboard {
	return &Dashboard{
		clipboardTimeout: stopwatch.NewWithInterval(time.Second * 2),
		viewport:         viewport.New(0, 0),
		options:          opts,
		elapsed:          stopwatch.New(),
		errorPanel:       component.NewErrorPanel(),
	}
}

func (m *Dashboard) Init() tea.Cmd {
	return nil
}

func (m *Dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case message.R53ZoneSelectedMsg:
		m.selected = msg.HostedZone
		m.domainName = m.resolveDomainName()

		cmds = append(cmds, m.initAssociation)
	case r53AssociatedMsg:
		m.connected = true
		cmds = append(cmds, m.elapsed.Start(), message.RefreshKeyMapCmd)
	case message.ErrorMsg:
		m.errorPanel = m.errorPanel.RaiseError(msg.Reason, msg.Cause)
		m.errorRaised = true
	case stopwatch.TickMsg:
		if msg.ID == m.clipboardTimeout.ID() {
			m.clipboardStatus = ""
			cmds = append(cmds, m.clipboardTimeout.Stop())
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymap.Quit):
			if m.connected {
				record := r53.ResourceRecord{
					PhzID:    m.selected.ID,
					Name:     m.domainName,
					Resource: m.options.Metadata.IPv4,
				}

				m.options.Client.DisassociateRecord(context.Background(), record)
			}
		case key.Matches(msg, keymap.Copy):
			m.options.Output.Copy(m.domainName)
			m.clipboardStatus = "(copied to clipboard)"

			// Trigger the timer to clear the clipboard message
			cmds = append(cmds, m.clipboardTimeout.Reset(), m.clipboardTimeout.Start())
		}
	}

	if m.connected {
		m.elapsed, cmd = m.elapsed.Update(msg)
		cmds = append(cmds, cmd)

		if m.clipboardTimeout.Running() {
			m.clipboardTimeout, cmd = m.clipboardTimeout.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Dashboard) View() string {
	phzData := lipgloss.JoinVertical(
		lipgloss.Top,
		fmt.Sprintf(dashboardLine, secondaryLabel.Render("Name:"), spacing, m.selected.Name),
		fmt.Sprintf(dashboardLine, secondaryLabel.Render("ID:"), spacing, m.selected.ID),
	)

	phz := lipgloss.JoinHorizontal(
		lipgloss.Left,
		primaryLabel.Render("PHZ:"),
		phzData,
	)

	ec2Data := lipgloss.JoinVertical(
		lipgloss.Top,
		fmt.Sprintf(dashboardLine, secondaryLabel.Render("IPv4:"), spacing, m.options.Metadata.IPv4),
		fmt.Sprintf(dashboardLine, secondaryLabel.Render("Region:"), spacing, m.options.Metadata.Region),
		fmt.Sprintf(dashboardLine, secondaryLabel.Render("VPC:"), spacing, m.options.Metadata.VPC),
	)

	ec2 := lipgloss.JoinHorizontal(
		lipgloss.Left,
		primaryLabel.Render("EC2:"),
		ec2Data,
	)

	status := pending.Render("pending")
	if m.connected {
		status = lipgloss.JoinHorizontal(lipgloss.Left, active.Render("active"), fmt.Sprintf(" (%s)", m.elapsed.View()))
	}

	dnsData := lipgloss.JoinVertical(
		lipgloss.Top,
		fmt.Sprintf(dashboardLine+" [A] %s",
			secondaryLabel.Render("Record:"),
			spacing,
			highlight.Render(m.domainName),
			faint.Render(m.clipboardStatus)),
		fmt.Sprintf(dashboardLine, secondaryLabel.Render("Status:"), spacing, status),
	)

	dns := lipgloss.JoinHorizontal(
		lipgloss.Left,
		primaryLabel.Render("DNS:"),
		dnsData,
	)

	dashboard := lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().MarginBottom(2).Render(phz),
		lipgloss.NewStyle().MarginBottom(2).Render(ec2),
		dns,
	)

	page := lipgloss.NewStyle().Render(dashboard)

	var err string
	if m.errorRaised {
		err = "\n\n" + m.errorPanel.View()
	}

	view := lipgloss.JoinVertical(lipgloss.Top, page, err)

	m.viewport.SetContent(view)
	return m.viewport.View()
}

func (m *Dashboard) ShortHelp() []key.Binding {
	bindings := []key.Binding{keymap.Quit}
	if m.connected {
		bindings = append(bindings, keymap.Copy)
	}

	return bindings
}

func (m *Dashboard) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

func (m *Dashboard) Resize(width, height int) Model {
	m.viewport.Width = width
	m.viewport.Height = height

	m.errorPanel = m.errorPanel.Resize(width, height).(component.ErrorPanel)
	return m
}

func (m *Dashboard) Width() int {
	return m.viewport.Width
}

func (m *Dashboard) Height() int {
	return m.viewport.Height
}

func (m *Dashboard) resolveDomainName() string {
	name := m.options.DomainName
	if name == "" {
		name = fmt.Sprintf("%s.dns53.%s", strings.ReplaceAll(m.options.Metadata.IPv4, ".", "-"), m.selected.Name)

		// If attaching to the dns53 domain, strip off the duplicate suffix
		if strings.Count(name, "dns53") > 1 {
			name = strings.TrimSuffix(name, ".dns53")
		}
	} else if !strings.HasSuffix(name, "."+m.selected.Name) {
		// Ensure root domain is appended as a suffix
		name = fmt.Sprintf("%s.%s", name, m.selected.Name)
	}

	return name
}

func (m *Dashboard) initAssociation() tea.Msg {
	record := r53.ResourceRecord{
		PhzID:    m.selected.ID,
		Name:     m.domainName,
		Resource: m.options.Metadata.IPv4,
	}

	if err := m.options.Client.AssociateRecord(context.Background(), record); err != nil {
		return message.ErrorMsg{
			Reason: fmt.Sprintf("associating EC2 with PHZ %s", m.selected.Name),
			Cause:  err,
		}
	}

	return r53AssociatedMsg{}
}
