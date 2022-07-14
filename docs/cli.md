# Command Line

## Usage

```sh
dns53 [options]
```

## Options

```sh
$ dns53 --help

Dynamic DNS within Amazon Route 53. Expose your EC2 quickly, easily and
privately within a Route 53 Private Hosted Zone (PHZ).

Your EC2 will be exposed through a dynamically generated resource record that
will automatically be deleted when dns53 exits. Let dns53 name your resource
record for you, or customise it to your needs.

Built using Bubbletea 🧋

Usage:
  dns53 [flags]
  dns53 [command]

Examples:
  # Launch the TUI and use the wizard to select a PHZ
  dns53

  # Launch the TUI using a chosen PHZ, effectively skipping the wizard
  dns53 --phz-id Z000000000ABCDEFGHIJK

  # Launch the TUI with a given domain name
  dns53 --domain-name custom.domain

  # Launch the TUI with a templated domain name
  dns53 --domain-name "{{.IPv4}}.{{.Region}}"

Available Commands:
  completion  Generate completion script for your target shell
  help        Help about any command
  version     Prints the build time version information

Flags:
      --domain-name string   assign a custom domain name when generating a
                             record set
  -h, --help                 help for dns53
      --phz-id string        an ID of a Route53 private hosted zone to use when
                             generating a record set
      --profile string       the AWS named profile to use when loading
                             credentials
      --region string        the AWS region to use when querying AWS

Use "dns53 [command] --help" for more information about a command.
```
