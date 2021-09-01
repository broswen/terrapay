

resource "aws_sqs_queue" "queue" {
  message_retention_seconds = var.message_retention_seconds
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.dlq.arn
    maxReceiveCount     = var.max_receive_count
  })
}

resource "aws_sqs_queue" "dlq" {
  message_retention_seconds = var.message_retention_seconds
}
