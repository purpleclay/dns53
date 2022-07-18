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

package imdstest

import (
	"context"
	"fmt"
	"io"
	"strings"

	awsimds "github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/purpleclay/dns53/pkg/imds"
	aemm "github.com/purpleclay/testcontainers-aemm"
)

// Stub implementation of the imds.MetadataClientAPI that ensures a
// consistent set of metadata values are always returned
type Stub struct{}

// GetMetadata will inspect the provided instance category and return a
// hardcoded value from the aemm testcontainers library:
//
// https://github.com/purpleclay/testcontainers-aemm
func (s Stub) GetMetadata(ctx context.Context, params *awsimds.GetMetadataInput, optFns ...func(*awsimds.Options)) (*awsimds.GetMetadataOutput, error) {
	var output string
	switch params.Path {
	case imds.PathIPv4:
		output = aemm.ValueLocalIPv4
	case imds.PathMacAddress:
		output = aemm.ValueMAC
	case fmt.Sprintf("network/interfaces/macs/%s/vpc-id", aemm.ValueMAC):
		output = aemm.ValueNetworkInterfaces0VPCID
	case imds.PathInstanceID:
		output = aemm.ValueInstanceID
	case imds.PathPlacementAZ:
		output = aemm.ValuePlacementAvailabilityZone
	case imds.PathPlacementRegion:
		output = aemm.ValuePlacementRegion
	}

	return &awsimds.GetMetadataOutput{
		Content: io.NopCloser(strings.NewReader(output)),
	}, nil
}
