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

package ec2

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
)

const (
	// IPv4Path defines the path to the private IPv4 address of the EC2
	// instance within IMDS
	IPv4Path = "local-ipv4"

	// MacAddress defines the path to the private MAC address of the EC2
	// instance within IMDS
	MacAddress = "mac"
)

// Metadata contains metadata associated with an EC2 instance
type Metadata struct {
	IPv4   string
	Region string
	VPC    string
}

// InstanceMetadata attempts to retrieve useful metadata associated with
// the current EC2 instance by querying IMDS
func InstanceMetadata(cfg aws.Config) (Metadata, error) {
	c := imds.NewFromConfig(cfg)

	region, err := Region(c)
	if err != nil {
		return Metadata{}, err
	}

	ipv4, err := IPv4Address(c)
	if err != nil {
		return Metadata{}, err
	}

	vpc, err := VPC(c)
	if err != nil {
		return Metadata{}, err
	}

	return Metadata{
		Region: region,
		IPv4:   ipv4,
		VPC:    vpc,
	}, nil
}

// Region retrieves the region associated with the current EC2
func Region(c *imds.Client) (string, error) {
	out, err := c.GetRegion(context.TODO(), &imds.GetRegionInput{})
	if err != nil {
		return "", err
	}

	return out.Region, nil
}

// IPv4Address retrieves the IPv4 address associated with the current EC2
func IPv4Address(c *imds.Client) (string, error) {
	out, err := c.GetMetadata(context.TODO(), &imds.GetMetadataInput{
		Path: IPv4Path,
	})
	if err != nil {
		return "", err
	}
	defer out.Content.Close()

	data, _ := ioutil.ReadAll(out.Content)
	return string(data), nil
}

// VPC retrieves the VPC associated with the current EC2
func VPC(c *imds.Client) (string, error) {
	mac, err := c.GetMetadata(context.TODO(), &imds.GetMetadataInput{
		Path: MacAddress,
	})
	if err != nil {
		return "", err
	}
	defer mac.Content.Close()
	md, _ := ioutil.ReadAll(mac.Content)

	// Use the MAC address to retrieve the VPC associated with the EC2 instance
	out, err := c.GetMetadata(context.TODO(), &imds.GetMetadataInput{
		Path: fmt.Sprintf("network/interfaces/macs/%s/vpc-id", string(md)),
	})
	if err != nil {
		return "", err
	}
	defer mac.Content.Close()
	data, _ := ioutil.ReadAll(out.Content)

	return string(data), nil
}
