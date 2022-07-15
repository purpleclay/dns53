# Go Templates

`dns53` supports Go [templating](https://pkg.go.dev/text/template) through a series of named fields. Using a templated field is as easy as writing `{{.IPv4}}`.

## EC2 Metadata Fields

| Field     | Description                                                    |
| --------- | -------------------------------------------------------------- |
| `.IPv4`   | the private IPv4 address assigned to the launched EC2 instance |
| `.Region` | the region where the EC2 instance was launched                 |
| `.VPC`    | the VPC ID of where the EC2 instance was launched              |
