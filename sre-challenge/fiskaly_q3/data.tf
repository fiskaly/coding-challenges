data "aws_availability_zones" "available" {
  region = var.aws_region
}

data "aws_subnets" "private" {
  tags = {
    Tier = "Private"
  }
}

data "aws_subnets" "public" {
  tags = {
    Tier = "Public"
  }
}