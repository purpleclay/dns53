# Quick Start

It's really easy to get up and running with `dns53`. You can expose your EC2 in a matter of seconds. ⚡

## Full Wizard

If you don't know which Amazon Private Hosted Zone to use, `dns53` provides a handy wizard.

```sh
dns53
```

## I have a PHZ ID

Skip the wizard and expose your EC2 straight away.

```sh
dns53 --phz-id AAAAAAAAAAAAAAAAA
```

## Custom Domain Name

If you want more control over the generated domain name for your exposed EC2, you have two options.

=== "As is"

    ```sh
    dns53 --dns-name my.ec2
    ```

=== "Templated"

    ```sh
    dns53 --dns-name "{{.IPv4}}.{{.Region}}"
    ```

!!!tip "Why not try templating?"

    A full list of fields supported by `dns53` templating can be found [here](./template.md)
