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
	"testing"

	"github.com/purpleclay/dns53/internal/imds"
	"github.com/purpleclay/dns53/internal/imds/imdsstub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAWSConfig(t *testing.T) {
	cfg, err := awsConfig(&globalOptions{
		AWSRegion:  "eu-west-2",
		AWSProfile: "testing",
	})

	require.NoError(t, err)
	assert.Equal(t, "eu-west-2", cfg.Region)
}

func TestDomainNameSupported(t *testing.T) {
	tests := []struct {
		name   string
		domain string
	}{
		{
			name:   "NoTemplating",
			domain: "custom.domain",
		},
		{
			name:   "WithNameField",
			domain: "custom.{{.Name}}",
		},
		{
			name:   "WithNameFieldSpaces",
			domain: "custom.{{ .Name }}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domainNameSupported(tt.domain, imds.NewFromAPI(imdsstub.New(t)))

			require.NoError(t, err)
		})
	}
}

func TestDomainNameSupported_NoInstanceTags(t *testing.T) {
	err := domainNameSupported("custom.{{.Name}}", imds.NewFromAPI(imdsstub.NewWithoutTags(t)))

	assert.EqualError(t, err, `to use metadata within a custom domain name, please enable IMDS instance tags support 
for your EC2 instance:

  $ dns53 imds --instance-metadata-tags on

Or read the official AWS documentation at: 
https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Tags.html#allow-access-to-tags-in-IMDS`)
}
