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
  strestream_arn = aws_kinesis_stream.transaction.arn
  table_name     = aws_dynamodb_table.accounts.name
}

# kinesis stream to hold dynamodb events
resource "aws_kinesis_stream" "transactions" {
  name             = "transactions-stream-${var.stage}"
  shard_count      = 1
  retention_period = 24 # hours
}

# kinesis firehose stream to put events in s3
resource "aws_kinesis_firehose_delivery_stream" "transactions" {
  destination = "extended_s3"

  kinesis_source_configuration {
    kinesis_stream_arn = aws_kinesis_stream.transactions.arn
    role_arn           = aws_iam_role.firehose_delivery.arn
  }

  extended_s3_configuration {
    role_arn   = aws_iam_role.firehose_delivery.arn
    bucket_arn = aws_s3_bucket.transactions.arn

    processing_configuration {
      enabled = "true"
      processors {
        type = "Lambda"
        parameters {
          parameter_name  = "LambdaArn"
          parameter_value = ":$LATEST"
        }
      }
    }
  }
}

# s3 bucket to store transactions event archive
resource "aws_s3_bucket" "transactions" {
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
#   "kinesis:PutRecord",
#   "kinesis:PutRecords",
#   "kinesis:DescribeStream",
#   "kinesis:DescribeStreamSummary",
#   "kinesis:GetRecords",
#   "kinesis:GetShardIterator",
#   "kinesis:ListShards",
#   "kinesis:ListStreams",
#   "kinesis:SubscribeToShared",
#   "lambda:InvokeFunction",
#   "lambda:GetFunctionConfiguration"
# ],
