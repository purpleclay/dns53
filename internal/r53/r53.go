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

package r53

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsr53 "github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

// DNSClientAPI defines the API for interacting with Amazon Route 53
type DNSClientAPI interface {
	// GetHostedZone retrieves information about a specified hosted zone including
	// the four name servers assigned to the hosted zone
	GetHostedZone(ctx context.Context, params *awsr53.GetHostedZoneInput, optFns ...func(*awsr53.Options)) (*awsr53.GetHostedZoneOutput, error)

	// ListHostedZonesByVPC lists all of the private hosted zones that a specified VPC
	// is associated with
	ListHostedZonesByVPC(ctx context.Context, params *awsr53.ListHostedZonesByVPCInput, optFns ...func(*awsr53.Options)) (*awsr53.ListHostedZonesByVPCOutput, error)

	// ChangeResourceRecordSets creates, changes, or deletes a resource record set,
	// which contains authoritative DNS information for a specified domain name or subdomain name
	ChangeResourceRecordSets(ctx context.Context, params *awsr53.ChangeResourceRecordSetsInput, optFns ...func(*awsr53.Options)) (*awsr53.ChangeResourceRecordSetsOutput, error)
}

// Client defines the client for interacting with Amazon Route 53
type Client struct {
	api DNSClientAPI
}

// PrivateHostedZone identifies an AWS Route53 Private Hosted Zone (PHZ)
type PrivateHostedZone struct {
	// ID of the AWS Route53 Private Hosted Zone (PHZ)
	ID string

	// Name of the AWS Route53 Hosted Zone (PHZ). This will be the CNAME
	// of the parent domain
	Name string
}

// ResourceRecord represents a DNS record type that is supported by an
// AWS Route53 Private Hosted Zone (PHZ)
type ResourceRecord struct {
	// PhzID of the AWS Route53 Private Hosted Zone (PHZ)
	PhzID string

	// Name of the resource record that will be either be created, updated
	// or deleted within the AWS Route53 Private Hosted Zone (PHZ)
	Name string

	// Resource contains the value associated with the resource record
	Resource string
}

// NewFromAPI returns a new client from the provided DNS API implementation
func NewFromAPI(api DNSClientAPI) *Client {
	return &Client{api: api}
}

// ByID attempts to retrieve a Route53 Private Hosted Zone by its given ID
//
// The equivalent operation can be achieved through the CLI using:
//
//	aws route53 get-hosted-zone --id <HOSTED_ZONE_ID>
func (r *Client) ByID(ctx context.Context, id string) (PrivateHostedZone, error) {
	resp, err := r.api.GetHostedZone(ctx, &awsr53.GetHostedZoneInput{
		Id: aws.String(id),
	})
	if err != nil {
		return PrivateHostedZone{}, err
	}

	// Trim off the static prefix from the Hosted Zone ID
	return PrivateHostedZone{
		ID:   strings.TrimPrefix(*resp.HostedZone.Id, "/hostedzone/"),
		Name: strings.TrimSuffix(*resp.HostedZone.Name, "."),
	}, nil
}

// ByVPC finds all Route53 Private Hosted Zones associated with a given VPC ID
//
// The equivalent operation can be achieved through the CLI using:
//
//	aws route53 list-hosted-zones-by-vpc --vpc-id <VPC_ID> --vpc-region <REGION>
func (r *Client) ByVPC(ctx context.Context, vpc, region string) ([]PrivateHostedZone, error) {
	resp, err := r.api.ListHostedZonesByVPC(ctx, &awsr53.ListHostedZonesByVPCInput{
		VPCId:     aws.String(vpc),
		VPCRegion: types.VPCRegion(region),
	})
	if err != nil {
		return []PrivateHostedZone{}, err
	}

	phzs := make([]PrivateHostedZone, 0, len(resp.HostedZoneSummaries))
	for _, hzs := range resp.HostedZoneSummaries {
		phzs = append(phzs, PrivateHostedZone{
			ID:   *hzs.HostedZoneId,
			Name: strings.TrimSuffix(*hzs.Name, "."),
		})
	}

	return phzs, nil
}

// AssociateRecord creates a new A-Record entry within a given Route53 Private Hosted Zone
// for the specified Record Name and target EC2 IPv4 address
func (r *Client) AssociateRecord(ctx context.Context, res ResourceRecord) error {
	_, err := r.api.ChangeResourceRecordSets(ctx, &awsr53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(res.PhzID),
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeActionCreate,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String(res.Name),
						Type: types.RRTypeA,
						ResourceRecords: []types.ResourceRecord{
							{
								Value: aws.String(res.Resource),
							},
						},
						TTL: aws.Int64(300),
					},
				},
			},
		},
	})

	return err
}

// DisassociateRecord attempts to delete an existing A-Record entry within a given Route53
// Private Hosted Zone, based on the specified Record Name and target EC2 IPv4 address
func (r *Client) DisassociateRecord(ctx context.Context, res ResourceRecord) error {
	_, err := r.api.ChangeResourceRecordSets(ctx, &awsr53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(res.PhzID),
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeActionDelete,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String(res.Name),
						Type: types.RRTypeA,
						ResourceRecords: []types.ResourceRecord{
							{
								Value: aws.String(res.Resource),
							},
						},
						TTL: aws.Int64(300),
					},
				},
			},
		},
	})

	return err
}
