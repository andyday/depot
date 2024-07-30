# module "dynamo" {
#   source = "./dynamo"
#   env    = var.env
# }

# module "datastore" {
#   source  = "./datastore"
#   env     = var.env
#   project = var.project
#   database = "(default)"
# }

module "firestore" {
  source = "./firestore"
  env = var.env
  project = var.project
  database = "(default)"
}