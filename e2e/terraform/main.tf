terraform {
  backend "s3" {
    bucket = "depot-terraform"
    key    = "state"
    region = "us-west-2"
  }
}

module "dynamo" {
  source = "./dynamo"
}

# module "datastore" {
#   source  = "./datastore"
#   env     = var.env
#   project = var.project
#   database = "(default)"
# }

module "firestore" {
  source = "./firestore"
  project = var.project
  database = "depot-e2e"
}