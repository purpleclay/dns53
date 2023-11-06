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

package component

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/purpleclay/dns53/internal/tui/keymap"
	"github.com/purpleclay/dns53/internal/tui/message"
	theme "github.com/purpleclay/lipgloss-theme"
	"gopkg.in/elazarl/goproxy.v1"
)

var (
	traceBanner    = faint.Copy().Margin(1, 0, 1, 0)
	traceSeparator = faint.Render("-")
)

var LogConnect goproxy.FuncHttpsHandler = func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	// Enable MITM (man in the middle) to capture the response output from HTTPS requests
	return goproxy.MitmConnect, host
}

type traceRequestMsg struct {
	method     string
	path       string
	statusCode int
}

func (t traceRequestMsg) String() string {
	statusCode := strconv.Itoa(t.statusCode)
	var status string
	if statusCode[0] == '2' {
		status = theme.H6.Render(statusCode)
	} else if statusCode[0] == '5' {
		status = theme.H1.Render(statusCode)
	} else {
		status = theme.H3.Render(statusCode)
	}

	return fmt.Sprintf("%s %s %-7s %s\n", status, traceSeparator, t.method, t.path)
}

type RequestTracer struct {
	proxyPort int
	server    *http.Server
	traces    chan traceRequestMsg
	viewport  viewport.Model
	tracingOn string
	trace     strings.Builder
}

func NewRequestTracer(port int) *RequestTracer {
	return &RequestTracer{
		proxyPort: port,
		traces:    make(chan traceRequestMsg),
		tracingOn: traceBanner.Render(fmt.Sprintf("proxing on :%d", port)),
		viewport:  viewport.New(0, 0),
	}
}

func (m *RequestTracer) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			proxy := goproxy.NewProxyHttpServer()
			proxy.OnRequest().HandleConnect(LogConnect)
			proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
				m.traces <- traceRequestMsg{
					method:     resp.Request.Method,
					path:       resp.Request.URL.Path,
					statusCode: resp.StatusCode,
				}
				return resp
			})
			m.server = &http.Server{
				Addr:    ":" + strconv.Itoa(m.proxyPort),
				Handler: proxy,
			}
			if err := m.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				return message.ErrorMsg{
					Reason: "failed to start reverse proxy for tracing requests",
					Cause:  err,
				}
			}

			return nil
		},
		m.waitForTrace())
}

func (m *RequestTracer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case traceRequestMsg:
		m.trace.WriteString(msg.String())
		m.viewport.GotoBottom()
		cmds = append(cmds, m.waitForTrace())
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymap.Quit):
			if m.server != nil {
				m.server.Shutdown(context.Background())
			}

			close(m.traces)
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *RequestTracer) View() string {
	m.viewport.SetContent(m.trace.String())

	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.tracingOn,
		m.viewport.View(),
	)
}

func (m *RequestTracer) Resize(width, height int) Model {
	m.viewport.Width = width
	m.viewport.Height = height - lipgloss.Height(m.tracingOn)
	return m
}

func (m *RequestTracer) Width() int {
	return m.viewport.Width
}

func (m *RequestTracer) Height() int {
	return m.viewport.Height + lipgloss.Height(m.tracingOn)
}

func (m *RequestTracer) waitForTrace() tea.Cmd {
	return func() tea.Msg {
		return <-m.traces
	}
}
