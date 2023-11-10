---
description: "Learn how to use the DNS 53 command line for privately broadcasting your EC2 instance"
icon: material/console
social:
  cards: false
status: new
---

# Command Line

Dynamic DNS within Amazon Route 53. Expose your EC2 quickly, easily, and privately within a Route 53 Private Hosted Zone (PHZ).

Your EC2 will be exposed through a dynamically generated resource record that will automatically be deleted when dns53 exits. Let dns53 name your resource record for you, or customise it to your needs.

## Usage

```{ .text .no-select .no-copy }
dns53 [flags]
dns53 [command]
```

## Flags

```{ .text .no-select .no-copy }
    --auto-attach          automatically create and attach a record set to a
                           default private hosted zone
    --domain-name string   assign a custom domain name when generating a record
                           set
-h, --help                 help for dns53
    --phz-id string        an ID of a Route53 private hosted zone to use when
                           generating a record set
    --profile string       the AWS named profile to use when loading credentials
    --proxy                enable a reverse proxy for tracing requests to this
                           ec2
    --proxy-port int       the port assigned to the proxy when enabled
                           (default 10080)
    --region string        the AWS region to use when querying AWS
```

## Commands

```{ .text .no-select .no-copy }
completion  Generate the autocompletion script for the specified shell
help        Help about any command
imds        Toggle EC2 IMDS features
tags        Lists all available EC2 instance tags and how to use them with Go
            templating
version     Print build time version information
```
