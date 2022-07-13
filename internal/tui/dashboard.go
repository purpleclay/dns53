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
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/purpleclay/dns53/pkg/imds"
	"github.com/purpleclay/dns53/pkg/r53"
	"golang.org/x/term"
)

// DashboardModel defines the underlying model used for updating
// and rendering the dashboard TUI
type DashboardModel struct {
	opts DashboardOptions

	// bubbles used for capturing input from the user
	phz     list.Model
	loading spinner.Model

	// data used to render final dashboard
	ec2       imds.Metadata
	connected *connection
	err       error
}

// DashboardOptions defines all of the supported options when initialising
// the Dashboard model
type DashboardOptions struct {
	IMDSClient *imds.Client
	R53Client  *r53.Client
	Version    string
	PhzID      string
	DNSName    string
}

type associationRequest struct {
	phz r53.PrivateHostedZone
}

type connection struct {
	phz r53.PrivateHostedZone
	dns string
}

type hostedZoneItem struct {
	name        string
	description string
}

func (i hostedZoneItem) Title() string       { return i.name }
func (i hostedZoneItem) Description() string { return i.description }
func (i hostedZoneItem) FilterValue() string { return i.description }

// Used to capture any error message that has been reported
type errMsg struct {
	err error
}

func (e errMsg) Error() string {
	return e.err.Error()
}

// Dashboard creates the initial model for the TUI
func Dashboard(opts DashboardOptions) (*DashboardModel, error) {
	width, _, _ := term.GetSize(int(os.Stdout.Fd()))

	m := &DashboardModel{opts: opts}

	m.phz = list.New([]list.Item{}, list.NewDefaultDelegate(), width, 20)
	m.phz.Styles.HelpStyle = listHelpStyle
	m.phz.SetShowFilter(false)
	m.phz.SetShowTitle(false)
	m.phz.DisableQuitKeybindings()

	m.loading = spinner.New()
	m.loading.Spinner = spinner.Dot
	m.loading.Style = spinnerStyle

	return m, nil
}

// Init initialises the model ready for its first update and render
func (m DashboardModel) Init() tea.Cmd {
	return tea.Batch(
		m.loading.Tick,
		func() tea.Msg {
			meta, err := m.opts.IMDSClient.InstanceMetadata(context.TODO())
			if err != nil {
				return errMsg{err}
			}

			return meta
		},
	)
}

// Update handles all IO operations
func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.loading, cmd = m.loading.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case imds.Metadata:
		m.ec2 = msg

		// If the PHZ is already known by this point, attempt an association
		if m.opts.PhzID != "" {
			cmds = append(cmds, m.queryHostedZone)
		} else {
			cmds = append(cmds, m.queryHostedZones)
		}
	case []r53.PrivateHostedZone:
		// PHZ have been successfully retrieved. Load them into the list
		items := make([]list.Item, 0, len(msg))
		for _, phz := range msg {
			items = append(items, hostedZoneItem{name: phz.ID, description: phz.Name})
		}
		m.phz.SetItems(items)
	case errMsg:
		m.err = msg
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			if m.connected != nil && m.connected.dns != "connecting..." {
				record := r53.ResourceRecord{
					PhzID:    m.connected.phz.ID,
					Name:     m.connected.dns,
					Resource: m.ec2.IPv4,
				}

				m.opts.R53Client.DisassociateRecord(context.TODO(), record)
			}

			return m, tea.Quit
		case "enter":
			i := m.phz.SelectedItem().(hostedZoneItem)

			cmds = append(cmds, func() tea.Msg {
				return associationRequest{
					phz: r53.PrivateHostedZone{ID: i.name, Name: i.description},
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
	width, _, _ := term.GetSize(int(os.Stdout.Fd()))
	rw := width - 5

	var b strings.Builder

	// Render the title bar
	name := titleItemStyle.Padding(0, 3).Render("dns53")
	version := titleItemStyle.Padding(0, 2).Render(m.opts.Version)
	menu := titleMenuStyle.Copy().
		Width(rw - lipgloss.Width(name) - lipgloss.Width(version)).
		PaddingLeft(2).
		Render("quit (ctrl+c)")

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		name,
		menu,
		version,
	)

	b.WriteString(titleBarStyle.Width(rw).Render(bar))
	b.WriteString(br)

	if m.connected != nil {
		phzLabel := dashboardLabel.Padding(0, 2).Render("PHZ:")
		ec2MetaLabel := dashboardLabel.Padding(0, 2).Render("EC2:")
		dnsLabel := dashboardLabel.Padding(0, 2).Render("DNS:")

		lbl := lipgloss.NewStyle().Width(20)

		phz := lipgloss.JoinHorizontal(lipgloss.Top,
			lbl.Render(phzLabel),
			fmt.Sprintf("%s [%s]", m.connected.phz.ID, m.connected.phz.Name))

		ec2Meta := lipgloss.JoinHorizontal(lipgloss.Top,
			lbl.Render(ec2MetaLabel),
			fmt.Sprintf("%s   :>   %s   :>   %s", m.ec2.IPv4, m.ec2.Region, m.ec2.VPC))

		dns := lipgloss.JoinHorizontal(lipgloss.Top,
			lbl.Render(dnsLabel),
			fmt.Sprintf("%s   ~>   localhost   [A]", m.connected.dns))

		dashboard := lipgloss.JoinVertical(lipgloss.Top,
			phz,
			br,
			ec2Meta,
			br,
			dns)

		b.WriteString(lipgloss.NewStyle().MarginTop(3).Render(dashboard))
	} else {
		// If phzs have been retrieved, no longer render the spinner
		if len(m.phz.Items()) == 0 {
			str := fmt.Sprintf("%s Retrieving phzs from AWS...\n\n", m.loading.View())
			b.WriteString(str)
		} else {
			b.WriteString(m.phz.View())
		}
	}

	if m.err != nil {
		errorPanelStyle := lipgloss.NewStyle().MarginLeft(1).Width(rw)

		errorPanel := lipgloss.JoinVertical(lipgloss.Top,
			fmt.Sprintf("\n\n%s", errorLabelStyle),
			fmt.Sprintf("\n%s\n", m.err.Error()),
		)

		b.WriteString(errorPanelStyle.Render(errorPanel))
	}

	return b.String()
}

func (m DashboardModel) queryHostedZones() tea.Msg {
	phzs, err := m.opts.R53Client.ByVPC(context.TODO(), m.ec2.VPC, m.ec2.Region)
	if err != nil {
		return errMsg{err}
	}

	return phzs
}

func (m DashboardModel) queryHostedZone() tea.Msg {
	phz, err := m.opts.R53Client.ByID(context.TODO(), m.opts.PhzID)
	if err != nil {
		return errMsg{err}
	}

	return associationRequest{phz: phz}
}

func (m DashboardModel) initAssociation() tea.Msg {
	// Sanitise the IPv4 within the EC2 Metadata Object
	ipv4 := m.ec2.IPv4
	m.ec2.IPv4 = strings.ReplaceAll(m.ec2.IPv4, ".", "-")

	var name string
	if m.opts.DNSName != "" {
		name = appendDNSSuffix(m.opts.DNSName, m.connected.phz.Name)

		// Check if the provided name contains a template
		if strings.Contains(name, "{{") {
			tmpl, err := template.New("dns").Parse(name)
			if err != nil {
				return errMsg{err}
			}

			var out bytes.Buffer
			if err := tmpl.Execute(&out, m.ec2); err != nil {
				return errMsg{err}
			}

			name = out.String()
		}
	} else {
		// By default include the dns53 suffix
		name = fmt.Sprintf("%s.dns53.%s", m.ec2.IPv4, m.connected.phz.Name)
	}

	record := r53.ResourceRecord{
		PhzID:    m.connected.phz.ID,
		Name:     name,
		Resource: ipv4,
	}

	if err := m.opts.R53Client.AssociateRecord(context.TODO(), record); err != nil {
		return errMsg{err}
	}

	return connection{dns: name, phz: m.connected.phz}
}

func appendDNSSuffix(dns, domain string) string {
	if strings.HasSuffix(dns, "dns53."+domain) {
		return dns
	}

	// If suffix has only been partially set, trim it
	dns = strings.TrimSuffix(dns, ".dns53")
	dns = strings.TrimSuffix(dns, "."+domain)

	return fmt.Sprintf("%s.dns53.%s", dns, domain)
}
