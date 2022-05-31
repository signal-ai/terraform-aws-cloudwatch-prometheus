# terraform-aws-cloudwatch-prometheus

![System Architecture](./images/system_architecture.png)

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 0.13 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 3.75.1 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | >= 3.75.1 |

## Resources

| Name | Type |
|------|------|
| [aws_cloudwatch_log_group.logs](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_log_group) | resource |
| [aws_cloudwatch_metric_stream.main](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_metric_stream) | resource |
| [aws_iam_role.cloudwatch_metrics_firehose_role](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role) | resource |
| [aws_iam_role.iam_for_lambda](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role) | resource |
| [aws_iam_role.metric_stream_to_firehose](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role) | resource |
| [aws_iam_role_policy.cloudwatch_metrics_firehose_lambda_policy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy) | resource |
| [aws_iam_role_policy.cloudwatch_metrics_s3_policy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy) | resource |
| [aws_iam_role_policy.metric_stream_to_firehose](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy) | resource |
| [aws_iam_role_policy_attachment.execution](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy_attachment) | resource |
| [aws_iam_role_policy_attachment.vpc](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy_attachment) | resource |
| [aws_kinesis_firehose_delivery_stream.cloudwatch_metrics_firehose_delivery_stream](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/kinesis_firehose_delivery_stream) | resource |
| [aws_lambda_function.cloudwatch_metrics_firehose_prometheus_remote_write](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function) | resource |
| [aws_s3_bucket.cloudwatch_metrics_firehose_bucket](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket) | resource |
| [aws_s3_bucket_acl.cloudwatch_metrics_firehose_bucket_acl](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_acl) | resource |
| [aws_security_group.cloudwatch_metrics_firehose_prometheus_remote_write](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group) | resource |
| [aws_security_group_rule.cloudwatch_metrics_firehose_prometheus_remote_write](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group_rule) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_aws_cloudwatch_metric_stream_name"></a> [aws\_cloudwatch\_metric\_stream\_name](#input\_aws\_cloudwatch\_metric\_stream\_name) | The desired cloudwatch metric stream name that will be created | `string` | n/a | yes |
| <a name="input_aws_firehose_lambda_name"></a> [aws\_firehose\_lambda\_name](#input\_aws\_firehose\_lambda\_name) | The lambda name that will attached to put events in the s3 bucket output of the firehose stream | `string` | n/a | yes |
| <a name="input_aws_firehose_s3_bucket_name"></a> [aws\_firehose\_s3\_bucket\_name](#input\_aws\_firehose\_s3\_bucket\_name) | The s3 bucket name that will be the output of the firehose stream | `string` | n/a | yes |
| <a name="input_aws_firehose_stream_name"></a> [aws\_firehose\_stream\_name](#input\_aws\_firehose\_stream\_name) | The desired firehose stream name that will be created and linked to the output of the cloudwatch metric stream | `string` | n/a | yes |
| <a name="input_prometheus_endpoints"></a> [prometheus\_endpoints](#input\_prometheus\_endpoints) | A list of prometheus remote write endpoints to write metrics | `list(string)` | n/a | yes |
| <a name="input_subnet_ids"></a> [subnet\_ids](#input\_subnet\_ids) | The subnet ids the create the lambda in (these should have network access to the prometheus remote write endpoints) | `list(string)` | n/a | yes |
| <a name="input_vpc_id"></a> [vpc\_id](#input\_vpc\_id) | The VPC to create the lambda in (this should have network access to the prometheusremote write endpoints) | `string` | n/a | yes |
<!-- END_TF_DOCS -->
