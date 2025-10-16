locals {
  # Return AZs per AWS region based on AZ limit (default 3)
  aws_azs = data.aws_availability_zones.available.names

  # Return private subnets for AZs
  subnet_cidrs_private = [for netnumber in range(0, length(local.aws_azs)) : cidrsubnet(var.cidr_block_routable, 4, netnumber)]
  subnet_cidrs_public  = [for netnumber in range(0, length(local.aws_azs)) : cidrsubnet(var.cidr_block_routable, 4, netnumber + 6)]

  common_tags = {
    "aws.region"      = var.aws_region
    "fiskaly.project" = "hiring-excercise"
  }
}
