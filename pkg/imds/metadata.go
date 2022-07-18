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
	"io/ioutil"

	awsimds "github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
)

// TODO: comment paths
const (
	PathIPv4            = "local-ipv4"
	PathMacAddress      = "mac"
	PathPlacementRegion = "placement/region"
	PathPlacementAZ     = "placement/availability-zone"
	PathInstanceID      = "instance-id"
)

// MetadataClientAPI defines the API for interacting with the Amazon
// EC2 Instance Metadata Service (IMDS)
type MetadataClientAPI interface {
	// GetMetadata uses the path provided to request information from the Amazon
	// EC2 Instance Metadata Service
	GetMetadata(ctx context.Context, params *awsimds.GetMetadataInput, optFns ...func(*awsimds.Options)) (*awsimds.GetMetadataOutput, error)
}

// Client defines the client for interacting with the Amazon EC2 Instance
// Metadata Service (IMDS)
type Client struct {
	api MetadataClientAPI
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
}

// NewFromAPI returns a new client from the provided IMDS API implementation
func NewFromAPI(api MetadataClientAPI) *Client {
	return &Client{api: api}
}

// InstanceMetadata attempts to retrieve useful metadata associated with
// the current EC2 instance by querying IMDS
func (c *Client) InstanceMetadata(ctx context.Context) (Metadata, error) {
	md := Metadata{}

	var err error
	if md.Region, err = get(ctx, c.api, PathPlacementRegion); err != nil {
		return Metadata{}, err
	}

	if md.IPv4, err = get(ctx, c.api, PathIPv4); err != nil {
		return Metadata{}, err
	}

	if md.VPC, err = vpc(ctx, c.api); err != nil {
		return Metadata{}, err
	}

	if md.AZ, err = get(ctx, c.api, PathPlacementAZ); err != nil {
		return Metadata{}, err
	}

	if md.InstanceID, err = get(ctx, c.api, PathInstanceID); err != nil {
		return Metadata{}, err
	}

	return md, nil
}

func get(ctx context.Context, api MetadataClientAPI, path string) (string, error) {
	out, err := api.GetMetadata(ctx, &awsimds.GetMetadataInput{
		Path: path,
	})
	if err != nil {
		return "", err
	}
	defer out.Content.Close()

	data, _ := ioutil.ReadAll(out.Content)
	return string(data), nil
}

func vpc(ctx context.Context, api MetadataClientAPI) (string, error) {
	mac, err := get(ctx, api, PathMacAddress)
	if err != nil {
		return "", err
	}

	vpc, err := get(ctx, api, fmt.Sprintf("network/interfaces/macs/%s/vpc-id", mac))
	if err != nil {
		return "", err
	}

	return vpc, nil
}
