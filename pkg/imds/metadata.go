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

const (
	pathIPv4            = "local-ipv4"
	pathMacAddress      = "mac"
	pathPlacementRegion = "placement/region"
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
}

// NewFromAPI returns a new client from the provided IMDS API implementation
func NewFromAPI(api MetadataClientAPI) *Client {
	return &Client{api: api}
}

// InstanceMetadata attempts to retrieve useful metadata associated with
// the current EC2 instance by querying IMDS
func (c *Client) InstanceMetadata(ctx context.Context) (Metadata, error) {
	region, err := region(ctx, c.api)
	if err != nil {
		return Metadata{}, err
	}

	ipv4, err := ipv4(ctx, c.api)
	if err != nil {
		return Metadata{}, err
	}

	vpcID, err := vpc(ctx, c.api)
	if err != nil {
		return Metadata{}, err
	}

	return Metadata{
		Region: region,
		IPv4:   ipv4,
		VPC:    vpcID,
	}, nil
}

func region(ctx context.Context, api MetadataClientAPI) (string, error) {
	out, err := api.GetMetadata(ctx, &awsimds.GetMetadataInput{
		Path: pathPlacementRegion,
	})
	if err != nil {
		return "", err
	}
	defer out.Content.Close()

	data, _ := ioutil.ReadAll(out.Content)
	return string(data), nil
}

func ipv4(ctx context.Context, api MetadataClientAPI) (string, error) {
	out, err := api.GetMetadata(ctx, &awsimds.GetMetadataInput{
		Path: pathIPv4,
	})
	if err != nil {
		return "", err
	}
	defer out.Content.Close()

	data, _ := ioutil.ReadAll(out.Content)
	return string(data), nil
}

func vpc(ctx context.Context, api MetadataClientAPI) (string, error) {
	mac, err := api.GetMetadata(ctx, &awsimds.GetMetadataInput{
		Path: pathMacAddress,
	})
	if err != nil {
		return "", err
	}
	defer mac.Content.Close()
	md, _ := ioutil.ReadAll(mac.Content)

	// Use the MAC address to retrieve the VPC associated with the EC2 instance
	out, err := api.GetMetadata(ctx, &awsimds.GetMetadataInput{
		Path: fmt.Sprintf("network/interfaces/macs/%s/vpc-id", string(md)),
	})
	if err != nil {
		return "", err
	}
	defer mac.Content.Close()
	data, _ := ioutil.ReadAll(out.Content)

	return string(data), nil
}
