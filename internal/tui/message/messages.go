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

package message

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/purpleclay/dns53/internal/r53"
)

// R53ZoneSelectedMsg signals that a R53 Private Hosted Zone has
// been selected and an association can now be made between it and
// the EC2
type R53ZoneSelectedMsg struct {
	HostedZone r53.PrivateHostedZone
}

// ErrorMsg should be sent to notify a user of an unrecoverable error
type ErrorMsg struct {
	Reason string
	Cause  error
}

// RefreshKeymapMsg should be sent when the keymap associated with the
// help view needs updating
type RefreshKeymapMsg struct{}

// RefreshKeyMapCmd provides a utility method that wraps the RefreshKeymapMsg
// message for use within a update method
func RefreshKeyMapCmd() tea.Msg {
	return RefreshKeymapMsg{}
}
