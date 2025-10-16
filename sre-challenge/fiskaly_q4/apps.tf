module "ec2-instances" {
  source = "terraform-aws-modules/ec2-instance/aws"

  for_each = {for idx, instance_name in local.ec2_instance_matrix: idx => instance_name}

  name = each.value.instance_name

  ami                     = each.value.ami
  instance_type           = each.value.instance_type
  vpc_security_group_ids  = [each.value.security_group_id]
  subnet_id               = each.value.subnet_id
  availability_zone       = each.value.az
  monitoring              = false
  disable_api_termination = false
  enable_volume_tags      = false
  key_name                = aws_key_pair.david.key_name
  associate_public_ip_address = true

  root_block_device = {
    volume_type = "gp3"
    volume_size = each.value.root_volume_size
  }

  # Extra to avoid conflicts with alternative setup
  create_security_group = false

  tags = {
    "Project" = "Task4"
    "NodeGroup" = each.value.node_group_name
  }
}
