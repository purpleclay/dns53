//go:build integration

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

package imds_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	awsimds "github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/purpleclay/dns53/pkg/imds"
	aemm "github.com/purpleclay/testcontainers-aemm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstanceMetadata(t *testing.T) {
	ctx := context.Background()
	container := aemm.MustStart(ctx)
	defer container.Terminate(ctx)

	cfg, err := config.LoadDefaultConfig(ctx, config.WithEC2IMDSEndpoint(container.URL))
	require.NoError(t, err)

	client := imds.NewFromAPI(awsimds.NewFromConfig(cfg))
	metadata, err := client.InstanceMetadata(ctx)
	require.NoError(t, err)

	assert.Equal(t, aemm.ValueLocalIPv4, metadata.IPv4)
	assert.Equal(t, aemm.ValuePlacementRegion, metadata.Region)
	assert.Equal(t, aemm.ValueNetworkInterfaces0VPCID, metadata.VPC)
}
