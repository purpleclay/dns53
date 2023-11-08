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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"runtime"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsimds "github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	awsec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	awsr53 "github.com/aws/aws-sdk-go-v2/service/route53"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gobeam/stringy"
	"github.com/purpleclay/dns53/internal/ec2"
	"github.com/purpleclay/dns53/internal/imds"
	"github.com/purpleclay/dns53/internal/r53"
	"github.com/purpleclay/dns53/internal/tui"
	"github.com/spf13/cobra"
)

type BuildDetails struct {
	Version   string `json:"version,omitempty"`
	GitBranch string `json:"git_branch,omitempty"`
	GitCommit string `json:"git_commit,omitempty"`
	Date      string `json:"build_date,omitempty"`
}

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

  # Launch the TUI, automatically creating and attaching to a default
  # PHZ. This will also skip the wizard
  dns53 --auto-attach

  # Launch the TUI with a given domain name
  dns53 --domain-name custom.domain

  # Launch the TUI with a templated domain name
  dns53 --domain-name "{{.IPv4}}.{{.Region}}"`
)

var domainRegex = regexp.MustCompile("[^a-zA-Z0-9-.]+")

type globalOptions struct {
	awsRegion    string
	awsProfile   string
	imdsBindAddr string
}

type options struct {
	phzID      string
	domainName string
	autoAttach bool
	proxy      bool
	proxyPort  int
}

type autoAttachment struct {
	phzID         string
	vpc           string
	region        string
	createdPhz    bool
	associatedPhz bool
}

// Command defines the root DNS 53 cobra command
type Command struct {
	ctx *globalContext
}

// New initialises the root DNS 53 command
func New() *Command {
	return &Command{
		ctx: &globalContext{
			Context: context.Background(),
		},
	}
}

// Execute the DNS 53 command
func (c *Command) Execute(out io.Writer, buildInfo BuildDetails) error {
	globalOpts := &globalOptions{}
	opts := options{}

	// Capture in PreRun lifecycle
	var metadata imds.Metadata

	rootCmd := &cobra.Command{
		Use: "dns53",
		Short: `Dynamic DNS within Amazon Route 53. Expose your EC2 quickly, easily and privately within a Route
53 Private Hosted Zone (PHZ)`,
		Long:          longDesc,
		Example:       examples,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Construct the required AWS Clients and execute any of the provided options
			cfg, err := awsConfig(globalOpts)
			if err != nil {
				return err
			}

			c.ctx.ec2Client = ec2.NewFromAPI(awsec2.NewFromConfig(cfg))
			c.ctx.imdsClient = imds.NewFromAPI(awsimds.NewFromConfig(cfg))
			c.ctx.r53Client = r53.NewFromAPI(awsr53.NewFromConfig(cfg))

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if metadata, err = c.ctx.imdsClient.InstanceMetadata(c.ctx); err != nil {
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
			if opts.autoAttach {
				attachment, err := autoAttachToZone(c.ctx, "dns53", metadata.VPC, metadata.Region)
				if err != nil {
					return err
				}
				opts.phzID = attachment.phzID

				defer removeAttachmentToZone(c.ctx, attachment)
			}

			// At the moment there isn't a nicer way to capture this
			options := tui.Options{
				About: tui.About{
					Name:             "dns53",
					Version:          buildInfo.Version,
					ShortDescription: "Dynamic DNS within Amazon Route 53. Expose your EC2 quickly, easily, and privately.",
				},
				R53Client:    c.ctx.r53Client,
				EC2Metadata:  metadata,
				HostedZoneID: opts.phzID,
				DomainName:   opts.domainName,
				Proxy:        opts.proxy,
				ProxyPort:    opts.proxyPort,
			}

			var err error
			p := tea.NewProgram(
				tui.New(options),
				tea.WithMouseCellMotion(),
				tea.WithOutput(out),
				tea.WithAltScreen(),
			)
			_, err = p.Run()
			return err
		},
	}

	pf := rootCmd.PersistentFlags()
	pf.StringVar(&globalOpts.imdsBindAddr, "imds-bind-addr", "", "the endpoint for all AWS IMDS requests")
	pf.StringVar(&globalOpts.awsProfile, "profile", "", "the AWS named profile to use when loading credentials")
	pf.StringVar(&globalOpts.awsRegion, "region", "", "the AWS region to use when querying AWS")

	// Allow the imds address to be changed at runtime through a hidden flag
	rootCmd.PersistentFlags().MarkHidden("imds-bind-addr")

	f := rootCmd.Flags()
	f.BoolVar(&opts.autoAttach, "auto-attach", false, "automatically create and attach a record set to a default private hosted zone")
	f.StringVar(&opts.domainName, "domain-name", "", "assign a custom domain name when generating a record set")
	f.StringVar(&opts.phzID, "phz-id", "", "an ID of a Route53 private hosted zone to use when generating a record set")
	f.BoolVar(&opts.proxy, "proxy", false, "enable a reverse proxy for tracing requests to this ec2")
	f.IntVar(&opts.proxyPort, "proxy-port", 10080, "the port assigned to the proxy when enabled")

	rootCmd.AddCommand(versionCmd(out, buildInfo))
	rootCmd.AddCommand(manPagesCmd(out))
	rootCmd.AddCommand(imdsCommand())
	rootCmd.AddCommand(tagsCommand(out))

	rootCmd.SetUsageTemplate(customUsageTemplate)
	return rootCmd.ExecuteContext(c.ctx)
}

func awsConfig(opts *globalOptions) (aws.Config, error) {
	var optsFn []func(*config.LoadOptions) error
	if opts.awsProfile != "" {
		optsFn = append(optsFn, config.WithSharedConfigProfile(opts.awsProfile))
	}

	if opts.awsRegion != "" {
		optsFn = append(optsFn, config.WithRegion(opts.awsRegion))
	}

	if opts.imdsBindAddr != "" {
		optsFn = append(optsFn, config.WithEC2IMDSEndpoint(opts.imdsBindAddr))
	}
	return config.LoadDefaultConfig(context.Background(), optsFn...)
}

//nolint:goerr113
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

func autoAttachToZone(ctx *globalContext, name, vpc, region string) (autoAttachment, error) {
	attachment := autoAttachment{
		vpc:    vpc,
		region: region,
	}

	zone, err := ctx.r53Client.ByName(ctx, name)
	if err != nil {
		return attachment, err
	}

	if zone == nil {
		newZone, err := ctx.r53Client.CreatePrivateHostedZone(ctx, "dns53", vpc, region)
		if err != nil {
			return attachment, err
		}

		zone = &newZone

		// Record that this PHZ was created during auto-attachment
		attachment.createdPhz = true
	} else {
		if err := ctx.r53Client.AssociateVPCWithZone(ctx, zone.ID, vpc, region); err != nil {
			return attachment, err
		}

		// An explicit association has been made between the EC2 VPC and the PHZ during auto-attachment
		attachment.associatedPhz = true
	}

	attachment.phzID = zone.ID
	return attachment, nil
}

func removeAttachmentToZone(ctx *globalContext, attach autoAttachment) error {
	if attach.createdPhz {
		return ctx.r53Client.DeletePrivateHostedZone(ctx, attach.phzID)
	}

	return ctx.r53Client.DisassociateVPCWithZone(ctx, attach.phzID, attach.vpc, attach.region)
}

func versionCmd(out io.Writer, buildInfo BuildDetails) *cobra.Command {
	var short bool
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print build time version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			if short {
				fmt.Fprintf(out, buildInfo.Version)
				return nil
			}

			ver := struct {
				Go     string `json:"go"`
				GoArch string `json:"go_arch"`
				GoOS   string `json:"go_os"`
				BuildDetails
			}{
				Go:           runtime.Version(),
				GoArch:       runtime.GOARCH,
				GoOS:         runtime.GOOS,
				BuildDetails: buildInfo,
			}
			return json.NewEncoder(out).Encode(&ver)
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&short, "short", false, "only print the version number")

	return cmd
}
