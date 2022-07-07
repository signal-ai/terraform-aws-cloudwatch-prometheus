
resource "aws_cloudwatch_metric_stream" "main" {
  name          = var.aws_cloudwatch_metric_stream_name
  role_arn      = aws_iam_role.metric_stream_to_firehose.arn
  firehose_arn  = aws_kinesis_firehose_delivery_stream.cloudwatch_metrics_firehose_delivery_stream.arn
  output_format = "json"

  include_filter {
    namespace = "AWS/SQS"
  }

  include_filter {
    namespace = "AWS/Lambda"
  }

  include_filter {
    namespace = "AWS/ElastiCache"
  }

  include_filter {
    namespace = "AWS/Redshift"
  }

  include_filter {
    namespace = "AWS/ES"
  }

  include_filter {
    namespace = "AWS/ES"
  }
  include_filter {
    namespace = "AWS/Redshift"
  }

}

resource "aws_iam_role" "metric_stream_to_firehose" {
  name = "${var.aws_cloudwatch_metric_stream_name}-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "streams.metrics.cloudwatch.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "metric_stream_to_firehose" {
  name = "${var.aws_cloudwatch_metric_stream_name}-firehose-policy"
  role = aws_iam_role.metric_stream_to_firehose.id

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "firehose:PutRecord",
                "firehose:PutRecordBatch"
            ],
            "Resource": "${aws_kinesis_firehose_delivery_stream.cloudwatch_metrics_firehose_delivery_stream.arn}"
        }
    ]
}
EOF
}
