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
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/purpleclay/dns53/pkg/ec2"
	"github.com/purpleclay/dns53/pkg/r53"
	"golang.org/x/term"
)

// TODO: store >
// metadata, phzid

// DashboardModel ...
type DashboardModel struct {
	cfg     aws.Config
	PHZ     list.Model
	EC2     ec2.Metadata
	Loading spinner.Model

	Err       error
	Connected *Connection
}

// Connection ...
type Connection struct {
	PHZ  string
	Name string
	DNS  string
}

type association struct{ dns string }

type initListItem struct {
	name        string
	description string
}

func (i initListItem) Title() string       { return i.name }
func (i initListItem) Description() string { return i.description }
func (i initListItem) FilterValue() string { return i.name }

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

// Dashboard ...
func Dashboard(cfg aws.Config) (*DashboardModel, error) {
	width, _, _ := term.GetSize(int(os.Stdout.Fd()))

	m := &DashboardModel{cfg: cfg}

	m.PHZ = list.New([]list.Item{}, list.NewDefaultDelegate(), width, 20)
	m.Loading = spinner.New()
	m.Loading.Spinner = spinner.Dot

	return m, nil
}

// Init ...
func (m DashboardModel) Init() tea.Cmd {
	return tea.Batch(
		m.Loading.Tick,
		func() tea.Msg {
			meta, err := ec2.InstanceMetadata(m.cfg)
			if err != nil {
				return errMsg{err}
			}

			return meta
		},
	)
}

// Update ...
func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.Loading, cmd = m.Loading.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case ec2.Metadata:
		// Dynamically retrieve R53 PHZ from AWS based on the EC2 metadata
		cmds = append(cmds, func() tea.Msg {
			time.Sleep(3 * time.Second)
			//phzs, _ := r53.ByVPC(m.cfg, meta.VPC, meta.Region)
			return []r53.PrivateHostedZone{
				{
					Name: "Zone1",
					ID:   "abcdef12345",
				},
				{
					Name: "Zone2",
					ID:   "abcdef54321",
				},
			}
		})
	case []r53.PrivateHostedZone:
		// PHZ have been successfully retrieved. Load them into the list
		items := make([]list.Item, 0, len(msg))
		for _, phz := range msg {
			items = append(items, initListItem{name: phz.ID, description: phz.Name})
		}
		m.PHZ.SetItems(items)
	case errMsg:
		m.Err = msg
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			i := m.PHZ.SelectedItem().(initListItem)
			m.Connected = &Connection{
				PHZ:  i.description,
				Name: i.name,
				DNS:  "connecting...",
			}

			cmds = append(cmds, func() tea.Msg {
				time.Sleep(3 * time.Second)

				// TODO: actually query AWS and generate the required association
				return association{"1111.1111.1111.1111 > 10.0.0.0"}
			})
		}
	case association:
		m.Connected.DNS = msg.dns
	}

	if len(m.PHZ.Items()) > 0 {
		m.PHZ, cmd = m.PHZ.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View ...
func (m DashboardModel) View() string {
	if m.Err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.Err)
	}

	var b strings.Builder

	if m.Connected != nil {
		b.WriteString("PHZ: " + m.Connected.PHZ)
		b.WriteRune('\n')
		b.WriteString("Name: " + m.Connected.Name)
		b.WriteRune('\n')
		b.WriteString("DNS: " + m.Connected.DNS)
		b.WriteRune('\n')
		return b.String()
	}

	// If PHZs have been retrieved, no longer render the spinner
	if len(m.PHZ.Items()) == 0 {
		str := fmt.Sprintf("%s Retrieving PHZs from AWS...\n\n", m.Loading.View())
		b.WriteString(str)
	} else {
		b.WriteString(m.PHZ.View())
	}

	return b.String()
}
