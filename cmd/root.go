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
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/config"
	awsimds "github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	awsr53 "github.com/aws/aws-sdk-go-v2/service/route53"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/purpleclay/dns53/internal/tui"
	"github.com/purpleclay/dns53/pkg/imds"
	"github.com/purpleclay/dns53/pkg/r53"
	"github.com/spf13/cobra"
)

type options struct {
	region  string
	profile string
	phzID   string
	dnsName string
}

func Execute(out io.Writer) error {
	opts := options{}

	rootCmd := &cobra.Command{
		Use:   "dns53",
		Short: "Dynamic DNS within Amazon Route53. Expose your EC2 quickly, easily and privately",
		RunE: func(cmd *cobra.Command, args []string) error {
			optsFn := []func(*config.LoadOptions) error{}
			if opts.profile != "" {
				optsFn = append(optsFn, config.WithSharedConfigProfile(opts.profile))
			}

			if opts.region != "" {
				optsFn = append(optsFn, config.WithRegion(opts.region))
			}

			optsFn = append(optsFn, config.WithEC2IMDSEndpoint("http://localhost:1338/latest/meta-data/"))

			cfg, err := config.LoadDefaultConfig(context.TODO(), optsFn...)
			if err != nil {
				return err
			}

			model, err := tui.Dashboard(tui.DashboardOptions{
				R53Client:  r53.NewFromAPI(awsr53.NewFromConfig(cfg)),
				IMDSClient: imds.NewFromAPI(awsimds.NewFromConfig(cfg)),
				Version:    version,
				PhzID:      opts.phzID,
				DNSName:    opts.dnsName,
			})
			if err != nil {
				return err
			}

			return tea.NewProgram(model, tea.WithAltScreen()).Start()
		},
	}

	f := rootCmd.Flags()
	f.StringVar(&opts.region, "region", "", "the AWS region to use when querying AWS")
	f.StringVar(&opts.profile, "profile", "", "the AWS named profile to use when loading credentials")
	f.StringVar(&opts.phzID, "phz-id", "", "an ID of a Route53 private hosted zone to use when generating a record set")
	f.StringVar(&opts.dnsName, "dns-name", "", "use the DNS name when generating a record set")

	rootCmd.AddCommand(newVersionCmd(out))
	rootCmd.AddCommand(newManPagesCmd(out))
	rootCmd.AddCommand(newCompletionCmd(out))

	return rootCmd.ExecuteContext(context.Background())
}
