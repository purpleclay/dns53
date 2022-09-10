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

package r53_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsr53 "github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/purpleclay/dns53/internal/r53"
	"github.com/purpleclay/dns53/internal/r53/r53mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errAPI = errors.New("api error")

func TestByIDStripsPrefix(t *testing.T) {
	id := "Z0011223344HHGHGH"

	out := &awsr53.GetHostedZoneOutput{
		HostedZone: &types.HostedZone{
			Id:   aws.String("/hostedzone/" + id),
			Name: aws.String("testing"),
		},
	}

	m := r53mock.New(t)
	m.On("GetHostedZone", mock.Anything, mock.MatchedBy(func(req *awsr53.GetHostedZoneInput) bool {
		return *req.Id == id
	}), mock.Anything).Return(out, nil)

	c := r53.NewFromAPI(m)
	phz, err := c.ByID(context.Background(), id)

	require.NoError(t, err)
	assert.Equal(t, id, phz.ID)
}

func TestByIDError(t *testing.T) {
	m := r53mock.New(t)
	m.On("GetHostedZone", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.GetHostedZoneOutput{}, errAPI)

	c := r53.NewFromAPI(m)
	_, err := c.ByID(context.Background(), "")

	assert.Error(t, err)
}

func TestByVPCTrimsDotSuffix(t *testing.T) {
	m := r53mock.New(t)
	m.On("ListHostedZonesByVPC", mock.Anything, mock.MatchedBy(func(req *awsr53.ListHostedZonesByVPCInput) bool {
		return *req.VPCId == "vpc-12345" &&
			req.VPCRegion == types.VPCRegion("eu-west-2")
	}), mock.Anything).Return(&awsr53.ListHostedZonesByVPCOutput{
		HostedZoneSummaries: []types.HostedZoneSummary{
			{
				HostedZoneId: aws.String("Z0011223344HHGHGH"),
				Name:         aws.String("testing1."),
			},
			{
				HostedZoneId: aws.String("Z0099887766JKJJKJ"),
				Name:         aws.String("testing2."),
			},
		},
	}, nil)

	c := r53.NewFromAPI(m)
	phzs, err := c.ByVPC(context.Background(), "vpc-12345", "eu-west-2")

	require.NoError(t, err)

	expected := []r53.PrivateHostedZone{
		{
			ID:   "Z0011223344HHGHGH",
			Name: "testing1",
		},
		{
			ID:   "Z0099887766JKJJKJ",
			Name: "testing2",
		},
	}
	assert.ElementsMatch(t, expected, phzs)
}

func TestByVPCError(t *testing.T) {
	m := r53mock.New(t)
	m.On("ListHostedZonesByVPC", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.ListHostedZonesByVPCOutput{}, errAPI)

	c := r53.NewFromAPI(m)
	_, err := c.ByVPC(context.Background(), "vpc-12345", "eu-west-2")

	assert.Error(t, err)
}

func TestAssociateRecord(t *testing.T) {
	res := r53.ResourceRecord{
		PhzID:    "Z0011223344HHGHGH",
		Name:     "testing",
		Resource: "testing.zone",
	}

	m := r53mock.New(t)
	m.On("ChangeResourceRecordSets", mock.Anything, mock.MatchedBy(func(req *awsr53.ChangeResourceRecordSetsInput) bool {
		change := req.ChangeBatch.Changes[0]

		return *req.HostedZoneId == res.PhzID &&
			change.Action == types.ChangeActionCreate &&
			*change.ResourceRecordSet.Name == res.Name &&
			change.ResourceRecordSet.Type == types.RRTypeA &&
			*change.ResourceRecordSet.ResourceRecords[0].Value == res.Resource &&
			*change.ResourceRecordSet.TTL == int64(300)
	}), mock.Anything).Return(&awsr53.ChangeResourceRecordSetsOutput{}, nil)

	c := r53.NewFromAPI(m)
	err := c.AssociateRecord(context.Background(), res)

	assert.NoError(t, err)
}

func TestAssociateRecordError(t *testing.T) {
	m := r53mock.New(t)
	m.On("ChangeResourceRecordSets", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.ChangeResourceRecordSetsOutput{}, errAPI)

	c := r53.NewFromAPI(m)
	err := c.AssociateRecord(context.Background(), r53.ResourceRecord{})

	assert.Error(t, err)
}

func TestDisassociateRecord(t *testing.T) {
	res := r53.ResourceRecord{
		PhzID:    "Z0011223344HHGHGH",
		Name:     "testing",
		Resource: "testing.zone",
	}

	m := r53mock.New(t)
	m.On("ChangeResourceRecordSets", mock.Anything, mock.MatchedBy(func(req *awsr53.ChangeResourceRecordSetsInput) bool {
		change := req.ChangeBatch.Changes[0]

		return *req.HostedZoneId == res.PhzID &&
			change.Action == types.ChangeActionDelete &&
			*change.ResourceRecordSet.Name == res.Name &&
			change.ResourceRecordSet.Type == types.RRTypeA &&
			*change.ResourceRecordSet.ResourceRecords[0].Value == res.Resource &&
			*change.ResourceRecordSet.TTL == int64(300)
	}), mock.Anything).Return(&awsr53.ChangeResourceRecordSetsOutput{}, nil)

	c := r53.NewFromAPI(m)
	err := c.DisassociateRecord(context.Background(), res)

	assert.NoError(t, err)
}

func TestDisassociateRecordError(t *testing.T) {
	m := r53mock.New(t)
	m.On("ChangeResourceRecordSets", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.ChangeResourceRecordSetsOutput{}, errAPI)

	c := r53.NewFromAPI(m)
	err := c.DisassociateRecord(context.Background(), r53.ResourceRecord{})

	assert.Error(t, err)
}
