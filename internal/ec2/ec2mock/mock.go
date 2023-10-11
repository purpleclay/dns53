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

package ec2mock

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/stretchr/testify/mock"
)

type ClientAPI struct {
	mock.Mock
}

func (m *ClientAPI) ModifyInstanceMetadataOptions(ctx context.Context, params *ec2.ModifyInstanceMetadataOptionsInput, optFns ...func(*ec2.Options)) (*ec2.ModifyInstanceMetadataOptionsOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*ec2.ModifyInstanceMetadataOptionsOutput), args.Error(1)
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
