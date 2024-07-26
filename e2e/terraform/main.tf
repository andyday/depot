module "aws" {
  source = "./aws"
  env    = var.env
}

module "gcp" {
  source  = "./gcp"
  env     = var.env
  project = var.project
}