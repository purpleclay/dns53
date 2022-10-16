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
	"bytes"
	"context"
	"errors"
	"io"
	"regexp"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsimds "github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gobeam/stringy"
	"github.com/purpleclay/dns53/internal/imds"
	"github.com/purpleclay/dns53/internal/tui"
	"github.com/spf13/cobra"
)

const (
	longDesc = `Dynamic DNS within Amazon Route 53. Expose your EC2 quickly, easily, and privately within a Route 
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

var (
	globalOpts = &globalOptions{}

	domainRegex = regexp.MustCompile("[^a-zA-Z0-9-.]+")
)

type options struct {
	phzID      string
	domainName string
}

func Execute(out io.Writer) error {
	opts := options{}

	// Capture in PreRun lifecycle
	var cfg aws.Config
	var metadata imds.Metadata

	rootCmd := &cobra.Command{
		Use: "dns53",
		Short: `Dynamic DNS within Amazon Route 53. Expose your EC2 quickly, easily and privately within a Route 
53 Private Hosted Zone (PHZ)`,
		Long:          longDesc,
		Example:       examples,
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			cfg, err = awsConfig(globalOpts)
			if err != nil {
				return err
			}

			imdsClient := imds.NewFromAPI(awsimds.NewFromConfig(cfg))
			if metadata, err = imdsClient.InstanceMetadata(context.Background()); err != nil {
				return err
			}

			if opts.domainName == "" {
				return nil
			}

			cleanTags(metadata.Tags)
			opts.domainName, err = resolveDomainName(opts.domainName, metadata)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// model := tui.Dashboard(tui.DashboardOptions{
			// 	R53Client:  r53.NewFromAPI(awsr53.NewFromConfig(cfg)),
			// 	Metadata:   metadata,
			// 	Version:    version,
			// 	PhzID:      opts.phzID,
			// 	DomainName: opts.domainName,
			// })

			model := tui.New(tui.Options{
				About: tui.About{
					Name:             "dns53",
					Version:          version,
					ShortDescription: "Dynamic DNS within Amazon Route 53. Expose your EC2 quickly, easily, and privately.",
				},
			})

			return tea.NewProgram(model, tea.WithAltScreen()).Start()
		},
	}

	pf := rootCmd.PersistentFlags()
	pf.StringVar(&globalOpts.AWSRegion, "region", "", "the AWS region to use when querying AWS")
	pf.StringVar(&globalOpts.AWSProfile, "profile", "", "the AWS named profile to use when loading credentials")

	f := rootCmd.Flags()
	f.StringVar(&opts.domainName, "domain-name", "", "assign a custom domain name when generating a record set")
	f.StringVar(&opts.phzID, "phz-id", "", "an ID of a Route53 private hosted zone to use when generating a record set")

	rootCmd.AddCommand(newVersionCmd(out))
	rootCmd.AddCommand(newManPagesCmd(out))
	rootCmd.AddCommand(newCompletionCmd(out))
	rootCmd.AddCommand(newIMDSCommand(out))
	rootCmd.AddCommand(newTagsCommand(out))

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

	optsFn = append(optsFn, config.WithEC2IMDSEndpoint("http://localhost:1338/latest/meta-data/"))

	return config.LoadDefaultConfig(context.Background(), optsFn...)
}

func resolveDomainName(domain string, metadata imds.Metadata) (string, error) {
	dmn := strings.ReplaceAll(domain, " ", "")

	if strings.Contains(dmn, "{{.Name}}") {
		if metadata.Name == "" {
			return "", errors.New(`to use metadata within a custom domain name, please enable IMDS instance tags support
for your EC2 instance:

  $ dns53 imds --instance-metadata-tags on

Or read the official AWS documentation at:
https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Tags.html#allow-access-to-tags-in-IMDS`)
		}

		name := stringy.New(metadata.Name)
		metadata.Name = name.KebabCase().ToLower()
	}

	// Sanitise the copy of the metadata before resolving the template
	metadata.IPv4 = strings.ReplaceAll(metadata.IPv4, ".", "-")

	// Execute the domain template
	tmpl, err := template.New("domain").Parse(domain)
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	if err := tmpl.Execute(&out, metadata); err != nil {
		return "", err
	}
	dmn = out.String()

	// Final tidy up of the domain
	dmn = strings.ReplaceAll(dmn, "--", "-")
	dmn = strings.ReplaceAll(dmn, "..", ".")
	dmn = strings.Trim(dmn, "-")
	dmn = strings.Trim(dmn, ".")
	dmn = domainRegex.ReplaceAllString(dmn, "")

	return dmn, nil
}
