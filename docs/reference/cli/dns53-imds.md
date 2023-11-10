---
description: "Learn how to use the DNS 53 command line to toggle IMDS features"
icon: material/console
social:
  cards: false
---

# Command Line

Toggle EC2 IMDS features

## Usage

```{ .text .no-select .no-copy }
dns53 imds [flags]
```

## Flags

```{ .text .no-select .no-copy }
-h, --help                            help for imds
    --instance-metadata-tags string   toggle the inclusion of EC2 instance tags
                                      within IMDS (on|off)
```

## Global Flags

```{ .text .no-select .no-copy }
--profile string   the AWS named profile to use when loading credentials
--region string    the AWS region to use when querying AWS
```
