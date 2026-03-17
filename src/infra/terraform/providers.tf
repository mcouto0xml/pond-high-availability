terraform {
  required_providers {
    supabase = {
      source  = "supabase/supabase"
      version = "~> 1.0"
    }

    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

provider "supabase" {
  access_token = var.supabase_access_token
}

provider "google" {
  project = var.project_id
  region  = var.region
}