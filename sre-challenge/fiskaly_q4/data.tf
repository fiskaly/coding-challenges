data "aws_availability_zones" "available" {
  region = var.aws_region
}

data "aws_vpc" "main" {
  region = var.aws_region
}

data "aws_subnets" "public_per_az" {
  for_each = toset(data.aws_availability_zones.available.names)

  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.main.id]
  }

  filter {
    name   = "tag:Tier"
    values = ["Public"]
  }

  filter {
    name   = "availability-zone"
    values = ["${each.value}"]
  }
}
