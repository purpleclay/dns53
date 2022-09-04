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

package ec2_test

import (
	"context"
	"errors"
	"testing"

	awsec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/purpleclay/dns53/internal/ec2"
	"github.com/purpleclay/dns53/internal/ec2/ec2mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestToggleInstanceMetadataTags(t *testing.T) {
	tests := []struct {
		name       string
		instanceID string
		toggle     ec2.InstanceMetadataToggle
	}{
		{
			name:       "On",
			instanceID: "12345",
			toggle:     ec2.InstanceMetadataToggleEnabled,
		},
		{
			name:       "Off",
			instanceID: "12345",
			toggle:     ec2.InstanceMetadataToggleDisabled,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := ec2mock.New(t)
			api.On("ModifyInstanceMetadataOptions",
				mock.Anything,
				mock.MatchedBy(func(req *awsec2.ModifyInstanceMetadataOptionsInput) bool {
					return *req.InstanceId == tt.instanceID &&
						req.InstanceMetadataTags == types.InstanceMetadataTagsState(tt.toggle)
				}),
				mock.Anything).Return(&awsec2.ModifyInstanceMetadataOptionsOutput{}, nil)

			client := ec2.NewFromAPI(api)
			err := client.ToggleInstanceMetadataTags(context.Background(), tt.instanceID, tt.toggle)

			assert.NoError(t, err)
			api.AssertExpectations(t)
		})
	}
}

func TestToggleInstanceMetadataTags_Error(t *testing.T) {
	api := ec2mock.New(t)
	api.On("ModifyInstanceMetadataOptions", mock.Anything, mock.Anything, mock.Anything).
		Return(&awsec2.ModifyInstanceMetadataOptionsOutput{}, errors.New("error"))

	client := ec2.NewFromAPI(api)
	err := client.ToggleInstanceMetadataTags(context.Background(), "12345", ec2.InstanceMetadataToggleEnabled)

	assert.Error(t, err)
}
