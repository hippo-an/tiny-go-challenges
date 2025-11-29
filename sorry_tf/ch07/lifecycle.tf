resource "aws_autoscaling_group" "this" {
  max_size         = 10
  min_size         = 1
  desired_capacity = 5

  lifecycle {
    ignore_changes = [
      desired_capacity,
    ]
  }
}
