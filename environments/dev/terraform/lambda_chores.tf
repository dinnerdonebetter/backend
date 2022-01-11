resource "aws_sqs_queue" "chores_dead_letter" {
  name                    = "chores_dead_letter"
  sqs_managed_sse_enabled = true
}

resource "aws_sqs_queue" "chores_queue" {
  name                    = "chores"
  sqs_managed_sse_enabled = true

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.chores_dead_letter.arn
    maxReceiveCount     = 1
  })
}

resource "aws_lambda_function" "chores_worker_lambda" {
  function_name = "chores_worker"
  handler       = "chores_worker"
  role          = aws_iam_role.worker_lambda_role.arn
  runtime       = local.lambda_runtime
  memory_size   = local.memory_size
  timeout       = 55

  tracing_config {
    mode = "Active"
  }

  vpc_config {
    subnet_ids = concat(
      [for x in aws_subnet.public_subnets : x.id],
      [for x in aws_subnet.private_subnets : x.id],
    )
    security_group_ids = [
      aws_security_group.lambda_workers.id,
    ]
  }

  layers = [
    local.collector_layer_arns.us-east-1,
  ]

  filename = data.archive_file.dummy_zip.output_path

  depends_on = [
    aws_cloudwatch_log_group.chores_worker_lambda_logs,
  ]
}

resource "aws_cloudwatch_event_rule" "every_minute" {
  name                = "every-minute"
  description         = "Fires every minute"
  schedule_expression = "rate(1 minute)"
}

resource "aws_cloudwatch_event_target" "run_chores_every_minute" {
  rule      = aws_cloudwatch_event_rule.every_minute.name
  target_id = "chores_worker"
  arn       = aws_lambda_function.chores_worker_lambda.arn
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_chores_worker" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.chores_worker_lambda.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.every_minute.arn
}

resource "aws_cloudwatch_log_group" "chores_worker_lambda_logs" {
  name              = "/aws/lambda/chores_worker"
  retention_in_days = local.log_retention_period_in_days
}