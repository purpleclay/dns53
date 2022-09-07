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
	"errors"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsimds "github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	awsr53 "github.com/aws/aws-sdk-go-v2/service/route53"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/purpleclay/dns53/internal/imds"
	"github.com/purpleclay/dns53/internal/r53"
	"github.com/purpleclay/dns53/internal/tui"
	"github.com/spf13/cobra"
)

const (
	longDesc = `Dynamic DNS within Amazon Route 53. Expose your EC2 quickly, easily and privately within a Route 
53 Private Hosted Zone (PHZ).

Your EC2 will be exposed through a dynamically generated resource record that will automatically 
be deleted when dns53 exits. Let dns53 name your resource record for you, or customise it to your needs. 

Built using Bubbletea 🧋`
	examples = `  # Launch the TUI and use the wizard to select a PHZ
  dns53

  # Launch the TUI using a chosen PHZ, effectively skipping the wizard
  dns53 --phz-id Z000000000ABCDEFGHIJK

  # Launch the TUI with a given domain name
  dns53 --domain-name custom.domain

  # Launch the TUI with a templated domain name
  dns53 --domain-name "{{.IPv4}}.{{.Region}}"`
)

// Global options set through persistent flags
type globalOptions struct {
	AWSRegion  string
	AWSProfile string
}

var globalOpts = &globalOptions{}

type options struct {
	phzID      string
	domainName string
}

func Execute(out io.Writer) error {
	opts := options{}

	rootCmd := &cobra.Command{
		Use: "dns53",
		Short: `Dynamic DNS within Amazon Route 53. Expose your EC2 quickly, easily and privately within a Route 
53 Private Hosted Zone (PHZ)`,
		Long:    longDesc,
		Example: examples,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := awsConfig(globalOpts)
			if err != nil {
				return err
			}

			// If a custom domain name has been provided, check that it can be resolved from IMDS
			imdsClient := imds.NewFromAPI(awsimds.NewFromConfig(cfg))

			if opts.domainName != "" {
				if err := domainNameSupported(opts.domainName, imdsClient); err != nil {
					return err
				}
			}

			model, err := tui.Dashboard(tui.DashboardOptions{
				R53Client:  r53.NewFromAPI(awsr53.NewFromConfig(cfg)),
				IMDSClient: imdsClient,
				Version:    version,
				PhzID:      opts.phzID,
				DomainName: opts.domainName,
			})
			if err != nil {
				return err
			}

			return tea.NewProgram(model, tea.WithAltScreen()).Start()
		},
	}

	pf := rootCmd.PersistentFlags()
	pf.StringVar(&globalOpts.AWSRegion, "region", "", "the AWS region to use when querying AWS")
	pf.StringVar(&globalOpts.AWSProfile, "profile", "", "the AWS named profile to use when loading credentials")

	f := rootCmd.Flags()
	f.StringVar(&opts.phzID, "phz-id", "", "an ID of a Route53 private hosted zone to use when generating a record set")
	f.StringVar(&opts.domainName, "domain-name", "", "assign a custom domain name when generating a record set")

	rootCmd.AddCommand(newVersionCmd(out))
	rootCmd.AddCommand(newManPagesCmd(out))
	rootCmd.AddCommand(newCompletionCmd(out))
	rootCmd.AddCommand(newIMDSCommand(out))

	return rootCmd.ExecuteContext(context.Background())
}

func awsConfig(opts *globalOptions) (aws.Config, error) {
	optsFn := []func(*config.LoadOptions) error{}
	if opts.AWSProfile != "" {
		optsFn = append(optsFn, config.WithSharedConfigProfile(opts.AWSProfile))
	}

	if opts.AWSRegion != "" {
		optsFn = append(optsFn, config.WithRegion(opts.AWSRegion))
	}

	return config.LoadDefaultConfig(context.Background(), optsFn...)
}

func domainNameSupported(domain string, imdsClient *imds.Client) error {
	dmn := strings.ReplaceAll(domain, " ", "")
	if strings.Contains(dmn, "{{.Name}}") {
		metadata, err := imdsClient.InstanceMetadata(context.Background())
		if err != nil {
			return err
		}

		if metadata.Name == "" {
			return errors.New(`to use metadata within a custom domain name, please enable IMDS instance tags support 
for your EC2 instance:

  $ dns53 imds --instance-metadata-tags on

Or read the official AWS documentation at: 
https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Tags.html#allow-access-to-tags-in-IMDS`)
		}
	}
	return nil
}
