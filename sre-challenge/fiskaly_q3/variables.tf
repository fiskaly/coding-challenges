# Assuming single value only as input due to data type, not looping
variable "aws_region" {
  type    = string
  default = "eu-central-1"
}

variable "vpc_name" {
  type    = string
  default = "fiskaly-task3-vpc"
}

variable "eks_name" {
  type    = string
  default = "fiskaly-task3-eks"
}

### Extra
variable "cidr_block_routable" {
  type    = string
  default = "10.1.0.0/16"
}
