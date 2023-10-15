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

package r53mock

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/stretchr/testify/mock"
)

type ClientAPI struct {
	mock.Mock
}

func (m *ClientAPI) CreateHostedZone(ctx context.Context, params *route53.CreateHostedZoneInput, optFns ...func(*route53.Options)) (*route53.CreateHostedZoneOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*route53.CreateHostedZoneOutput), args.Error(1)
}

func (m *ClientAPI) DeleteHostedZone(ctx context.Context, params *route53.DeleteHostedZoneInput, optFns ...func(*route53.Options)) (*route53.DeleteHostedZoneOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*route53.DeleteHostedZoneOutput), args.Error(1)
}

func (m *ClientAPI) GetHostedZone(ctx context.Context, params *route53.GetHostedZoneInput, optFns ...func(*route53.Options)) (*route53.GetHostedZoneOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*route53.GetHostedZoneOutput), args.Error(1)
}

func (m *ClientAPI) ListHostedZonesByVPC(ctx context.Context, params *route53.ListHostedZonesByVPCInput, optFns ...func(*route53.Options)) (*route53.ListHostedZonesByVPCOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*route53.ListHostedZonesByVPCOutput), args.Error(1)
}

func (m *ClientAPI) ListHostedZonesByName(ctx context.Context, params *route53.ListHostedZonesByNameInput, optFns ...func(*route53.Options)) (*route53.ListHostedZonesByNameOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*route53.ListHostedZonesByNameOutput), args.Error(1)
}

func (m *ClientAPI) ChangeResourceRecordSets(ctx context.Context, params *route53.ChangeResourceRecordSetsInput, optFns ...func(*route53.Options)) (*route53.ChangeResourceRecordSetsOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*route53.ChangeResourceRecordSetsOutput), args.Error(1)
}

func (m *ClientAPI) AssociateVPCWithHostedZone(ctx context.Context, params *route53.AssociateVPCWithHostedZoneInput, optFns ...func(*route53.Options)) (*route53.AssociateVPCWithHostedZoneOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*route53.AssociateVPCWithHostedZoneOutput), args.Error(1)
}

func (m *ClientAPI) DisassociateVPCFromHostedZone(ctx context.Context, params *route53.DisassociateVPCFromHostedZoneInput, optFns ...func(*route53.Options)) (*route53.DisassociateVPCFromHostedZoneOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*route53.DisassociateVPCFromHostedZoneOutput), args.Error(1)
}

func New(tb testing.TB) *ClientAPI {
	tb.Helper()

	mock := &ClientAPI{}
	mock.Mock.Test(tb)

	tb.Cleanup(func() {
		mock.AssertExpectations(tb)
	})

	return mock
}
