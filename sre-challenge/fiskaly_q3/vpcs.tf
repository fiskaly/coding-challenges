# # Create VPC in specified region
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "6.4.0"

  name = var.vpc_name
  cidr = var.cidr_block_routable
  azs = local.aws_azs

  private_subnets = local.subnet_cidrs_private
  public_subnets  = local.subnet_cidrs_public

  tags                 = local.common_tags
  enable_dns_hostnames = true
  enable_dns_support   = true
  enable_nat_gateway   = true
  single_nat_gateway   = true

  public_subnet_tags = {
    "Tier" = "Public"
    "kubernetes.io/role/elb" = 1
  }

  private_subnet_tags = {
    "Tier" = "Private"
    "kubernetes.io/role/internal-elb" = 1
  }
}
