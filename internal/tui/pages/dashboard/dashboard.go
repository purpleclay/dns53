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

package dashboard

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/purpleclay/dns53/internal/imds"
	"github.com/purpleclay/dns53/internal/r53"
	"github.com/purpleclay/dns53/internal/tui/components/errorpanel"
	"github.com/purpleclay/dns53/internal/tui/message"
	"github.com/purpleclay/dns53/internal/tui/pages"
)

type r53AssociatedMsg struct{}

// Options ...
type Options struct {
	Client     *r53.Client
	Metadata   imds.Metadata
	DomainName string
}

type Model struct {
	viewport    viewport.Model
	options     Options
	domainName  string
	selected    r53.PrivateHostedZone
	connected   bool
	elapsed     stopwatch.Model
	errorPanel  errorpanel.Model
	errorRaised bool
	styles      *Styles
}

func New(opts Options) Model {
	return Model{
		viewport:   viewport.New(0, 0),
		options:    opts,
		elapsed:    stopwatch.New(),
		errorPanel: errorpanel.New(),
		styles:     DefaultStyles(),
	}
}

// Init ...
func (m Model) Init() tea.Cmd {
	return nil
}

// Update ...
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		cmds = append(cmds, m.elapsed.Start())
	case message.ErrorMsg:
		m.errorPanel = m.errorPanel.RaiseError(msg.Reason, msg.Cause)
		m.errorRaised = true
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			if m.connected {
				record := r53.ResourceRecord{
					PhzID:    m.selected.ID,
					Name:     m.domainName,
					Resource: m.options.Metadata.IPv4,
				}

				m.options.Client.DisassociateRecord(context.Background(), record)
			}
		}
	}

	if m.connected {
		m.elapsed, cmd = m.elapsed.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View ...
func (m Model) View() string {
	phzData := lipgloss.JoinVertical(
		lipgloss.Top,
		fmt.Sprintf("%s %s %s", m.styles.SecondaryLabel.Render("Name:"), m.styles.Spacing, m.selected.Name),
		fmt.Sprintf("%s %s %s", m.styles.SecondaryLabel.Render("ID:"), m.styles.Spacing, m.selected.ID),
	)

	phz := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.styles.PrimaryLabel.Copy().Render("PHZ:"),
		phzData,
	)

	ec2Data := lipgloss.JoinVertical(
		lipgloss.Top,
		fmt.Sprintf("%s %s %s", m.styles.SecondaryLabel.Render("IPv4:"), m.styles.Spacing, m.options.Metadata.IPv4),
		fmt.Sprintf("%s %s %s", m.styles.SecondaryLabel.Render("Region:"), m.styles.Spacing, m.options.Metadata.Region),
		fmt.Sprintf("%s %s %s", m.styles.SecondaryLabel.Render("VPC:"), m.styles.Spacing, m.options.Metadata.VPC),
	)

	ec2 := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.styles.PrimaryLabel.Copy().Render("EC2:"),
		ec2Data,
	)

	status := m.styles.PendingStatus.Render("pending")
	if m.connected {
		status = lipgloss.JoinHorizontal(lipgloss.Left, m.styles.ActiveStatus.Render("active"), fmt.Sprintf(" (%s)", m.elapsed.View()))
	}

	dnsData := lipgloss.JoinVertical(
		lipgloss.Top,
		fmt.Sprintf("%s %s %s [A]", m.styles.SecondaryLabel.Render("Record:"), m.styles.Spacing, m.styles.Highlight.Render(m.domainName)),
		fmt.Sprintf("%s %s %s", m.styles.SecondaryLabel.Render("Status:"), m.styles.Spacing, status),
	)

	dns := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.styles.PrimaryLabel.Copy().Render("DNS:"),
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

// ShortHelp ...
func (m Model) ShortHelp() []key.Binding {
	return []key.Binding{}
}

// FullHelp ...
func (m Model) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

// Resize ...
func (m Model) Resize(width, height int) pages.Model {
	m.viewport.Width = width
	m.viewport.Height = height

	m.errorPanel = m.errorPanel.Resize(width, height)
	return m
}

// Width ...
func (m Model) Width() int {
	return m.viewport.Width
}

// Height ...
func (m Model) Height() int {
	return m.viewport.Height
}

func (m Model) resolveDomainName() string {
	name := m.options.DomainName
	if name == "" {
		name = fmt.Sprintf("%s.dns53.%s", strings.ReplaceAll(m.options.Metadata.IPv4, ".", "-"), m.selected.Name)

		// If attaching to the dns53 domain, strip off the duplicate suffix
		if strings.Count(name, "dns53") > 1 {
			name = strings.TrimSuffix(name, ".dns53")
		}
	} else {
		// Ensure root domain is appended as a suffix
		if !strings.HasSuffix(name, "."+m.selected.Name) {
			name = fmt.Sprintf("%s.%s", name, m.selected.Name)
		}
	}

	return name
}

func (m Model) initAssociation() tea.Msg {
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
