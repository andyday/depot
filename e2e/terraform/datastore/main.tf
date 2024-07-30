provider "google" {
  project     = var.project
  region      = "us-west1"
}

locals {
  widget = "depot-${var.env}-widget"
  message = "depot-${var.env}-message"
  database = "(default)"
}

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
  properties {
    name      = "expiration"
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

resource "google_datastore_index" "message" {
  kind = local.message
  properties {
    name      = "tenantId"
    direction = "ASCENDING"
  }
  properties {
    name      = "id"
    direction = "DESCENDING"
  }
}
