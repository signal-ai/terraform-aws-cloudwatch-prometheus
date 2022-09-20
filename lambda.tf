resource "aws_lambda_function" "cloudwatch_metrics_firehose_prometheus_remote_write" {
  filename         = "${path.module}/lambda_code/payload.zip"
  source_code_hash = filebase64sha256("${path.module}/lambda_code/payload.zip")
  function_name    = var.aws_firehose_lambda_name
  role             = aws_iam_role.iam_for_lambda.arn
  handler          = "main"
  timeout          = 60
  memory_size      = 256

  runtime = "go1.x"

  vpc_config {
    subnet_ids = var.subnet_ids

    security_group_ids = [aws_security_group.cloudwatch_metrics_firehose_prometheus_remote_write.id]
  }

  environment {
    variables = {
      PROMETHEUS_REMOTE_WRITE_URLS = join(",", var.prometheus_endpoints)
    }
  }

  tags = var.tags
}

resource "aws_security_group" "cloudwatch_metrics_firehose_prometheus_remote_write" {
  name   = "${var.aws_firehose_lambda_name}-security-group"
  vpc_id = var.vpc_id

  tags = var.tags
}

resource "aws_security_group_rule" "cloudwatch_metrics_firehose_prometheus_remote_write" {
  security_group_id = aws_security_group.cloudwatch_metrics_firehose_prometheus_remote_write.id
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "all"
  cidr_blocks       = ["0.0.0.0/0"]
}


resource "aws_iam_role" "iam_for_lambda" {
  name = "${var.aws_firehose_lambda_name}-lambda-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
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

data "aws_iam_policy" "lambda_basic_execution_role_policy_vpc" {
  arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"

  tags = var.tags
}

data "aws_iam_policy" "lambda_basic_execution_role_policy" {
  arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"

  tags = var.tags
}

resource "aws_iam_role_policy_attachment" "vpc" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = data.aws_iam_policy.lambda_basic_execution_role_policy_vpc.arn
}

resource "aws_iam_role_policy_attachment" "execution" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = data.aws_iam_policy.lambda_basic_execution_role_policy.arn
}


resource "aws_cloudwatch_log_group" "logs" {
  name              = "/aws/lambda/${aws_lambda_function.cloudwatch_metrics_firehose_prometheus_remote_write.function_name}"
  retention_in_days = 30

  tags = var.tags
}
