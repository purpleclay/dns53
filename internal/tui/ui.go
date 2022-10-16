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

	tea "github.com/charmbracelet/bubbletea"
	"github.com/purpleclay/dns53/internal/tui/component"
	"github.com/purpleclay/dns53/internal/tui/component/header"
	"github.com/purpleclay/dns53/internal/tui/theme"
)

// // Basic page management
// type page int

// const (
// 	wizardPage page = iota
// 	dashboardPage
// )

/*
AppDetails:{name, version, description} -> forwarded onto the header to render
*/

// UI ...
type UI struct {
	component.Model
	header header.Model
	//currentPage page
}

// New ...
func New(opts Options) UI {
	hdr := header.New(opts.About.Name, opts.About.Version, opts.About.ShortDescription)

	return UI{
		header: hdr,
	}
}

// TODO: Get rid of pointer receivers
func (u *UI) Resize(width, height int) {
	u.Model.Resize(width, height)
	wm, hm := u.getMargins()
	u.header.Resize(width-wm, height-hm)
	// TODO: resize the footer
}

func (ui *UI) getMargins() (int, int) {
	style := theme.AppStyle.Copy()
	wm := style.GetHorizontalFrameSize()
	hm := style.GetVerticalFrameSize()
	return wm, hm
}

// Init ...
func (u UI) Init() tea.Cmd {
	// create pages -> bundle their init commands
	return nil
}

// Update ...
func (u UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		fmt.Println(msg.Width)
		u.header.Resize(msg.Width, msg.Height)
		//u.Resize(msg.Width, msg.Height)
		// TODO: issue updates on the pages
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return u, tea.Quit
		}
	}

	return u, tea.Batch(cmds...)
}

// View ...
func (u UI) View() string {
	return u.header.View()
}
