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

package tui_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsr53 "github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/charmbracelet/bubbletea/teatest"
	"github.com/purpleclay/dns53/internal/tui"
	"github.com/purpleclay/dns53/pkg/imds"
	"github.com/purpleclay/dns53/pkg/imdstest"
	"github.com/purpleclay/dns53/pkg/r53"
)

func TestDashboard(t *testing.T) {
	model, _ := tui.Dashboard(tui.DashboardOptions{
		R53Client:  r53.NewFromAPI(R53Stub{}),
		IMDSClient: imds.NewFromAPI(imdstest.Stub{}),
		Version:    "v0.1.0",
		PhzID:      "AAAAAAAAAAAAAAAAAAAAAAAA",
		DomainName: "custom.ec2",
	})

	teatest.TestModel(t,
		model,
		func(p teatest.Program, in io.Writer) {
			// Enforced sleep to ensure Init method retrieves EC2 metadata
			time.Sleep(200 * time.Millisecond)
		},
		func(out []byte) {
			teatest.RequireEqualOutput(t, out)
		})
}

// Stub implementation of R53
type R53Stub struct{}

func (s R53Stub) GetHostedZone(ctx context.Context, params *awsr53.GetHostedZoneInput, optFns ...func(*awsr53.Options)) (*awsr53.GetHostedZoneOutput, error) {
	return &awsr53.GetHostedZoneOutput{
		HostedZone: &types.HostedZone{
			Id:   aws.String("AAAAAAAAAAAAAAAAAAAAAAAA"),
			Name: aws.String("testing"),
		},
	}, nil
}

func (s R53Stub) ListHostedZonesByVPC(ctx context.Context, params *awsr53.ListHostedZonesByVPCInput, optFns ...func(*awsr53.Options)) (*awsr53.ListHostedZonesByVPCOutput, error) {
	return &awsr53.ListHostedZonesByVPCOutput{}, nil
}

func (s R53Stub) ChangeResourceRecordSets(ctx context.Context, params *awsr53.ChangeResourceRecordSetsInput, optFns ...func(*awsr53.Options)) (*awsr53.ChangeResourceRecordSetsOutput, error) {
	return &awsr53.ChangeResourceRecordSetsOutput{}, nil
}
