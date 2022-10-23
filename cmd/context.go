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
)

// GlobalContext can be used to share state between cobra commands
type GlobalContext struct {
	context.Context

	out        io.Writer
	imdsClient *imds.Client
	ec2Client  *ec2.Client
	r53Client  *r53.Client
}

// GlobalContextOption is used to set an option within the GlobalContext at
// initialisation time. Initialisation of the context will occur during
// PersistentPreRunE within the root command
type GlobalContextOption func(*GlobalContext)

// WithOutput sets the output which, by default, is stdout
func WithOutput(output io.Writer) GlobalContextOption {
	return func(c *GlobalContext) {
		c.out = output
	}
}

// WithEC2Client ...
func WithEC2Client(client *ec2.Client) GlobalContextOption {
	return func(c *GlobalContext) {
		c.ec2Client = client
	}
}

// WithIMDSClient ...
func WithIMDSClient(client *imds.Client) GlobalContextOption {
	return func(c *GlobalContext) {
		c.imdsClient = client
	}
}

// WithR53Client ...
func WithR53Client(client *r53.Client) GlobalContextOption {
	return func(c *GlobalContext) {
		c.r53Client = client
	}
}
