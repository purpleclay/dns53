# Named Templates

`dns53` supports templating through a series of named templates

## EC2 Metadata Fields

| Field     | Description                                                    |
| --------- | -------------------------------------------------------------- |
| `.IPv4`   | the private IPv4 address assigned to the launched EC2 instance |
| `.Region` | the region where the EC2 instance was launched                 |
| `.VPC`    | the VPC ID of where the EC2 instance was launched              |
