variable "vpc_id" {
  type        = string
  description = "The VPC to create the lambda in (this should have network access to the prometheusremote write endpoints)"
}

variable "subnet_ids" {
  type        = list(string)
  description = "The subnet ids the create the lambda in (these should have network access to the prometheus remote write endpoints)"
}

variable "included_aws_namespaces" {
  type        = list(string)
  description = "The list of AWS Namespaces to include in the stream. All metrics from this namespace will be collected via the firehose"
}

variable "included_aws_namespace_metrics" {
  type        = map(list(string))
  description = "The list of AWS Namespaces and specific metrics from these Namespaces to include in the stream. Only the specified metrics will be collected via the firehose"
}

variable "aws_cloudwatch_metric_stream_name" {
  type        = string
  description = "The desired cloudwatch metric stream name that will be created"
}

variable "aws_firehose_stream_name" {
  type        = string
  description = "The desired firehose stream name that will be created and linked to the output of the cloudwatch metric stream"
}

variable "aws_firehose_s3_bucket_name" {
  type        = string
  description = "The s3 bucket name that will be the output of the firehose stream"
}

variable "aws_firehose_lambda_name" {
  type        = string
  description = "The lambda name that will attached to put events in the s3 bucket output of the firehose stream"
}

variable "prometheus_endpoints" {
  type        = list(string)
  description = "A list of prometheus remote write endpoints to write metrics"
}

variable "tags" {
  type        = map(string)
  description = "The standard tags to apply to every AWS resource."

  default = {}
}
