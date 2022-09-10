---
icon: material/shield-lock-outline
---

# IAM Permissions

Limited access to Route53 and EC2 is required for `dns53` to work. Your IAM persona must have the following permissions granted:

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
    },
    {
      "Effect": "Allow",
      "Action": ["ec2:ModifyInstanceMetadataOptions"],
      "Resource": "arn:aws:ec2:<REGION>:<ACCOUNT>:instance/*" // (1)!
    }
  ]
}
```

1. Don't forget to replace the `<REGION>` and `<ACCOUNT>` placeholders with your specific AWS details, e.g. `arn:aws:ec2:eu-west-2:112233445566:instance/*`. You could also lock it down to a specific EC2 instance if you wanted :lock:

!!! warning "Aim for Least Privilege :lock:"

    It would be best if you fine-tuned this policy further to restrict access and adopt the mantra of "**least privilege**". You accept this policy at your own risk
