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
	"net"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsr53 "github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

// TODO: define an interface that can be used during testing

// PrivateHostedZone ...
type PrivateHostedZone struct {
	ID   string
	Name string
}

// ByVPC ...
func ByVPC(cfg aws.Config) ([]PrivateHostedZone, error) {
	c := awsr53.NewFromConfig(cfg)
	resp, err := c.ListHostedZonesByVPC(context.TODO(), &awsr53.ListHostedZonesByVPCInput{})
	if err != nil {
		return []PrivateHostedZone{}, err
	}

	phz := make([]PrivateHostedZone, 0, len(resp.HostedZoneSummaries))
	for _, hzs := range resp.HostedZoneSummaries {
		phz = append(phz, PrivateHostedZone{ID: *hzs.HostedZoneId, Name: *hzs.Name})
	}

	return phz, nil
}

// AssociateRecord ...
func AssociateRecord(cfg aws.Config, phzID, ip string) error {
	c := awsr53.NewFromConfig(cfg)

	// Select the correct resource record set type based on the IP address
	var rrtype types.RRType = types.RRTypeA
	pip := net.ParseIP(ip)
	if pip.To4() == nil {
		rrtype = types.RRTypeAaaa
	}

	_, err := c.ChangeResourceRecordSets(context.TODO(), &awsr53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(phzID),
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeActionCreate,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String(""),
						Type: rrtype,
						ResourceRecords: []types.ResourceRecord{
							{
								Value: aws.String(ip),
							},
						},
					},
				},
			},
		},
	})

	return err
}
