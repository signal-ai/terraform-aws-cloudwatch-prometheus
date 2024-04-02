resource "aws_s3_bucket" "cloudwatch_metrics_firehose_bucket" {
  bucket = var.aws_firehose_s3_bucket_name

  tags = var.tags
}

resource "aws_s3_bucket_acl" "cloudwatch_metrics_firehose_bucket_acl" {
  bucket = aws_s3_bucket.cloudwatch_metrics_firehose_bucket.id
  acl    = "private"
  depends_on = [
    aws_s3_bucket_ownership_controls.bucket_ownership_cloudwatch_firehose
  ]
}

resource "aws_s3_bucket_lifecycle_configuration" "cloudwatch_metrics_firehose_remove_after_10_days" {
  bucket = aws_s3_bucket.cloudwatch_metrics_firehose_bucket.id
  rule {
    status = "Enabled"
    id     = "remove_older_than_10_days"
    expiration {
      days = 10
    }
  }
}

resource "aws_s3_bucket_ownership_controls" "bucket_ownership_cloudwatch_firehose" {
  bucket = aws_s3_bucket.cloudwatch_metrics_firehose_bucket.id

  rule {
    object_ownership = "ObjectWriter"
  }
}

resource "aws_kinesis_firehose_delivery_stream" "cloudwatch_metrics_firehose_delivery_stream" {
  name        = var.aws_firehose_stream_name
  destination = "extended_s3"

  extended_s3_configuration {
    role_arn   = aws_iam_role.cloudwatch_metrics_firehose_role.arn
    bucket_arn = aws_s3_bucket.cloudwatch_metrics_firehose_bucket.arn

    processing_configuration {
      enabled = "true"

      processors {
        type = "Lambda"

        parameters {
          parameter_name  = "LambdaArn"
          parameter_value = "${aws_lambda_function.cloudwatch_metrics_firehose_prometheus_remote_write.arn}:$LATEST"
        }
        parameters {
          parameter_name  = "BufferSizeInMBs"
          parameter_value = "1"
        }
        parameters {
          parameter_name  = "BufferIntervalInSeconds"
          parameter_value = "300"
        }
      }
    }
  }

  tags = var.tags
}

resource "aws_iam_role" "cloudwatch_metrics_firehose_role" {
  name = "${var.aws_firehose_stream_name}-firehose-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "firehose.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    },
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF

  tags = var.tags
}

resource "aws_iam_role_policy" "cloudwatch_metrics_s3_policy" {
  name = "${var.aws_firehose_stream_name}-firehose-s3-policy"
  role = aws_iam_role.cloudwatch_metrics_firehose_role.id

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "s3:AbortMultipartUpload",
                "s3:GetBucketLocation",
                "s3:GetObject",
                "s3:ListBucket",
                "s3:ListBucketMultipartUploads",
                "s3:PutObject"
            ],
            "Resource": [
                "${aws_s3_bucket.cloudwatch_metrics_firehose_bucket.arn}",
                "${aws_s3_bucket.cloudwatch_metrics_firehose_bucket.arn}/*"
            ]
        }
    ]
}
EOF
}

resource "aws_iam_role_policy" "cloudwatch_metrics_firehose_lambda_policy" {
  name = "${var.aws_firehose_stream_name}-firehose-lambda-policy"
  role = aws_iam_role.cloudwatch_metrics_firehose_role.id

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "lambda:*"
            ],
            "Resource": [
                "${aws_lambda_function.cloudwatch_metrics_firehose_prometheus_remote_write.arn}",
                "${aws_lambda_function.cloudwatch_metrics_firehose_prometheus_remote_write.arn}:*"
            ]
        }
    ]
}
EOF
}
