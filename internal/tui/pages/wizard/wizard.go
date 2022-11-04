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

package wizard

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/purpleclay/dns53/internal/imds"
	"github.com/purpleclay/dns53/internal/r53"
	"github.com/purpleclay/dns53/internal/tui/components/errorpanel"
	"github.com/purpleclay/dns53/internal/tui/components/filteredlist"
	"github.com/purpleclay/dns53/internal/tui/message"
	"github.com/purpleclay/dns53/internal/tui/pages"
	"github.com/purpleclay/dns53/internal/tui/theme"
)

type zoneSelectionMsg struct {
	hostedZones []r53.PrivateHostedZone
}

type hostedZoneItem struct {
	name string
	id   string
}

func (i hostedZoneItem) Title() string       { return i.name }
func (i hostedZoneItem) Description() string { return i.id }
func (i hostedZoneItem) FilterValue() string { return i.name }

type Options struct {
	Client       *r53.Client
	Metadata     imds.Metadata
	HostedZoneID string
	DomainName   string
}

type Model struct {
	viewport    viewport.Model
	loading     spinner.Model
	selection   list.Model
	errorPanel  errorpanel.Model
	errorRaised bool
	options     Options
	styles      *Styles
	keymap      *KeyMap
}

func New(opts Options) Model {
	styles := DefaultStyles()

	loading := spinner.New()
	loading.Spinner = spinner.Dot
	loading.Style = styles.Spinner

	return Model{
		viewport:   viewport.New(0, 0),
		loading:    loading,
		selection:  filteredlist.New([]list.Item{}, 30, 20),
		errorPanel: errorpanel.New(),
		options:    opts,
		styles:     styles,
		keymap:     DefaultKeyMap(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.viewport.Init(),
		m.loading.Tick,
		func() tea.Msg {
			if m.options.HostedZoneID == "" {
				return m.queryHostedZones()
			}
			return m.queryHostedZone()
		},
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case zoneSelectionMsg:
		// PHZ have been successfully retrieved. Load them into the list
		items := make([]list.Item, 0, len(msg.hostedZones))
		for _, phz := range msg.hostedZones {
			items = append(items, hostedZoneItem{name: phz.Name, id: phz.ID})
		}
		m.selection.SetItems(items)

		// TODO: this should be made into a function
		cmds = append(cmds, func() tea.Msg {
			return message.RefreshKeymapMsg{}
		})
	case message.ErrorMsg:
		m.errorPanel = m.errorPanel.RaiseError(msg.Reason, msg.Cause)
		m.errorRaised = true
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selected := m.selection.SelectedItem().(hostedZoneItem)

			cmds = append(cmds, func() tea.Msg {
				return message.R53ZoneSelectedMsg{
					HostedZone: r53.PrivateHostedZone{ID: selected.id, Name: selected.name},
				}
			})
		case "/":
			fallthrough
		case "esc":
			// Refresh the keymap based on the list being populated
			cmds = append(cmds, func() tea.Msg {
				return message.RefreshKeymapMsg{}
			})
		}
	}

	m.loading, cmd = m.loading.Update(msg)
	cmds = append(cmds, cmd)

	m.selection, cmd = m.selection.Update(msg)
	cmds = append(cmds, cmd)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var page string
	if len(m.selection.Items()) == 0 {
		zoneLabel := "Zones"
		if m.options.HostedZoneID != "" {
			zoneLabel = "Zone"
		}

		page = lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.loading.View(),
			theme.TextStyle.MarginBottom(1).Render(fmt.Sprintf(" Retrieving Private Hosted %s from AWS...", zoneLabel)),
		)
	} else {
		page = lipgloss.JoinVertical(
			lipgloss.Top,
			theme.TextStyle.Render("Please select a Private Hosted Zone:"),
			m.selection.View(),
		)
	}

	var err string
	if m.errorRaised {
		err = "\n" + m.errorPanel.View()
	}

	view := lipgloss.JoinVertical(lipgloss.Top, page, err)

	m.viewport.SetContent(view)
	return m.viewport.View()
}

func (m Model) ShortHelp() []key.Binding {
	kb := make([]key.Binding, 0)

	kb = append(kb, m.keymap.Quit)

	// Respond to the selection being populated with items
	if len(m.selection.Items()) > 0 {
		if m.selection.FilterState() == list.Filtering {
			kb = append(kb, m.keymap.Enter, m.keymap.Escape)
		} else {
			kb = append(kb, m.keymap.UpDown)

			if m.selection.Paginator.TotalPages > 1 {
				kb = append(kb, m.keymap.LeftRight)
			}

			kb = append(kb, m.keymap.Enter, m.keymap.ForwardSlash)
		}
	}

	return kb
}

func (m Model) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

func (m Model) Resize(width, height int) pages.Model {
	m.viewport.Width = width
	m.viewport.Height = height

	m.errorPanel = m.errorPanel.Resize(width, height)
	return m
}

func (m Model) Width() int {
	return m.viewport.Width
}

func (m Model) Height() int {
	return m.viewport.Height
}

func (m Model) queryHostedZones() tea.Msg {
	phzs := []r53.PrivateHostedZone{
		{
			ID:   "Z0000000000000001",
			Name: "testing1",
		},
		{
			ID:   "Z0000000000000002",
			Name: "testing2",
		},
		{
			ID:   "Z0000000000000003",
			Name: "testing3",
		},
		{
			ID:   "Z0000000000000004",
			Name: "testing4",
		},
		{
			ID:   "Z0000000000000005",
			Name: "testing5",
		},
		{
			ID:   "Z0000000000000006",
			Name: "testing6",
		},
	}

	// metadata := m.options.Metadata
	// phzs, err := m.options.Client.ByVPC(context.Background(), metadata.VPC, metadata.Region)
	// if err != nil {
	// 	return message.ErrorMsg{
	// 		Reason: fmt.Sprintf("querying private hosted zones for VPC %s in region %s", metadata.VPC, metadata.Region),
	// 		Cause:  err,
	// 	}
	// }

	return zoneSelectionMsg{hostedZones: phzs}
}

func (m Model) queryHostedZone() tea.Msg {
	phz, err := m.options.Client.ByID(context.Background(), m.options.HostedZoneID)
	if err != nil {
		return message.ErrorMsg{
			Reason: fmt.Sprintf("querying private hosted zone %s", m.options.HostedZoneID),
			Cause:  err,
		}
	}

	return message.R53ZoneSelectedMsg{HostedZone: phz}
}
