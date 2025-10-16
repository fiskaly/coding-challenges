locals {
  # Return AZs per AWS region based on AZ limit (default 3)
  aws_azs = data.aws_availability_zones.available.names
  # public_subnet_ids_az = {for k, v in data.aws_subnets.public_per_az : tolist(v.filter)[0].values[0] => v.ids if length(v.ids) > 0}

  # Create map with matrix of applications and AZs
  ec2_instance_matrix = flatten([
    for ec2_node_group, ec2_node_group_values in var.ec2_node_groups : [
      for az in var.aws_azs : {
        instance_name  = "${ec2_node_group_values.name}-${az}"
        node_group_name = ec2_node_group_values.name
        ami               = ec2_node_group_values.ami
        instance_type     = ec2_node_group_values.instance_type
        root_volume_size  = ec2_node_group_values.root_volume_size
        az                = az
        subnet_id         = data.aws_subnets.public_per_az[az].ids[0]
        security_group_id = aws_security_group.fiskaly-q4.id
      }
    ]
  ])

  common_tags = {
    "aws.region"      = var.aws_region
    "fiskaly.project" = "hiring-excercise"
  }
}
