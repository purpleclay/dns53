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

package imds

import (
	"context"
	"fmt"
	"io"
	"strings"

	awsimds "github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
)

const (
	pathIPv4            = "local-ipv4"
	pathMacAddress      = "mac"
	pathPlacementRegion = "placement/region"
	pathPlacementAZ     = "placement/availability-zone"
	pathInstanceID      = "instance-id"
	pathTagsInstance    = "tags/instance"
)

// ClientAPI defines the API for interacting with the Amazon
// EC2 Instance Metadata Service (IMDS)
type ClientAPI interface {
	// GetMetadata uses the path provided to request information from the Amazon
	// EC2 Instance Metadata Service
	GetMetadata(ctx context.Context, params *awsimds.GetMetadataInput, optFns ...func(*awsimds.Options)) (*awsimds.GetMetadataOutput, error)
}

// Client defines the client for interacting with the Amazon EC2 Instance
// Metadata Service (IMDS)
type Client struct {
	api ClientAPI
}

// Metadata contains metadata associated with an EC2 instance
type Metadata struct {
	// IPv4 is the private IPv4 address of the launched instance
	IPv4 string

	// Region of where the EC2 instance was launched
	Region string

	// VPC ID of where the EC2 instance was launched
	VPC string

	// AZ is the availability zone where the instance was launched
	AZ string

	// InstanceID is the unique ID of this instance
	InstanceID string

	// Name associated with the EC2 instance. This will be blank unless
	// tags have been enabled within IMDS for this EC2 instance
	Name string

	// Tags contains ...
	Tags map[string]string
}

// NewFromAPI returns a new client from the provided IMDS API implementation
func NewFromAPI(api ClientAPI) *Client {
	return &Client{api: api}
}

// InstanceMetadata attempts to retrieve useful metadata associated with
// the current EC2 instance by querying IMDS
func (c *Client) InstanceMetadata(ctx context.Context) (Metadata, error) {
	if err := checkRoot(ctx, c.api); err != nil {
		return Metadata{}, err
	}

	md := Metadata{}
	md.AZ, _ = get(ctx, c.api, pathPlacementAZ)
	md.InstanceID, _ = get(ctx, c.api, pathInstanceID)
	md.IPv4, _ = get(ctx, c.api, pathIPv4)
	md.Tags = tags(ctx, c.api)
	md.Region, _ = get(ctx, c.api, pathPlacementRegion)
	md.VPC = vpc(ctx, c.api)

	// Extract the name from the map if it exists
	md.Name = md.Tags["Name"]

	return md, nil
}

func checkRoot(ctx context.Context, api ClientAPI) error {
	_, err := api.GetMetadata(ctx, &awsimds.GetMetadataInput{})
	return err
}

func get(ctx context.Context, api ClientAPI, path string) (string, error) {
	out, err := api.GetMetadata(ctx, &awsimds.GetMetadataInput{
		Path: path,
	})
	if err != nil {
		return "", err
	}
	defer out.Content.Close()

	data, _ := io.ReadAll(out.Content)
	return string(data), nil
}

func vpc(ctx context.Context, api ClientAPI) string {
	mac, _ := get(ctx, api, pathMacAddress)
	vpcID, _ := get(ctx, api, fmt.Sprintf("network/interfaces/macs/%s/vpc-id", mac))
	return vpcID
}

func tags(ctx context.Context, api ClientAPI) map[string]string {
	tagPaths, err := get(ctx, api, pathTagsInstance)
	if err != nil {
		return map[string]string{}
	}

	tags := map[string]string{}
	for _, tagName := range strings.Split(tagPaths, "\n") {
		tag, _ := get(ctx, api, pathTagsInstance+"/"+tagName)
		tags[tagName] = tag
	}

	return tags
}
