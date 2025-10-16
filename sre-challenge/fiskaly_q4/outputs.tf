output "instance_matrix" {
  value = toset(local.ec2_instance_matrix)
}