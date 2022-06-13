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

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/purpleclay/dns53/pkg/ec2"
	"github.com/spf13/cobra"
)

func newExposeCmd(out io.Writer) *cobra.Command {
	expCmd := &cobra.Command{
		Use:                   "expose",
		Short:                 "Generate a DNS A-Record and privately expose this EC2",
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: create a connection to AWS
			// TODO: create model
			// TODO: create bubbletea program

			cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEC2IMDSEndpoint("http://localhost:1338"))
			if err != nil {
				return err
			}

			meta, err := ec2.InstanceMetadata(cfg)
			fmt.Printf("%#v\n", meta)
			return err
		},
	}

	return expCmd
}
