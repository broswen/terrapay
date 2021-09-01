terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.56.0"
    }

  }
  backend "s3" {
    bucket = "broswen-terraform"
    key    = "terraform_test/terraform.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  profile = "default"
  region  = var.region
}

# dynamodb table: single table design to hold accounts and transactions
resource "aws_dynamodb_table" "accounts" {
  name         = "${var.name}-${var.stage}-accounts"
  billing_mode = "PAY_PER_REQUEST"
  attribute {
    name = "PK"
    type = "S"
  }
  attribute {
    name = "SK"
    type = "S"
  }
  hash_key         = "PK"
  range_key        = "SK"
  stream_enabled   = true
  stream_view_type = "NEW_IMAGE"
}

# maps dynamodb stream events to kinesis data stream
resource "aws_dynamodb_kinesis_streaming_destination" "transactions" {
  stream_arn = aws_kinesis_stream.transactions.arn
  table_name = aws_dynamodb_table.accounts.name
}

# kinesis stream to hold dynamodb events
resource "aws_kinesis_stream" "transactions" {
  name             = "${var.name}-${var.stage}-transactions"
  shard_count      = 1
  retention_period = 24 # hours
}

# kinesis firehose stream to put events in s3
resource "aws_kinesis_firehose_delivery_stream" "transactions" {
  name        = "${var.name}-${var.stage}-transactions"
  destination = "extended_s3"

  kinesis_source_configuration {
    kinesis_stream_arn = aws_kinesis_stream.transactions.arn
    role_arn           = aws_iam_role.delivery_role.arn
  }

  extended_s3_configuration {
    role_arn   = aws_iam_role.delivery_role.arn
    bucket_arn = aws_s3_bucket.transactions.arn

    processing_configuration {
      enabled = "true"
      processors {
        type = "Lambda"
        parameters {
          parameter_name  = "LambdaArn"
          parameter_value = "${module.event_processor.lambda_function_arn}:$LATEST"
        }
      }
    }
  }
}

# s3 bucket to store transactions event archive
resource "aws_s3_bucket" "transactions" {
}

module "register_lambda" {
  source        = "terraform-aws-modules/lambda/aws"
  function_name = "${var.name}-${var.stage}-register"
  handler       = "register"
  runtime       = "go1.x"
  memory_size   = 128
  timeout       = 3

  source_path = "./bin/register"

  environment_variables = {
    ACCOUNTSDB = aws_dynamodb_table.accounts.name
  }

  attach_policies = true
  policies = [
    aws_iam_policy.read_write_accounts.arn
  ]
}

module "login_lambda" {
  source        = "terraform-aws-modules/lambda/aws"
  function_name = "${var.name}-${var.stage}-login"
  handler       = "login"
  runtime       = "go1.x"
  memory_size   = 128
  timeout       = 3

  source_path = "./bin/login"

  environment_variables = {
    ACCOUNTSDB = aws_dynamodb_table.accounts.name
  }

  attach_policies = true
  policies = [
    aws_iam_policy.read_write_accounts.arn
  ]
}

module "post_transaction" {
  source        = "terraform-aws-modules/lambda/aws"
  function_name = "${var.name}-${var.stage}-post-transaction"
  handler       = "post_transaction"
  runtime       = "go1.x"
  memory_size   = 128
  timeout       = 3

  source_path = "./bin/post_transaction"

  environment_variables = {
    ACCOUNTSDB = aws_dynamodb_table.accounts.name
  }

  attach_policies = true
  policies = [
    aws_iam_policy.read_write_accounts.arn
  ]
}

module "get_transactions" {
  source        = "terraform-aws-modules/lambda/aws"
  function_name = "${var.name}-${var.stage}-get-transactions"
  handler       = "get_transactions"
  runtime       = "go1.x"
  memory_size   = 128
  timeout       = 3

  source_path = "./bin/get_transactions"

  environment_variables = {
    ACCOUNTSDB = aws_dynamodb_table.accounts.name
  }

  attach_policies = true
  policies = [
    aws_iam_policy.read_write_accounts.arn
  ]
}

module "send_notification" {
  source        = "terraform-aws-modules/lambda/aws"
  function_name = "${var.name}-${var.stage}-send-notification"
  handler       = "send_notification"
  runtime       = "go1.x"
  memory_size   = 128
  timeout       = 3

  source_path = "./bin/send_notification"

  environment_variables = {
    ACCOUNTSDB = aws_dynamodb_table.accounts.name
  }
}

module "event_processor" {
  source        = "terraform-aws-modules/lambda/aws"
  function_name = "${var.name}-${var.stage}-event-processor"
  handler       = "event_processor"
  runtime       = "go1.x"
  memory_size   = 128
  timeout       = 3

  source_path = "./bin/event_processor"
}

resource "aws_iam_policy" "read_write_accounts" {
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:Query",
          "dynamodb:TransactWriteItems"
        ]
        Effect = "Allow"
        Resource = [
          "${aws_dynamodb_table.accounts.arn}"
        ]
      },
    ]
  })
}

resource "aws_iam_policy" "kinesis_read" {
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "kinesis:PutRecord",
          "kinesis:PutRecords",
          "kinesis:DescribeStream",
          "kinesis:DescribeStreamSummary",
          "kinesis:GetRecords",
          "kinesis:GetShardIterator",
          "kinesis:ListShards",
          "kinesis:ListStreams",
          "kinesis:SubscribeToShard"
        ]
        Effect = "Allow"
        Resource = [
          "${aws_kinesis_firehose_delivery_stream.transactions.arn}",
        ]
      },
    ]
  })
}

resource "aws_iam_policy" "s3_put" {
  policy = jsonencode({
    Version = "2021-10-17"
    Statement = [
      {
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject",
          "s3:CopyObject",
        ]
        Effect = "Allow"
        Resource = [
          "${aws_s3_bucket.transactions.arn}",
        ]
      }
    ]
  })
}

resource "aws_iam_role" "delivery_role" {
  assume_role_policy = jsonencode({
    Version = "2021-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "firehose.amazonaws.com"
        }
      }
    ]
  })
}

# iam role for firehose to read from kinesis and write to s3
# iam role for lambda event processor to run
# iam role for api lambdas to read/write to dynamodb
# iam role for transaction notification to publish to sns and read from dynamodb stream

# iam role 
# "Action": [
#   "s3:GetObject",
#   "s3:PutObject",
#   "s3:DeleteObject",
#   "s3:CopyObject",

#   "lambda:InvokeFunction",
#   "lambda:GetFunctionConfiguration"
# ],

