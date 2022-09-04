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

package imdsstub

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
)

type Client struct {
	t testing.TB
}

func New(t testing.TB) *Client {
	return &Client{t: t}
}

func (c *Client) GetMetadata(ctx context.Context, params *imds.GetMetadataInput, optFns ...func(*imds.Options)) (*imds.GetMetadataOutput, error) {
	c.t.Helper()

	switch params.Path {
	case "":
		return wrapOutput("local-ipv4\nmac\nplacement-region\nplacement/availability-zone\ninstance-id\ntags/instance"), nil
	case "local-ipv4":
		return wrapOutput("10.0.1.100"), nil
	case "mac":
		return wrapOutput("06:e5:43:29:8f:08"), nil
	case "network/interfaces/macs/06:e5:43:29:8f:08/vpc-id":
		return wrapOutput("vpc-016d173db537793d1"), nil
	case "placement/region":
		return wrapOutput("us-east-1"), nil
	case "placement/availability-zone":
		return wrapOutput("us-east-1a"), nil
	case "instance-id":
		return wrapOutput("i-0decb1524582da041"), nil
	case "tags/instance":
		return wrapOutput("Name\nEnvironment"), nil
	case "tags/instance/Name":
		return wrapOutput("stub-ec2"), nil
	case "tags/instance/Environment":
		return wrapOutput("dev"), nil
	}

	return &imds.GetMetadataOutput{}, fmt.Errorf("unexpected instance category %s", params.Path)
}

func wrapOutput(value string) *imds.GetMetadataOutput {
	return &imds.GetMetadataOutput{
		Content: io.NopCloser(strings.NewReader(value)),
	}
}
