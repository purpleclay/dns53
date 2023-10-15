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

package r53

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsr53 "github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

const (
	dotSuffix        = "."
	hostedZonePrefix = "/hostedzone/"
)

// DNSClientAPI defines the API for interacting with Amazon Route 53
type DNSClientAPI interface {
	// CreateHostedZone creates a new private hosted zone
	CreateHostedZone(ctx context.Context, params *awsr53.CreateHostedZoneInput, optFns ...func(*awsr53.Options)) (*awsr53.CreateHostedZoneOutput, error)

	// DeleteHostedZone deletes an existing private hosted zone
	DeleteHostedZone(ctx context.Context, params *awsr53.DeleteHostedZoneInput, optFns ...func(*awsr53.Options)) (*awsr53.DeleteHostedZoneOutput, error)

	// GetHostedZone retrieves information about a specified hosted zone including
	// the four name servers assigned to the hosted zone
	GetHostedZone(ctx context.Context, params *awsr53.GetHostedZoneInput, optFns ...func(*awsr53.Options)) (*awsr53.GetHostedZoneOutput, error)

	// ListHostedZonesByVPC lists all of the private hosted zones that a specified VPC
	// is associated with
	ListHostedZonesByVPC(ctx context.Context, params *awsr53.ListHostedZonesByVPCInput, optFns ...func(*awsr53.Options)) (*awsr53.ListHostedZonesByVPCOutput, error)

	// ListHostedZonesByName lists all of the private hosted zones that contain a specific
	// name in lexicographic order
	ListHostedZonesByName(ctx context.Context, params *awsr53.ListHostedZonesByNameInput, optFns ...func(*awsr53.Options)) (*awsr53.ListHostedZonesByNameOutput, error)

	// ChangeResourceRecordSets creates, changes, or deletes a resource record set,
	// which contains authoritative DNS information for a specified domain name or subdomain name
	ChangeResourceRecordSets(ctx context.Context, params *awsr53.ChangeResourceRecordSetsInput, optFns ...func(*awsr53.Options)) (*awsr53.ChangeResourceRecordSetsOutput, error)

	// AssociateVPCWithHostedZone will associate a VPC with a Route53 Private Hosted Zone
	AssociateVPCWithHostedZone(ctx context.Context, params *awsr53.AssociateVPCWithHostedZoneInput, optFns ...func(*awsr53.Options)) (*awsr53.AssociateVPCWithHostedZoneOutput, error)

	// DisassociateVPCFromHostedZone will disassociate a VPC with a Route53 Private Hosted Zone
	DisassociateVPCFromHostedZone(ctx context.Context, params *awsr53.DisassociateVPCFromHostedZoneInput, optFns ...func(*awsr53.Options)) (*awsr53.DisassociateVPCFromHostedZoneOutput, error)
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

// CreatePrivateHostedZone will attempt to create a Route53 Private Hosted Zone
// with the given domain name and associate it with the required VPC
//
// The equivalent operation can be achieved through the CLI using:
//
//	aws route53 create-hosted-zone --name <DOMAIN_NAME> --vpc VPCId=<VPC_ID>,VPCRegion=<VPC_REGION> \
//	 --hosted-zone-config PrivateZone=true --caller-reference <REFERENCE>
func (r *Client) CreatePrivateHostedZone(ctx context.Context, name, vpc, region string) (PrivateHostedZone, error) {
	resp, err := r.api.CreateHostedZone(ctx, &awsr53.CreateHostedZoneInput{
		Name: aws.String(name),
		VPC: &types.VPC{
			VPCId:     aws.String(vpc),
			VPCRegion: types.VPCRegion(region),
		},
		HostedZoneConfig: &types.HostedZoneConfig{
			PrivateZone: true,
		},
		CallerReference: aws.String(time.Now().Format(time.RFC3339)),
	})
	if err != nil {
		return PrivateHostedZone{}, err
	}

	return PrivateHostedZone{
		ID:   strings.TrimPrefix(*resp.HostedZone.Id, hostedZonePrefix),
		Name: strings.TrimSuffix(*resp.HostedZone.Name, dotSuffix),
	}, nil
}

// DeletePrivateHostedZone will attempt to delete an existing Route53 Private Hosted Zone
// by its ID. If the hosted zone contains any record sets, the deletion will fail
//
// The equivalent operation can be achieved through the CLI using:
//
//	aws route53 delete-hosted-zone --id <HOSTED_ZONE_ID>
func (r *Client) DeletePrivateHostedZone(ctx context.Context, id string) error {
	_, err := r.api.DeleteHostedZone(ctx, &awsr53.DeleteHostedZoneInput{
		Id: aws.String(id),
	})
	if err != nil {
		// The hosted zone may be owned by another process and contain record sets.
		// It is not deemed a failure when deletion fails in this scenario
		var errNotEmpty *types.HostedZoneNotEmpty
		if errors.As(err, &errNotEmpty) {
			return nil
		}
	}

	return err
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
		ID:   strings.TrimPrefix(*resp.HostedZone.Id, hostedZonePrefix),
		Name: strings.TrimSuffix(*resp.HostedZone.Name, dotSuffix),
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
			Name: strings.TrimSuffix(*hzs.Name, dotSuffix),
		})
	}

	return phzs, nil
}

// ByName will attempt to find a Route53 Private Hosted Zone that exactly matches
// the given domain name
//
// The equivalent operation can be achieved through thr CLI using:
//
//	aws route53 list-hosted-zones-by-name --dns-name <DOMAIN_NAME>
func (r *Client) ByName(ctx context.Context, name string) (*PrivateHostedZone, error) {
	resp, err := r.api.ListHostedZonesByName(ctx, &awsr53.ListHostedZonesByNameInput{
		DNSName: aws.String(name),
	})
	if err != nil {
		return nil, err
	}

	for _, hzs := range resp.HostedZones {
		// Stop on the first matching private hosted zone
		if hzs.Config.PrivateZone {
			hzName := strings.TrimSuffix(*hzs.Name, dotSuffix)

			if hzName == name {
				return &PrivateHostedZone{
					ID:   strings.TrimPrefix(*hzs.Id, hostedZonePrefix),
					Name: strings.TrimSuffix(*hzs.Name, dotSuffix),
				}, nil
			}
		}
	}

	return nil, nil
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

// AssociateVPCWithZone attempts to associate a given VPC with a Route53 Private
// Hosted Zone. This is required to support any future DNS resolution queries.
// If an association already exists between the VPC and PHZ, a ConflictingDomainExists
// error is thrown by the AWS SDK and handled accordingly
//
// The equivalent operation can be achieved through the CLI using:
//
//	aws route53 associate-vpc-with-hosted-zone --hosted-zone-id <HOSTED_ZONE_ID> \
//	 --vpc VPCId=<VPC_ID>,VPCRegion=<VPC_REGION>
func (r *Client) AssociateVPCWithZone(ctx context.Context, id, vpc, region string) error {
	_, err := r.api.AssociateVPCWithHostedZone(ctx, &awsr53.AssociateVPCWithHostedZoneInput{
		HostedZoneId: aws.String(id),
		VPC: &types.VPC{
			VPCId:     aws.String(vpc),
			VPCRegion: types.VPCRegion(region),
		},
	})
	if err != nil {
		// If an association already exists between the VPC and PHZ, swallow the error
		var errAssocExists *types.ConflictingDomainExists
		if errors.As(err, &errAssocExists) {
			return nil
		}
	}

	return err
}

// DisassociateVPCWithZone attempts to disassociate a given VPC with a Route53
// Private Hosted Zone. If no association exists between the VPC and PHZ, a
// VPCAssociationNotFound error is thrown by the AWS SDK and handled accordingly
//
// The equivalent operation can be achieved through the CLI using:
//
//	aws route53 disassociate-vpc-with-hosted-zone --hosted-zone-id <HOSTED_ZONE_ID> \
//	 --vpc VPCId=<VPC_ID>,VPCRegion=<VPC_REGION>
func (r *Client) DisassociateVPCWithZone(ctx context.Context, id, vpc, region string) error {
	_, err := r.api.DisassociateVPCFromHostedZone(ctx, &awsr53.DisassociateVPCFromHostedZoneInput{
		HostedZoneId: aws.String(id),
		VPC: &types.VPC{
			VPCId:     aws.String(vpc),
			VPCRegion: types.VPCRegion(region),
		},
	})
	if err != nil {
		// If no association exists, swallow the error
		var errNoAssocExists *types.VPCAssociationNotFound
		if errors.As(err, &errNoAssocExists) {
			return nil
		}
	}

	return err
}
