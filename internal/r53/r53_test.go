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

const (
	vpcID  = "vpc-12345"
	region = "eu-west-2"
)

var errAPI = errors.New("api error")

func TestCreatePrivateHostedZone(t *testing.T) {
	domain := "dns53"

	m := r53mock.New(t)
	m.On("CreateHostedZone", mock.Anything, mock.MatchedBy(func(req *awsr53.CreateHostedZoneInput) bool {
		return true
	}), mock.Anything).Return(&awsr53.CreateHostedZoneOutput{
		HostedZone: &types.HostedZone{
			Id:   aws.String("/hostedzone/Z00000000000001"),
			Name: aws.String(domain + "."),
		},
	}, nil)

	c := r53.NewFromAPI(m)
	hz, err := c.CreatePrivateHostedZone(context.Background(), "dns53", vpcID, region)

	require.NoError(t, err)
	assert.Equal(t, "Z00000000000001", hz.ID)
	assert.Equal(t, domain, hz.Name)
}

func TestCreatePrivateHostedZoneError(t *testing.T) {
	m := r53mock.New(t)
	m.On("CreateHostedZone", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.CreateHostedZoneOutput{}, errAPI)

	c := r53.NewFromAPI(m)
	_, err := c.CreatePrivateHostedZone(context.Background(), "Z00000000000002", vpcID, region)

	assert.Error(t, err)
}

func TestDeletePrivateHostedZone(t *testing.T) {
	m := r53mock.New(t)
	m.On("DeleteHostedZone", mock.Anything, mock.MatchedBy(func(req *awsr53.DeleteHostedZoneInput) bool {
		return *req.Id == "Z00000000000003"
	}), mock.Anything).Return(&awsr53.DeleteHostedZoneOutput{}, nil)

	c := r53.NewFromAPI(m)
	err := c.DeletePrivateHostedZone(context.Background(), "Z00000000000003")

	assert.NoError(t, err)
}

func TestDeletePrivateHostedZoneError(t *testing.T) {
	m := r53mock.New(t)
	m.On("DeleteHostedZone", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.DeleteHostedZoneOutput{}, errAPI)

	c := r53.NewFromAPI(m)
	err := c.DeletePrivateHostedZone(context.Background(), "Z00000000000004")

	assert.Error(t, err)
}

func TestDeletePrivateHostedZoneNotEmpty(t *testing.T) {
	m := r53mock.New(t)
	m.On("DeleteHostedZone", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.DeleteHostedZoneOutput{}, &types.HostedZoneNotEmpty{})

	c := r53.NewFromAPI(m)
	err := c.DeletePrivateHostedZone(context.Background(), "Z00000000000005")

	assert.NoError(t, err)
}

func TestByIDStripsPrefix(t *testing.T) {
	id := "Z00000000000006"

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
		return *req.VPCId == vpcID &&
			req.VPCRegion == types.VPCRegion(region)
	}), mock.Anything).Return(&awsr53.ListHostedZonesByVPCOutput{
		HostedZoneSummaries: []types.HostedZoneSummary{
			{
				HostedZoneId: aws.String("Z00000000000007"),
				Name:         aws.String("testing1."),
			},
			{
				HostedZoneId: aws.String("Z00000000000008"),
				Name:         aws.String("testing2."),
			},
		},
	}, nil)

	c := r53.NewFromAPI(m)
	phzs, err := c.ByVPC(context.Background(), vpcID, region)

	require.NoError(t, err)

	expected := []r53.PrivateHostedZone{
		{
			ID:   "Z00000000000007",
			Name: "testing1",
		},
		{
			ID:   "Z00000000000008",
			Name: "testing2",
		},
	}
	assert.ElementsMatch(t, expected, phzs)
}

func TestByVPCError(t *testing.T) {
	m := r53mock.New(t)
	m.On("ListHostedZonesByVPC", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.ListHostedZonesByVPCOutput{}, errAPI)

	c := r53.NewFromAPI(m)
	_, err := c.ByVPC(context.Background(), vpcID, region)

	assert.Error(t, err)
}

func TestByNameExcludesPublicZones(t *testing.T) {
	domain := "dns53"

	m := r53mock.New(t)
	m.On("ListHostedZonesByName", mock.Anything, mock.MatchedBy(func(req *awsr53.ListHostedZonesByNameInput) bool {
		return *req.DNSName == domain
	}), mock.Anything).Return(&awsr53.ListHostedZonesByNameOutput{
		HostedZones: []types.HostedZone{
			{
				Id:   aws.String("/hostedzone/Z00000000000009"),
				Name: aws.String("dns53."),
				Config: &types.HostedZoneConfig{
					PrivateZone: false,
				},
			},
			{
				Id:   aws.String("/hostedzone/Z00000000000010"),
				Name: aws.String("dns53."),
				Config: &types.HostedZoneConfig{
					PrivateZone: true,
				},
			},
		},
	}, nil)

	c := r53.NewFromAPI(m)
	hz, err := c.ByName(context.Background(), domain)

	require.NoError(t, err)
	require.NotNil(t, hz)
	assert.Equal(t, "Z00000000000010", hz.ID)
	assert.Equal(t, domain, hz.Name)
}

func TestByNameExactMatchOnly(t *testing.T) {
	domain := "dns53"

	m := r53mock.New(t)
	m.On("ListHostedZonesByName", mock.Anything, mock.MatchedBy(func(req *awsr53.ListHostedZonesByNameInput) bool {
		return *req.DNSName == domain
	}), mock.Anything).Return(&awsr53.ListHostedZonesByNameOutput{
		HostedZones: []types.HostedZone{
			{
				Id:   aws.String("/hostedzone/Z00000000000011"),
				Name: aws.String("dns53zone."),
				Config: &types.HostedZoneConfig{
					PrivateZone: true,
				},
			},
			{
				Id:   aws.String("/hostedzone/Z00000000000012"),
				Name: aws.String("dns53."),
				Config: &types.HostedZoneConfig{
					PrivateZone: true,
				},
			},
		},
	}, nil)

	c := r53.NewFromAPI(m)
	hz, err := c.ByName(context.Background(), domain)

	require.NoError(t, err)
	require.NotNil(t, hz)
	assert.Equal(t, "Z00000000000012", hz.ID)
	assert.Equal(t, domain, hz.Name)
}

func TestByNameNoMatch(t *testing.T) {
	m := r53mock.New(t)
	m.On("ListHostedZonesByName", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.ListHostedZonesByNameOutput{}, nil)

	c := r53.NewFromAPI(m)
	hz, err := c.ByName(context.Background(), "Z00000000000013")

	require.NoError(t, err)
	assert.Nil(t, hz)
}

func TestByNameError(t *testing.T) {
	m := r53mock.New(t)
	m.On("ListHostedZonesByName", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.ListHostedZonesByNameOutput{}, errAPI)

	c := r53.NewFromAPI(m)
	hz, err := c.ByName(context.Background(), "error")

	require.Error(t, err)
	assert.Nil(t, hz)
}

func TestAssociateRecord(t *testing.T) {
	res := r53.ResourceRecord{
		PhzID:    "Z00000000000014",
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
		PhzID:    "Z00000000000015",
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

func TestAssociateVPCWithZone(t *testing.T) {
	m := r53mock.New(t)
	m.On("AssociateVPCWithHostedZone", mock.Anything, mock.MatchedBy(func(req *awsr53.AssociateVPCWithHostedZoneInput) bool {
		return *req.HostedZoneId == "Z00000000000016" &&
			*req.VPC.VPCId == vpcID &&
			req.VPC.VPCRegion == types.VPCRegion(region)
	}), mock.Anything).Return(&awsr53.AssociateVPCWithHostedZoneOutput{}, nil)

	c := r53.NewFromAPI(m)
	err := c.AssociateVPCWithZone(context.Background(), "Z00000000000016", vpcID, region)

	require.NoError(t, err)
}

func TestAssociateVPCWithZoneError(t *testing.T) {
	m := r53mock.New(t)
	m.On("AssociateVPCWithHostedZone", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.AssociateVPCWithHostedZoneOutput{}, errAPI)

	c := r53.NewFromAPI(m)
	err := c.AssociateVPCWithZone(context.Background(), "Z00000000000017", vpcID, region)

	assert.Error(t, err)
}

func TestAssociateVPCWithZoneDuplicateAssociationIgnored(t *testing.T) {
	m := r53mock.New(t)
	m.On("AssociateVPCWithHostedZone", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.AssociateVPCWithHostedZoneOutput{}, &types.ConflictingDomainExists{})

	c := r53.NewFromAPI(m)
	err := c.AssociateVPCWithZone(context.Background(), "Z00000000000018", vpcID, region)

	assert.NoError(t, err)
}
