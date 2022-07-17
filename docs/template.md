# Go Templates

`dns53` supports Go [templating](https://pkg.go.dev/text/template) through a series of named fields. Using a templated field is as easy as writing `{{.IPv4}}`. Some fields are `formatted` to ensure they are URL compliant.

## EC2 Metadata Fields

| Field         | Description                                                                | Example                 |
| ------------- | -------------------------------------------------------------------------- | ----------------------- |
| `.IPv4`       | the private IPv4 address assigned to the launched EC2 instance `formatted` | `10-0-1-182`            |
| `.Region`     | the region where the EC2 instance was launched                             | `eu-west-2`             |
| `.VPC`        | the VPC ID of where the EC2 instance was launched                          | `vpc-016d173db537793d1` |
| `.AZ`         | the availability zone (AZ) of where the EC2 instance was launched          | `eu-west-2a`            |
| `.InstanceID` | the unique ID of the launched EC2 instance                                 | `i-03e092f544905abb2`   |
