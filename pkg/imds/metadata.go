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
	// IPv4Path defines the path to the private IPv4 address of the EC2
	// instance within IMDS
	IPv4Path = "local-ipv4"

	// MacAddress defines the path to the private MAC address of the EC2
	// instance within IMDS
	MacAddress = "mac"
)

// MetadataClientAPI ...
type MetadataClientAPI interface {
	// GetRegion ...
	GetRegion(ctx context.Context, params *awsimds.GetRegionInput, optFns ...func(*awsimds.Options)) (*awsimds.GetRegionOutput, error)

	// GetMetadata ...
	GetMetadata(ctx context.Context, params *awsimds.GetMetadataInput, optFns ...func(*awsimds.Options)) (*awsimds.GetMetadataOutput, error)
}

// Client ...
type Client struct {
	api MetadataClientAPI
}

// Metadata contains metadata associated with an EC2 instance
type Metadata struct {
	IPv4   string
	Region string
	VPC    string
}

// NewFromAPI ...
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
	out, err := api.GetRegion(context.TODO(), &awsimds.GetRegionInput{})
	if err != nil {
		return "", err
	}

	return out.Region, nil
}

func ipv4(ctx context.Context, api MetadataClientAPI) (string, error) {
	out, err := api.GetMetadata(context.TODO(), &awsimds.GetMetadataInput{
		Path: IPv4Path,
	})
	if err != nil {
		return "", err
	}
	defer out.Content.Close()

	data, _ := ioutil.ReadAll(out.Content)
	return string(data), nil
}

func vpc(ctx context.Context, api MetadataClientAPI) (string, error) {
	mac, err := api.GetMetadata(context.TODO(), &awsimds.GetMetadataInput{
		Path: MacAddress,
	})
	if err != nil {
		return "", err
	}
	defer mac.Content.Close()
	md, _ := ioutil.ReadAll(mac.Content)

	// Use the MAC address to retrieve the VPC associated with the EC2 instance
	out, err := api.GetMetadata(context.TODO(), &awsimds.GetMetadataInput{
		Path: fmt.Sprintf("network/interfaces/macs/%s/vpc-id", string(md)),
	})
	if err != nil {
		return "", err
	}
	defer mac.Content.Close()
	data, _ := ioutil.ReadAll(out.Content)

	return string(data), nil
}
