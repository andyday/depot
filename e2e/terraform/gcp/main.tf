provider "google" {
  project     = var.project
  region      = "us-west1"
}

locals {
  widget = "depot-${var.env}-widget"
  database = "(default)"
}

# resource "google_firestore_database" "database" {
#   project                 = var.project
#   name                    = "(default)"
#   location_id             = "us-west1"
#   type                    = "DATASTORE_MODE"
#   delete_protection_state = "DELETE_PROTECTION_DISABLED"
#   deletion_policy         = "DELETE"
# }

# resource "google_firestore_index" "created" {
#   project    = var.project
#   database   = local.database
#   collection = local.widget
#
#   fields {
#     field_path = "tenantId"
#     order      = "ASCENDING"
#   }
#
#   fields {
#     field_path = "createdAt"
#     order      = "ASCENDING"
#   }
# }
#
# resource "google_firestore_index" "named" {
#   project    = var.project
#   database   = local.database
#   collection = local.widget
#
#   fields {
#     field_path = "tenantId"
#     order      = "ASCENDING"
#   }
#
#   fields {
#     field_path = "name"
#     order      = "ASCENDING"
#   }
# }
#
# resource "google_firestore_index" "category" {
#   project    = var.project
#   database   = local.database
#   collection = local.widget
#
#   fields {
#     field_path = "tenantId"
#     order      = "ASCENDING"
#   }
#
#   fields {
#     field_path = "category"
#     order      = "ASCENDING"
#   }
# }
#
# resource "google_firestore_index" "expired" {
#   project    = var.project
#   database   = local.database
#   collection = local.widget
#
#   fields {
#     field_path = "expirationPartition"
#     order      = "ASCENDING"
#   }
#
#   fields {
#     field_path = "expiration"
#     order      = "ASCENDING"
#   }
# }



resource "google_datastore_index" "created" {
  kind = local.widget
  properties {
    name      = "tenantId"
    direction = "ASCENDING"
  }
  properties {
    name      = "createdAt"
    direction = "DESCENDING"
  }
}

resource "google_datastore_index" "named" {
  kind = local.widget
  properties {
    name      = "tenantId"
    direction = "ASCENDING"
  }
  properties {
    name      = "name"
    direction = "ASCENDING"
  }
}

resource "google_datastore_index" "category" {
  kind = local.widget
  properties {
    name      = "tenantId"
    direction = "ASCENDING"
  }
  properties {
    name      = "category"
    direction = "ASCENDING"
  }
}

resource "google_datastore_index" "expired" {
  kind = local.widget
  properties {
    name      = "expirationPartition"
    direction = "ASCENDING"
  }
  properties {
    name      = "expiration"
    direction = "ASCENDING"
  }
}