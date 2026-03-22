output "storage_name" {
  value = google_storage_bucket.function_source.name
}

output "object_name" {
  value = google_storage_bucket_object.function_source.name
}