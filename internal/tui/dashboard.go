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
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/purpleclay/dns53/internal/imds"
	"github.com/purpleclay/dns53/internal/r53"
	"github.com/purpleclay/dns53/internal/tui/component/errorpanel"
	"github.com/purpleclay/dns53/internal/tui/component/filteredlist"
	"github.com/purpleclay/dns53/internal/tui/component/header"
	"github.com/purpleclay/dns53/internal/tui/styles"
	"golang.org/x/term"
)

// DashboardModel defines the underlying model used for updating
// and rendering the dashboard TUI
type DashboardModel struct {
	opts DashboardOptions

	// bubbles used for capturing input from the user
	phz        list.Model
	loading    spinner.Model
	banner     header.Model
	errorPanel errorpanel.Model

	// data used to render final dashboard
	ec2       imds.Metadata
	connected *connection
}

// DashboardOptions defines all of the supported options when initialising
// the Dashboard model
type DashboardOptions struct {
	Metadata   imds.Metadata
	R53Client  *r53.Client
	Version    string
	PhzID      string
	DomainName string
}

type associationRequest struct {
	phz r53.PrivateHostedZone
}

type connection struct {
	phz r53.PrivateHostedZone
	dns string
}

type hostedZoneItem struct {
	name string
	id   string
}

func (i hostedZoneItem) Title() string       { return i.name }
func (i hostedZoneItem) Description() string { return i.id }
func (i hostedZoneItem) FilterValue() string { return i.name }

// Used to capture any error message that has been reported
type errMsg struct {
	reason string
	cause  error
}

func (e errMsg) Error() string {
	return e.cause.Error()
}

// Dashboard creates the initial model for the TUI
func Dashboard(opts DashboardOptions) *DashboardModel {
	width, _, _ := term.GetSize(int(os.Stdout.Fd()))

	m := &DashboardModel{opts: opts, ec2: opts.Metadata}
	m.phz = filteredlist.New([]list.Item{}, width, 20)

	m.loading = spinner.New()
	m.loading.Spinner = spinner.Dot
	m.loading.Style = styles.SpinnerStyle

	m.banner = header.New("dns53", "v0.1.0", "Dynamic DNS within Amazon Route53. Expose your EC2 quickly, easily and privately.", width)
	m.errorPanel = errorpanel.New()

	return m
}

// Init initialises the model ready for its first update and render
func (m DashboardModel) Init() tea.Cmd {
	return tea.Batch(
		m.loading.Tick,
		func() tea.Msg {
			if m.opts.PhzID != "" {
				return m.queryHostedZone()
			}

			return m.queryHostedZones()
		},
	)
}

// Update handles all IO operations
func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// TODO: fix this by batching up commands

	m.loading, cmd = m.loading.Update(msg)
	cmds = append(cmds, cmd)

	m.banner, cmd = m.banner.Update(msg)
	cmds = append(cmds, cmd)

	m.errorPanel, cmd = m.errorPanel.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case []r53.PrivateHostedZone:
		// PHZ have been successfully retrieved. Load them into the list
		items := make([]list.Item, 0, len(msg))
		for _, phz := range msg {
			items = append(items, hostedZoneItem{name: phz.Name, id: phz.ID})
		}
		m.phz.SetItems(items)
	case errMsg:
		m.errorPanel.RaiseError(msg.reason, msg.cause)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			if m.connected != nil && m.connected.dns != "connecting..." {
				record := r53.ResourceRecord{
					PhzID:    m.connected.phz.ID,
					Name:     m.connected.dns,
					Resource: m.ec2.IPv4,
				}

				m.opts.R53Client.DisassociateRecord(context.Background(), record)
			}

			return m, tea.Quit
		case "enter":
			i := m.phz.SelectedItem().(hostedZoneItem)

			cmds = append(cmds, func() tea.Msg {
				return associationRequest{
					phz: r53.PrivateHostedZone{ID: i.id, Name: i.name},
				}
			})
		}
	case associationRequest:
		m.connected = &connection{
			phz: msg.phz,
			dns: "connecting...",
		}

		cmds = append(cmds, m.initAssociation)
	case connection:
		m.connected.dns = msg.dns
	}

	if len(m.phz.Items()) > 0 && m.connected == nil {
		m.phz, cmd = m.phz.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View will attempt to render the current dashboard based on the model
func (m DashboardModel) View() string {
	viewport := ""

	if m.connected != nil {
		// Build each section of the dashboard in turn
		phzData := lipgloss.JoinVertical(
			lipgloss.Top,
			fmt.Sprintf("%s %s %s", styles.SecondaryLabel.Render("Name:"), styles.Spacing, m.connected.phz.Name),
			fmt.Sprintf("%s %s %s", styles.SecondaryLabel.Render("ID:"), styles.Spacing, m.connected.phz.ID),
		)

		phz := lipgloss.JoinHorizontal(
			lipgloss.Left,
			styles.PrimaryLabel.Copy().Render("PHZ:"),
			phzData,
		)

		ec2Data := lipgloss.JoinVertical(
			lipgloss.Top,
			fmt.Sprintf("%s %s %s", styles.SecondaryLabel.Render("IPv4:"), styles.Spacing, m.ec2.IPv4),
			fmt.Sprintf("%s %s %s", styles.SecondaryLabel.Render("Region:"), styles.Spacing, m.ec2.Region),
			fmt.Sprintf("%s %s %s", styles.SecondaryLabel.Render("VPC:"), styles.Spacing, m.ec2.VPC),
		)

		ec2 := lipgloss.JoinHorizontal(
			lipgloss.Left,
			styles.PrimaryLabel.Copy().Render("EC2:"),
			ec2Data,
		)

		dnsData := lipgloss.JoinVertical(
			lipgloss.Top,
			fmt.Sprintf("%s %s %s [A]", styles.SecondaryLabel.Render("Record:"), styles.Spacing, styles.Highlight.Render(m.connected.dns)),
			fmt.Sprintf("%s %s %s", styles.SecondaryLabel.Render("Status:"), styles.Spacing, styles.PendingStatus.Render("pending")),
		)

		dns := lipgloss.JoinHorizontal(
			lipgloss.Left,
			styles.PrimaryLabel.Copy().Render("DNS:"),
			dnsData,
		)

		dashboard := lipgloss.JoinVertical(
			lipgloss.Top,
			lipgloss.NewStyle().MarginBottom(2).Render(phz),
			lipgloss.NewStyle().MarginBottom(2).Render(ec2),
			dns,
		)

		viewport = lipgloss.NewStyle().Render(dashboard)
	} else {
		if len(m.phz.Items()) == 0 {
			viewport = lipgloss.JoinHorizontal(
				lipgloss.Left,
				m.loading.View(),
				styles.TextStyle.MarginBottom(1).Render(" Retrieving PHZs from AWS..."),
			)
		} else {
			viewport = lipgloss.JoinVertical(
				lipgloss.Top,
				styles.TextStyle.Render("Please select a PHZ from the list:"),
				m.phz.View(),
			)
		}
	}

	errorMsg := m.errorPanel.View()
	footer := ""

	var b strings.Builder

	dashboard := lipgloss.JoinVertical(lipgloss.Top,
		m.banner.View(),
		viewport,
		errorMsg,
		footer,
	)

	b.WriteString(lipgloss.NewStyle().Margin(1).Render(dashboard))
	return b.String()
}

func (m DashboardModel) queryHostedZones() tea.Msg {
	return []r53.PrivateHostedZone{
		{
			ID:   "AAXXCCVDDD11234432",
			Name: "test1",
		},
		{
			ID:   "KJFHFHD8847575HFHF",
			Name: "test2",
		},
		{
			ID:   "LLOOWEE6645453HHDH",
			Name: "test3",
		},
		{
			ID:   "AAXXCCVDDD11234432",
			Name: "test4",
		},
		{
			ID:   "KJFHFHD8847575HFHF",
			Name: "test5",
		},
		{
			ID:   "LLOOWEE6645453HHDH",
			Name: "test6",
		},
	}

	// phzs, err := m.opts.R53Client.ByVPC(context.Background(), m.ec2.VPC, m.ec2.Region)
	// if err != nil {
	// 	return errMsg{
	// 		reason: fmt.Sprintf("querying private hosted zones for VPC %s in region %s", m.ec2.VPC, m.ec2.Region),
	// 		cause:  err,
	// 	}
	// }

	// return phzs
}

func (m DashboardModel) queryHostedZone() tea.Msg {
	phz, err := m.opts.R53Client.ByID(context.Background(), m.opts.PhzID)
	if err != nil {
		return errMsg{
			reason: fmt.Sprintf("querying private hosted zone %s", m.opts.PhzID),
			cause:  err,
		}
	}

	return phz
}

func (m DashboardModel) initAssociation() tea.Msg {
	name := m.opts.DomainName
	if name == "" {
		name = fmt.Sprintf("%s.dns53.%s", strings.ReplaceAll(m.ec2.IPv4, ".", "-"), m.connected.phz.Name)
	} else {
		// Ensure root domain is appended as a suffix
		if !strings.HasSuffix(name, "."+m.connected.phz.Name) {
			name = fmt.Sprintf("%s.%s", name, m.connected.phz.Name)
		}
	}

	// record := r53.ResourceRecord{
	// 	PhzID:    m.connected.phz.ID,
	// 	Name:     name,
	// 	Resource: m.ec2.IPv4,
	// }

	// if err := m.opts.R53Client.AssociateRecord(context.Background(), record); err != nil {
	// 	return errMsg{
	// 		reason: fmt.Sprintf("associating EC2 with PHZ %s", m.connected.phz.Name),
	// 		cause:  err,
	// 	}
	// }

	return connection{dns: name, phz: m.connected.phz}
}
