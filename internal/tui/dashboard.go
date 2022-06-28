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
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/purpleclay/dns53/pkg/ec2"
	"github.com/purpleclay/dns53/pkg/r53"
	"golang.org/x/term"
)

type DashboardModel struct {
	cfg       aws.Config
	version   string
	phz       list.Model
	ec2       ec2.Metadata
	loading   spinner.Model
	err       error
	connected *connection
}

type connection struct {
	phz  string
	name string
	dns  string
}

type association struct {
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
func Dashboard(cfg aws.Config, version string) (*DashboardModel, error) {
	width, _, _ := term.GetSize(int(os.Stdout.Fd()))

	m := &DashboardModel{cfg: cfg, version: version}

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
			meta, err := ec2.InstanceMetadata(m.cfg)
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
	case ec2.Metadata:
		m.ec2 = msg

		// Dynamically retrieve R53 phz from AWS based on the EC2 metadata
		cmds = append(cmds, func() tea.Msg {
			phzs, err := r53.ByVPC(m.cfg, msg.VPC, msg.Region)
			if err != nil {
				return errMsg{err}
			}

			return phzs
		})
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
				r53.DisassociateRecord(m.cfg, m.connected.phz, m.connected.dns, m.ec2.IPv4)
			}

			return m, tea.Quit
		case "enter":
			i := m.phz.SelectedItem().(hostedZoneItem)
			m.connected = &connection{
				phz:  i.name,
				name: i.description,
				dns:  "connecting...",
			}

			cmds = append(cmds, func() tea.Msg {
				name := fmt.Sprintf("%s.dns53.%s", strings.ReplaceAll(m.ec2.IPv4, ".", "-"), m.connected.name)

				if err := r53.AssociateRecord(m.cfg, m.connected.phz, name, m.ec2.IPv4); err != nil {
					return errMsg{err}
				}
				return association{name}
			})
		}
	case association:
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
	version := titleItemStyle.Padding(0, 2).Render(m.version)
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
			fmt.Sprintf("%s [%s]", m.connected.name, m.connected.phz))

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
