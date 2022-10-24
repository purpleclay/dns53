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

package cmd

import (
	"context"
	"io"

	"github.com/purpleclay/dns53/internal/ec2"
	"github.com/purpleclay/dns53/internal/imds"
	"github.com/purpleclay/dns53/internal/r53"
	"github.com/purpleclay/dns53/internal/tui"
)

// Can be used to share internal state between cobra commands
type globalContext struct {
	context.Context

	out io.Writer

	// Only here to prevent the starting of bubbletea, useful when
	// testing just the cobra command
	skipTea bool

	// Capture the tea model for validation
	teaModel *tui.DashboardModel

	// Support overwriting the clients during within the
	// PersistentPreRunE hook for testing
	imdsClient *imds.Client
	ec2Client  *ec2.Client
	r53Client  *r53.Client
}

// Options provide a way of overriding parts of the context during
// execution, especially useful when configuring mocks during testing
type globalContextOption func(*globalContext)

func withIMDSClient(client *imds.Client) globalContextOption {
	return func(c *globalContext) {
		c.imdsClient = client
	}
}

func withR53Client(client *r53.Client) globalContextOption {
	return func(c *globalContext) {
		c.r53Client = client
	}
}

func withSkipTea() globalContextOption {
	return func(c *globalContext) {
		c.skipTea = true
	}
}
