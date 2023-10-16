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

package cmd

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsr53 "github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/purpleclay/dns53/internal/imds"
	"github.com/purpleclay/dns53/internal/imds/imdsstub"
	"github.com/purpleclay/dns53/internal/r53"
	"github.com/purpleclay/dns53/internal/r53/r53mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

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

func TestRootCommand(t *testing.T) {
	options := []globalContextOption{
		withIMDSClient(imds.NewFromAPI(imdsstub.New(t))),
		withSkipTea(),
	}

	cmd := newWithOptions(options...)
	err := cmd.Execute([]string{})

	require.NoError(t, err)
	assert.Equal(t, cmd.ctx.teaModelOptions.HostedZoneID, "")
	assert.Equal(t, cmd.ctx.teaModelOptions.DomainName, "")
}

func TestRootCommandWithPrivateHostedZoneID(t *testing.T) {
	options := []globalContextOption{
		withIMDSClient(imds.NewFromAPI(imdsstub.New(t))),
		withSkipTea(),
	}

	cmd := newWithOptions(options...)
	err := cmd.Execute([]string{"--phz-id", "Z00000000001"})

	require.NoError(t, err)
	assert.Equal(t, cmd.ctx.teaModelOptions.HostedZoneID, "Z00000000001")
	assert.Equal(t, cmd.ctx.teaModelOptions.DomainName, "")
}

func TestRootCommandWithCustomDomain(t *testing.T) {
	options := []globalContextOption{
		withIMDSClient(imds.NewFromAPI(imdsstub.New(t))),
		withSkipTea(),
	}

	cmd := newWithOptions(options...)
	err := cmd.Execute([]string{"--domain-name", "custom.{{.Name}}"})

	require.NoError(t, err)
	assert.Equal(t, cmd.ctx.teaModelOptions.HostedZoneID, "")
	assert.Equal(t, cmd.ctx.teaModelOptions.DomainName, "custom.stub-ec2")
}

func TestRootCommandAutoAttachToZone(t *testing.T) {
	m := r53mock.New(t)
	m.On("ListHostedZonesByName", mock.Anything, mock.MatchedBy(func(req *awsr53.ListHostedZonesByNameInput) bool {
		return *req.DNSName == "dns53"
	}), mock.Anything).Return(&awsr53.ListHostedZonesByNameOutput{}, nil)
	m.On("CreateHostedZone", mock.Anything, mock.MatchedBy(func(req *awsr53.CreateHostedZoneInput) bool {
		return true
	}), mock.Anything).Return(&awsr53.CreateHostedZoneOutput{
		HostedZone: &types.HostedZone{
			Id:   aws.String("/hostedzone/Z00000000002"),
			Name: aws.String("dns53."),
		},
	}, nil)
	m.On("DeleteHostedZone", mock.Anything, mock.MatchedBy(func(req *awsr53.DeleteHostedZoneInput) bool {
		return *req.Id == "Z00000000002"
	}), mock.Anything).Return(&awsr53.DeleteHostedZoneOutput{}, nil)

	// Configure the command to run in test mode
	options := []globalContextOption{
		withIMDSClient(imds.NewFromAPI(imdsstub.New(t))),
		withR53Client(r53.NewFromAPI(m)),
		withSkipTea(),
	}

	cmd := newWithOptions(options...)
	err := cmd.Execute([]string{"--auto-attach"})

	require.NoError(t, err)
}

func TestRootCommandAutoAttachToZoneExisting(t *testing.T) {
	m := r53mock.New(t)
	m.On("ListHostedZonesByName", mock.Anything, mock.MatchedBy(func(req *awsr53.ListHostedZonesByNameInput) bool {
		return *req.DNSName == "dns53"
	}), mock.Anything).Return(&awsr53.ListHostedZonesByNameOutput{
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
	m.On("AssociateVPCWithHostedZone", mock.Anything, mock.MatchedBy(func(req *awsr53.AssociateVPCWithHostedZoneInput) bool {
		return *req.HostedZoneId == "Z00000000003"
	}), mock.Anything).Return(&awsr53.AssociateVPCWithHostedZoneOutput{}, nil)
	m.On("DisassociateVPCFromHostedZone", mock.Anything, mock.MatchedBy(func(req *awsr53.DisassociateVPCFromHostedZoneInput) bool {
		return *req.HostedZoneId == "Z00000000003"
	}), mock.Anything).Return(&awsr53.DisassociateVPCFromHostedZoneOutput{}, nil)

	// Configure the command to run in test mode
	options := []globalContextOption{
		withIMDSClient(imds.NewFromAPI(imdsstub.New(t))),
		withR53Client(r53.NewFromAPI(m)),
		withSkipTea(),
	}

	cmd := newWithOptions(options...)
	err := cmd.Execute([]string{"--auto-attach"})

	require.NoError(t, err)
}

//nolint:goerr113
func TestRootCommandAutoAttachToZoneSearchError(t *testing.T) {
	errMsg := "failed to search"

	m := r53mock.New(t)
	m.On("ListHostedZonesByName", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.ListHostedZonesByNameOutput{}, errors.New(errMsg))

	// Configure the command to run in test mode
	options := []globalContextOption{
		withIMDSClient(imds.NewFromAPI(imdsstub.New(t))),
		withR53Client(r53.NewFromAPI(m)),
		withSkipTea(),
	}

	cmd := newWithOptions(options...)
	err := cmd.Execute([]string{"--auto-attach"})

	require.EqualError(t, err, errMsg)
	m.AssertNotCalled(t, "AssociateVPCWithHostedZone")
}

//nolint:goerr113
func TestRootCommandAutoAttachToZoneCreationError(t *testing.T) {
	errMsg := "failed to create"

	m := r53mock.New(t)
	m.On("ListHostedZonesByName", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.ListHostedZonesByNameOutput{}, nil)
	m.On("CreateHostedZone", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.CreateHostedZoneOutput{}, errors.New(errMsg))

	// Configure the command to run in test mode
	options := []globalContextOption{
		withIMDSClient(imds.NewFromAPI(imdsstub.New(t))),
		withR53Client(r53.NewFromAPI(m)),
		withSkipTea(),
	}

	cmd := newWithOptions(options...)
	err := cmd.Execute([]string{"--auto-attach"})

	require.EqualError(t, err, errMsg)
}

//nolint:goerr113
func TestRootCommandAutoAttachToZoneAssociationError(t *testing.T) {
	errMsg := "failed to associate"

	m := r53mock.New(t)
	m.On("ListHostedZonesByName", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.ListHostedZonesByNameOutput{
		HostedZones: []types.HostedZone{
			{
				Id:   aws.String("/hostedzone/Z00000000004"),
				Name: aws.String("dns53"),
				Config: &types.HostedZoneConfig{
					PrivateZone: true,
				},
			},
		},
	}, nil)
	m.On("AssociateVPCWithHostedZone", mock.Anything, mock.Anything, mock.Anything).Return(&awsr53.AssociateVPCWithHostedZoneOutput{}, errors.New(errMsg))

	// Configure the command to run in test mode
	options := []globalContextOption{
		withIMDSClient(imds.NewFromAPI(imdsstub.New(t))),
		withR53Client(r53.NewFromAPI(m)),
		withSkipTea(),
	}

	cmd := newWithOptions(options...)
	err := cmd.Execute([]string{"--auto-attach"})

	require.EqualError(t, err, errMsg)
}
