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

var defaultMetadata = map[string]string{
	"":           "local-ipv4\nmac\nplacement-region\nplacement/availability-zone\ninstance-id\ntags/instance",
	"local-ipv4": "10.0.1.100",
	"mac":        "06:e5:43:29:8f:08",
	"network/interfaces/macs/06:e5:43:29:8f:08/vpc-id": "vpc-016d173db537793d1",
	"placement/region":            "us-east-1",
	"placement/availability-zone": "us-east-1a",
	"instance-id":                 "i-0decb1524582da041",
	"tags/instance":               "Name\nEnvironment",
	"tags/instance/Name":          "stub-ec2",
	"tags/instance/Environment":   "dev",
}

type Client struct {
	t        testing.TB
	metadata map[string]string
}

func New(t testing.TB) *Client {
	return &Client{t: t, metadata: defaultMetadata}
}

func NewWithoutTags(t testing.TB) *Client {
	// Remove all traces of tags from the default metadata
	noTags := defaultMetadata
	noTags[""] = "local-ipv4\nmac\nplacement-region\nplacement/availability-zone\ninstance-id"
	delete(noTags, "tags/instance")
	delete(noTags, "tags/instance/Name")
	delete(noTags, "tags/instance/Environment")

	return &Client{t: t, metadata: noTags}
}

func (c *Client) GetMetadata(ctx context.Context, params *imds.GetMetadataInput, optFns ...func(*imds.Options)) (*imds.GetMetadataOutput, error) {
	c.t.Helper()

	if category, ok := c.metadata[params.Path]; ok {
		return wrapOutput(category), nil
	}

	return &imds.GetMetadataOutput{}, fmt.Errorf("unexpected instance category %s", params.Path)
}

func wrapOutput(value string) *imds.GetMetadataOutput {
	return &imds.GetMetadataOutput{
		Content: io.NopCloser(strings.NewReader(value)),
	}
}