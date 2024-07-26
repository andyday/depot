terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# Configure the AWS Provider
provider "aws" {
  region = "us-west-2"
}

resource "aws_dynamodb_table" "widget" {
  name         = "depot-${var.env}-widget"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "tenantId"
  range_key    = "id"

  attribute {
    name = "tenantId"
    type = "S"
  }
  attribute {
    name = "id"
    type = "S"
  }
  attribute {
    name = "name"
    type = "S"
  }
  attribute {
    name = "category"
    type = "S"
  }
  attribute {
    name = "createdAt"
    type = "S"
  }
  attribute {
    name = "expirationPartition"
    type = "N"
  }
  attribute {
    name = "expiration"
    type = "S"
  }

  global_secondary_index {
    name            = "created"
    hash_key        = "tenantId"
    range_key       = "createdAt"
    projection_type = "ALL"
  }
  global_secondary_index {
    name            = "named"
    hash_key        = "tenantId"
    range_key       = "name"
    projection_type = "ALL"
  }
  global_secondary_index {
    name            = "category"
    hash_key        = "tenantId"
    range_key       = "category"
    projection_type = "ALL"
  }
  global_secondary_index {
    name            = "expired"
    hash_key        = "expirationPartition"
    range_key       = "expiration"
    projection_type = "ALL"
  }
}