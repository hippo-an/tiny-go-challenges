output "bucket_seoul_id" {
  value       = aws_s3_bucket.bucket_seoul.id
  description = "서울 버킷 ID"
}
output "bucket_tokyo_id" {
  value       = aws_s3_bucket.bucket_tokyo.id
  description = "도쿄 버킷 ID"
}
