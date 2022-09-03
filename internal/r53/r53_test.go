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
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsr53 "github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/purpleclay/dns53/internal/r53"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock implementation of the AWS R53 SDK
type mockAPI struct {
	mock.Mock
}

func (m *mockAPI) GetHostedZone(ctx context.Context, params *awsr53.GetHostedZoneInput, optFns ...func(*awsr53.Options)) (*awsr53.GetHostedZoneOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*awsr53.GetHostedZoneOutput), args.Error(1)
}

func (m *mockAPI) ListHostedZonesByVPC(ctx context.Context, params *awsr53.ListHostedZonesByVPCInput, optFns ...func(*awsr53.Options)) (*awsr53.ListHostedZonesByVPCOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*awsr53.ListHostedZonesByVPCOutput), args.Error(1)
}

func (m *mockAPI) ChangeResourceRecordSets(ctx context.Context, params *awsr53.ChangeResourceRecordSetsInput, optFns ...func(*awsr53.Options)) (*awsr53.ChangeResourceRecordSetsOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*awsr53.ChangeResourceRecordSetsOutput), args.Error(1)
}

func TestByID_StripsHostedZonePrefix(t *testing.T) {
	id := "Z0011223344HHGHGH"

	out := &awsr53.GetHostedZoneOutput{
		HostedZone: &types.HostedZone{
			Id:   aws.String("/hostedzone/" + id),
			Name: aws.String("testing"),
		},
	}

	m := &mockAPI{}
	m.On("GetHostedZone", mock.Anything, mock.Anything, mock.Anything).Return(out, nil)

	c := r53.NewFromAPI(m)
	phz, err := c.ByID(context.TODO(), id)

	require.NoError(t, err)
	assert.Equal(t, id, phz.ID)
}
