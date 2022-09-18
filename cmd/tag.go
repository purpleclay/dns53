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
	"fmt"
	"io"

	awsimds "github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/olekukonko/tablewriter"
	"github.com/purpleclay/dns53/internal/imds"
	"github.com/spf13/cobra"
)

func newTagsCommand(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "tags",
		Short:         "Lists all available EC2 instance tags and how to use them with Go templating",
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := awsConfig(globalOpts)
			if err != nil {
				return err
			}

			imdsClient := imds.NewFromAPI(awsimds.NewFromConfig(cfg))
			metadata, err := imdsClient.InstanceMetadata(context.Background())
			if err != nil {
				return err
			}

			table := tablewriter.NewWriter(out)
			table.SetHeader([]string{"Tag", "Value", "Property Chaining", "Indexed"})

			for k, v := range metadata.Tags {
				cleanedTag, cleanedValue := cleanTag(k, v)

				table.Append([]string{
					k,
					cleanedValue,
					fmt.Sprintf("{{.Tags.%s}}", cleanedTag),
					fmt.Sprintf("{{index .Tags \"%s\"}}", k),
				})
			}

			table.Render()
			return nil
		},
	}

	return cmd
}
