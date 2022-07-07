# IAM

For `dns53` to successfully manage a record set within a Route53 Private Hosted Zone, your IAM persona must have the following permissions granted:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["route53:GetHostedZone", "route53:ChangeResourceRecordSets"],
      "Resource": "arn:aws:route53:::hostedzone/*"
    },
    {
      "Effect": "Allow",
      "Action": ["route53:ListHostedZonesByVPC", "ec2:DescribeVpcs"],
      "Resource": "*"
    }
  ]
}
```
