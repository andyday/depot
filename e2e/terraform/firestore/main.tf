provider "google" {
  project = var.project
  region  = "us-west1"
}

locals {
  widget   = "depot-widget"
  message  = "depot-message"
  database = var.database
}

resource "google_firestore_database" "database" {
  project                 = var.project
  name                    = local.database
  location_id             = "us-west1"
  type                    = "FIRESTORE_NATIVE"
  delete_protection_state = "DELETE_PROTECTION_DISABLED"
  deletion_policy         = "DELETE"
}

resource "google_firestore_index" "created" {
  project    = var.project
  database   = google_firestore_database.database.name
  collection = local.widget

  fields {
    field_path = "tenantId"
    order      = "ASCENDING"
  }

  fields {
    field_path = "createdAt"
    order      = "DESCENDING"
  }
}

resource "google_firestore_index" "named" {
  project    = var.project
  database   = google_firestore_database.database.name
  collection = local.widget
  fields {
    field_path = "tenantId"
    order      = "ASCENDING"
  }
  fields {
    field_path = "expiration"
    order      = "ASCENDING"
  }
  fields {
    field_path = "name"
    order      = "ASCENDING"
  }
}

resource "google_firestore_index" "category" {
  project    = var.project
  database   = google_firestore_database.database.name
  collection = local.widget

  fields {
    field_path = "tenantId"
    order      = "ASCENDING"
  }

  fields {
    field_path = "category"
    order      = "ASCENDING"
  }
}

resource "google_firestore_index" "expired" {
  project    = var.project
  database   = google_firestore_database.database.name
  collection = local.widget

  fields {
    field_path = "expirationPartition"
    order      = "ASCENDING"
  }

  fields {
    field_path = "expiration"
    order      = "ASCENDING"
  }
}

resource "google_firestore_index" "message" {
  project    = var.project
  database   = google_firestore_database.database.name
  collection = local.message

  fields {
    field_path = "tenantId"
    order      = "ASCENDING"
  }

  fields {
    field_path = "id"
    order      = "DESCENDING"
  }
}
