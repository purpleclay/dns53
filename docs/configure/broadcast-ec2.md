---
description: "Broadcasting your EC2 privately within your VPC couldn't be easier"
icon: material/bullhorn-variant-outline
---

# Privately Broadcast your EC2

To broadcast your EC2 privately within your VPC couldn't be easier. Launch the wizard and follow the on-screen prompts:

```sh
dns53
```

## Default Domain Name

A default domain name will be assigned to your EC2 when `dns53` adds an A-Record to the chosen Route53 Private Hosted Zone (PHZ).

It follows the format:

`<EC2_PRIVATE_IPv4>.dns53.<PHZ_ROOT_DOMAIN>` ~> `10-0-1-182.dns53.testing`

## Skipping the Wizard

If you have the ID of your Route53 PHZ handy, you can skip the wizard and immediately broadcast your EC2:

```sh
dns53 --phz-id Z05504861FO8RFR02KU72
```

<div>
    <video controls>
        <source src="../../static/dns53-phzid.webm" type="video/webm">
        <source src="../../static/dns53-phzid.mp4" type="video/mp4">
    </video>
    <sub>recorded using <a href="https://github.com/charmbracelet/vhs" target="_blank">VHS</a> 💜</sub>
</div>
