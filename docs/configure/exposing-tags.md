---
description: "Expand your domain name customisation options by exposing your EC2 instance tags"
icon: material/tag-outline
status: new
---

# Exposing EC2 Instance Tags

By default, EC2 tags are not accessible through the Instance Metadata Service (IMDS) and subsequently by `dns53`. Granting access to EC2 instance tags can be carried out manually[^1] or with the following custom command:

```sh
dns53 imds --instance-metadata-tags on
```

## Cleaning Tag Names

All tags accessible through IMDS will subsequently be "cleaned" and made available for crafting custom domain names. A Pascal Case naming convention is applied to all tags when stored within an internal map alongside their originally named counterpart, and both are accessible through [templating](../reference/templating.md#dynamic-tags).

[^1]: Access to EC2 instance tags can be granted directly through the AWS Console or by using the CLI as documented [here](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Tags.html#allow-access-to-tags-in-IMDS)
