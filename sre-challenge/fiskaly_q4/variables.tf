variable "aws_region" {
  type    = string
  default = "eu-central-1"
}

variable "aws_azs" {
  type    = list(string)
  default     = [
    "eu-central-1a",
    "eu-central-1b",
    "eu-central-1c"
  ]
}

variable "ec2_node_groups" {
  type = map(object({
    name             = string
    root_volume_size = number
    ami              = string
    instance_type    = string
  }))
  default = {
    "ubuntu" = {
      name             = "ubuntu"
      root_volume_size = 50
      ami              = "ami-0a116fa7c861dd5f9"
      instance_type    = "t2.micro"
    },
    "redhat" = {
      name             = "redhat"
      root_volume_size = 50
      ami              = "ami-005c89b47f40aa0e1"
      instance_type    = "t2.micro"
    }
  }
}
