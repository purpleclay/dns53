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

package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// InstanceMetadataToggle allows the enabling and disabling of EC2 instance tags within IMDS
type InstanceMetadataToggle string

const (
	InstanceMetadataToggleEnabled  InstanceMetadataToggle = "enabled"
	InstanceMetadataToggleDisabled InstanceMetadataToggle = "disabled"
)

// ClientAPI defines the API for interacting with the Amazon EC2 service
type ClientAPI interface {
	// ModifyInstanceMetadataOptions modifies the parameters of a running EC2 instance,
	// by toggling the availability of EC2 instances tags within the Instance Metadata
	// Service (IMDS)
	ModifyInstanceMetadataOptions(ctx context.Context, params *awsec2.ModifyInstanceMetadataOptionsInput, optFns ...func(*awsec2.Options)) (*awsec2.ModifyInstanceMetadataOptionsOutput, error)
}

// Client defines the client for interacting with the Amazon EC2 service
type Client struct {
	api ClientAPI
}

// NewFromAPI returns a new client from the provided EC2 API implementation
func NewFromAPI(api ClientAPI) *Client {
	return &Client{api: api}
}

// ToggleInstanceMetadataTags will modify the parameters of a running EC2 instance,
// by toggling the availability of EC2 instance tags within the Instance Metadata
// Service.
//
// The equivalent operation can be achieved through the CLI using:
//
//	aws ec2 modify-instance-metadata-options --instance-id <INSTANCE_ID> --instance-metadata-tags enabled
func (c *Client) ToggleInstanceMetadataTags(ctx context.Context, id string, toggle InstanceMetadataToggle) error {
	_, err := c.api.ModifyInstanceMetadataOptions(ctx, &awsec2.ModifyInstanceMetadataOptionsInput{
		InstanceId:           aws.String(id),
		InstanceMetadataTags: types.InstanceMetadataTagsState(toggle),
	})

	return err
}
