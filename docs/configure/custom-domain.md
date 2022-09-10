---
icon: material/web
---

# Using a Custom Domain Name

If you want complete control of the domain name associated with your EC2, you can customise it in one of two ways.

!!! tip "Route53 Root Domain is Optional"

    `dns53` will automatically append the Route53 root domain when creating the A-Record. Feel free to omit this when providing a custom domain

## Static Domain

```sh
dns53 --domain-name "my.ec2"
```

## Templated Domain

A templated domain leverages the text templating capabilities of the Go language to replace named fields with concrete values. A list of supported named fields can be found [here](../reference/templating.md).

```sh
dns53 --domain-name "{{.IPv4}}.{{.Region}}"
```
