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

package cmd

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsr53 "github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/purpleclay/dns53/internal/imds"
	"github.com/purpleclay/dns53/internal/r53"
	"github.com/purpleclay/dns53/internal/r53/r53mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// func TestAWSConfig(t *testing.T) {
// 	cfg, err := awsConfig(&globalOptions{
// 		AWSRegion:  "eu-west-2",
// 		AWSProfile: "testing",
// 	})

// 	require.NoError(t, err)
// 	assert.Equal(t, "eu-west-2", cfg.Region)
// }

func TestResolveDomainName(t *testing.T) {
	metadata := imds.Metadata{
		Name: "my-ec2",
	}

	tests := []struct {
		name     string
		domain   string
		expected string
	}{
		{
			name:     "NoTemplating",
			domain:   "custom.domain",
			expected: "custom.domain",
		},
		{
			name:     "WithNameField",
			domain:   "custom.{{.Name}}",
			expected: "custom.my-ec2",
		},
		{
			name:     "WithNameFieldSpaces",
			domain:   "custom.{{ .Name }}",
			expected: "custom.my-ec2",
		},
		{
			name:     "ReplacesDoubleHyphens",
			domain:   "another--custom.domain",
			expected: "another-custom.domain",
		},
		{
			name:     "ReplacesDoubleDots",
			domain:   "my-custom123..domain",
			expected: "my-custom123.domain",
		},
		{
			name:     "RemoveLeadingTrailingHyphen",
			domain:   "-this-is-a-custom.domain-",
			expected: "this-is-a-custom.domain",
		},
		{
			name:     "RemoveLeadingTrailingDot",
			domain:   ".a-custom.domain.",
			expected: "a-custom.domain",
		},
		{
			name:     "TrimUnsupportedCharacters",
			domain:   "custom@#.doma**in-123",
			expected: "custom.domain-123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain, err := resolveDomainName(tt.domain, metadata)

			require.NoError(t, err)
			require.Equal(t, tt.expected, domain)
		})
	}
}

func TestResolveDomainNameNoInstanceTags(t *testing.T) {
	_, err := resolveDomainName("custom.{{.Name}}", imds.Metadata{})

	assert.EqualError(t, err, `to use metadata within a custom domain name, please enable IMDS instance tags support
for your EC2 instance:

  $ dns53 imds --instance-metadata-tags on

Or read the official AWS documentation at:
https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Tags.html#allow-access-to-tags-in-IMDS`)
}

func TestResolveDomainNameTransformsNameTagToKebabCase(t *testing.T) {
	domain, err := resolveDomainName("first.custom.{{.Name}}", imds.Metadata{Name: "MyEc2 123"})

	require.NoError(t, err)
	assert.Equal(t, "first.custom.my-ec2-123", domain)
}

func TestResolveDomainNameStripsLeadingTrailingHyphenFromNameTag(t *testing.T) {
	domain, err := resolveDomainName("second.custom.{{.Name}}", imds.Metadata{Name: "-MyEc2 123-"})

	require.NoError(t, err)
	assert.Equal(t, "second.custom.my-ec2-123", domain)
}

func TestResolveDomainNameInvalidGoTemplate(t *testing.T) {
	_, err := resolveDomainName("custom.{{.Name}", imds.Metadata{Name: "MyEc2 123"})

	assert.Error(t, err)
}

func TestResolveDomainNameUnrecognisedTemplateFields(t *testing.T) {
	_, err := resolveDomainName("custom.{{.Unknown}}", imds.Metadata{})

	assert.Error(t, err)
}

func TestCleanTagsAppendsToMap(t *testing.T) {
	tags := map[string]string{
		"My+@-key_=,.:1": "A value",
	}
	cleanTags(tags)

	expected := map[string]string{
		"My+@-key_=,.:1": "a-value",
		"MyKey1":         "a-value",
	}
	for k, v := range expected {
		assert.Contains(t, tags, k)
		assert.Equal(t, v, tags[k])
	}
}

func TestAutoAttachToZone(t *testing.T) {
	m := r53mock.New(t)
	m.On("ListHostedZonesByName", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.ListHostedZonesByNameOutput{}, nil)
	m.On("CreateHostedZone", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.CreateHostedZoneOutput{
		HostedZone: &types.HostedZone{
			Id:   aws.String("/hostedzone/Z00000000001"),
			Name: aws.String("dns53."),
		},
	}, nil)

	attachment, err := autoAttachToZone(r53.NewFromAPI(m), "dns53", "vpc-12345", "eu-west-2")

	require.NoError(t, err)
	assert.Equal(t, "Z00000000001", attachment.phzID)
	assert.True(t, attachment.createdPhz)
	assert.False(t, attachment.associatedPhz)
}

func TestAutoAttachToZoneExisting(t *testing.T) {
	m := r53mock.New(t)
	m.On("ListHostedZonesByName", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.ListHostedZonesByNameOutput{
		HostedZones: []types.HostedZone{
			{
				Id:   aws.String("/hostedzone/Z00000000002"),
				Name: aws.String("dns53"),
				Config: &types.HostedZoneConfig{
					PrivateZone: true,
				},
			},
		},
	}, nil)
	m.On("AssociateVPCWithHostedZone", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.AssociateVPCWithHostedZoneOutput{}, nil)

	attachment, err := autoAttachToZone(r53.NewFromAPI(m), "dns53", "vpc-12345", "eu-west-2")

	require.NoError(t, err)
	assert.Equal(t, "Z00000000002", attachment.phzID)
	assert.False(t, attachment.createdPhz)
	assert.True(t, attachment.associatedPhz)
}

func TestAutoAttachToZoneSearchError(t *testing.T) {
	errMsg := "failed to search"

	m := r53mock.New(t)
	m.On("ListHostedZonesByName", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.ListHostedZonesByNameOutput{}, errors.New(errMsg))

	_, err := autoAttachToZone(r53.NewFromAPI(m), "dns53", "vpc-12345", "eu-west-2")

	require.EqualError(t, err, errMsg)
	m.AssertNotCalled(t, "AssociateVPCWithHostedZone")
}

func TestAutoAttachToZoneCreationError(t *testing.T) {
	errMsg := "failed to create"

	m := r53mock.New(t)
	m.On("ListHostedZonesByName", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.ListHostedZonesByNameOutput{}, nil)
	m.On("CreateHostedZone", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.CreateHostedZoneOutput{}, errors.New(errMsg))

	_, err := autoAttachToZone(r53.NewFromAPI(m), "dns53", "vpc-12345", "eu-west-2")

	require.EqualError(t, err, errMsg)
}

func TestAutoAttachToZoneAssociationError(t *testing.T) {
	errMsg := "failed to associate"

	m := r53mock.New(t)
	m.On("ListHostedZonesByName", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.ListHostedZonesByNameOutput{
		HostedZones: []types.HostedZone{
			{
				Id:   aws.String("/hostedzone/Z00000000003"),
				Name: aws.String("dns53"),
				Config: &types.HostedZoneConfig{
					PrivateZone: true,
				},
			},
		},
	}, nil)
	m.On("AssociateVPCWithHostedZone", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.AssociateVPCWithHostedZoneOutput{}, errors.New(errMsg))

	_, err := autoAttachToZone(r53.NewFromAPI(m), "dns53", "vpc-12345", "eu-west-2")

	require.EqualError(t, err, errMsg)
}
