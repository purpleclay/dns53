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

// TODO: turn this into an interface that can be used to mock tests

// PrivateHostedZone identifies an AWS Route53 Private Hosted Zone (PHZ)
type PrivateHostedZone struct {
	ID   string
	Name string
}

// ByID attempts to retrieve a Route53 Private Hosted Zone by its given ID
func ByID(cfg aws.Config, id string) (PrivateHostedZone, error) {
	c := awsr53.NewFromConfig(cfg)

	resp, err := c.GetHostedZone(context.TODO(), &awsr53.GetHostedZoneInput{
		Id: aws.String(id),
	})
	if err != nil {
		return PrivateHostedZone{}, err
	}

	return PrivateHostedZone{
		ID:   strings.TrimPrefix(*resp.HostedZone.Id, "/hostedzone/"),
		Name: *resp.HostedZone.Name,
	}, nil
}

// ByVPC finds all Route53 Private Hosted Zones associated with a given VPC ID
func ByVPC(cfg aws.Config, vpc, region string) ([]PrivateHostedZone, error) {
	c := awsr53.NewFromConfig(cfg)

	resp, err := c.ListHostedZonesByVPC(context.TODO(), &awsr53.ListHostedZonesByVPCInput{
		VPCId:     aws.String(vpc),
		VPCRegion: types.VPCRegion(region),
	})
	if err != nil {
		return []PrivateHostedZone{}, err
	}

	phz := make([]PrivateHostedZone, 0, len(resp.HostedZoneSummaries))
	for _, hzs := range resp.HostedZoneSummaries {
		phz = append(phz, PrivateHostedZone{ID: *hzs.HostedZoneId, Name: *hzs.Name})
	}

	return phz, nil
}

// AssociateRecord creates a new A-Record entry within a given Route53 Private Hosted Zone
// for the specified Record Name and target EC2 IPv4 address
func AssociateRecord(cfg aws.Config, phzID, name, ip string) error {
	c := awsr53.NewFromConfig(cfg)

	_, err := c.ChangeResourceRecordSets(context.TODO(), &awsr53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(phzID),
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeActionCreate,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String(name),
						Type: types.RRTypeA,
						ResourceRecords: []types.ResourceRecord{
							{
								Value: aws.String(ip),
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
func DisassociateRecord(cfg aws.Config, phzID, name, ip string) error {
	c := awsr53.NewFromConfig(cfg)

	_, err := c.ChangeResourceRecordSets(context.TODO(), &awsr53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(phzID),
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeActionDelete,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String(name),
						Type: types.RRTypeA,
						ResourceRecords: []types.ResourceRecord{
							{
								Value: aws.String(ip),
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
