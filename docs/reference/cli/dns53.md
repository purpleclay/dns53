---
description: "Learn how to use the DNS 53 command line for privately broadcasting your EC2 instance"
icon: material/console
---

# Command Line

Dynamic DNS within Amazon Route 53. Expose your EC2 quickly, easily, and privately within a Route 53 Private Hosted Zone (PHZ).

Your EC2 will be exposed through a dynamically generated resource record that will automatically be deleted when dns53 exits. Let dns53 name your resource record for you, or customise it to your needs.

## Usage

```text
dns53 [flags]
dns53 [command]
```

## Flags

```text
    --domain-name string   assign a custom domain name when generating a record set
-h, --help                 help for dns53
    --phz-id string        an ID of a Route53 private hosted zone to use when generating a record set
    --profile string       the AWS named profile to use when loading credentials
    --region string        the AWS region to use when querying AWS
```

## Commands

```text
completion  Generate completion script for your target shell
help        Help about any command
imds        Toggle EC2 IMDS features
version     Prints the build time version information
```
