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
	"strings"

	"github.com/purpleclay/dns53/internal/ec2"
	"github.com/purpleclay/dns53/internal/imds"
	"github.com/spf13/cobra"
)

// Custom type used to toggle any setting "on" or "off"
type toggleSetting string

const (
	toggleSettingOn  toggleSetting = "on"
	toggleSettingOff toggleSetting = "off"
)

func (t *toggleSetting) String() string {
	return string(*t)
}

func (t *toggleSetting) Set(v string) error {
	setting := strings.ToLower(v)

	switch setting {
	case "on", "off":
		*t = toggleSetting(setting)
		return nil
	default:
		return errors.New(`supported values are "on" or "off" (case-insensitive)`)
	}
}

func (t *toggleSetting) Type() string {
	return "string"
}

type imdsOptions struct {
	InstanceMetadataTags toggleSetting
}

func newIMDSCommand() *cobra.Command {
	opt := imdsOptions{}

	cmd := &cobra.Command{
		Use:           "imds",
		Short:         "Toggle EC2 IMDS features",
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context().(*globalContext)

			return toggleMetadataTags(ctx.ec2Client, ctx.imdsClient, opt.InstanceMetadataTags)
		},
	}

	f := cmd.Flags()
	f.Var(&opt.InstanceMetadataTags, "instance-metadata-tags", "toggle the inclusion of EC2 instance tags within IMDS (on|off)")

	cmd.MarkFlagRequired("instance-metadata-tags")
	return cmd
}

func toggleMetadataTags(ec2Client *ec2.Client, imdsClient *imds.Client, setting toggleSetting) error {
	metadata, err := imdsClient.InstanceMetadata(context.Background())
	if err != nil {
		return err
	}

	var toggle ec2.InstanceMetadataToggle

	switch setting {
	case toggleSettingOn:
		toggle = ec2.InstanceMetadataToggleEnabled
	case toggleSettingOff:
		toggle = ec2.InstanceMetadataToggleDisabled
	}

	return ec2Client.ToggleInstanceMetadataTags(context.Background(), metadata.InstanceID, toggle)
}
