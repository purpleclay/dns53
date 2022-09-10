---
icon: material/application-cog-outline
---

# Go Templates

Full support for Go [templates](https://pkg.go.dev/text/template) through a series of predefined named fields allows `dns53`` to support a dynamic configuration where needed.

!!! info "Table Key"

    This is a living table and will change as new features are released.

    - :material-pencil-plus-outline:: the metadata was formatted to ensure it is URL compliant
    - :material-tag-outline:: the metadata was retrieved from EC2 instance tags; this feature must be [enabled](../configure/exposing-tags.md)

## Named Fields

The following named fields directly access metadata about your EC2 from the Instance Metadata Service (IMDS).

| Named Field                                                                  | Description                                       | Example                                                                        |
| ---------------------------------------------------------------------------- | ------------------------------------------------- | ------------------------------------------------------------------------------ |
| `{{.IPv4}}`                                                                  | the private IPv4 address of the EC2 instance      | `10-0-1-182` :material-pencil-plus-outline:{title="formatted from 10.0.1.182"} |
| `{{.Region}}`                                                                | the region of the EC2 instance                    | `eu-west-2`                                                                    |
| `{{.VPC}}`                                                                   | the VPC ID of where the EC2 instance was launched | `vpc-016d173db537793d1`                                                        |
| `{{.AZ}}`                                                                    | the availability zone (AZ) of the EC2 instance    | `eu-west-2a`                                                                   |
| `{{.InstanceID}}`                                                            | the unique ID of the EC2 instance                 | `i-03e092f544905abb2`                                                          |
| `{{.Name}}` :material-tag-outline:{title="retrieved from EC2 instance tags"} | a name assigned to the EC2 instance               | `dev-ec2` :material-pencil-plus-outline:                                       |
