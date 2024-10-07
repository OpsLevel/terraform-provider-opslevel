resource "random_string" "external_id" {
  length  = 16
  special = false
}

resource "aws_iam_policy" "opslevel" {
  name = "opslevel-demo"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "autoscaling:Describe*",
          "cloudfront:GetDistribution",
          "cloudfront:ListDistributions",
          "dynamodb:Describe*",
          "dynamodb:List*",
          "ec2:Describe*",
          "ec2:Get*",
          "ecs:Describe*",
          "ecs:List*",
          "elasticloadbalancing:Describe*",
          "elasticache:Describe*",
          "elasticache:List*",
          "es:Describe*",
          "es:List*",
          "es:Get*",
          "grafana:Describe*",
          "grafana:List*",
          "kafka:Describe*",
          "kafka:List*",
          "kinesis:Get*",
          "kinesis:List*",
          "kinesis:Describe*",
          "lambda:Get*",
          "lambda:List*",
          "rds:Describe*",
          "rds:List*",
          "redshift:Describe*",
          "redshift:List*",
          "route53domains:Get*",
          "route53domains:List*",
          "s3:Describe*",
          "s3:List*",
          "s3:GetBucketLocation",
          "s3:GetBucketTagging",
          "s3:GetBucketPolicyStatus",
          "s3:GetBucketVersioning",
          "sns:Get*",
          "sns:List*",
          "storagegateway:List*",
          "storagegateway:Describe*",
          "sqs:Get*",
          "sqs:List*",
          "tag:Get*",
          "waf:Get*",
          "waf:List*",
          "wafv2:Get*",
          "wafv2:List*",
          "wafv2:Describe*"
        ]
        Effect   = "Allow"
        Resource = "*"
      },
    ]
  })
}

resource "aws_iam_role" "opslevel" {
  name                = "opslevel-demo"
  managed_policy_arns = [aws_iam_policy.opslevel.arn]
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "ec2.amazonaws.com",
          AWS     = "arn:aws:iam::746108190720:user/opslevel-integration",
        }
        Condition = {
          StringEquals = {
            "sts:ExternalId" = random_string.external_id.result
          }
        }
      },
    ]
  })
}

resource "opslevel_integration_aws" "dev" {
  name                    = "dev"
  iam_role                = aws_iam_role.opslevel.arn
  external_id             = random_string.external_id.result
  ownership_tag_overrides = true
  ownership_tag_keys      = ["owner", "team", "group"]
  region_override         = ["eu-west-1", "us-east-1"]
}
