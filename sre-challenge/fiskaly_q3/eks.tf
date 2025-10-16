module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "21.3.2"

  name               = var.eks_name
  kubernetes_version = "1.33"

  addons = {
    coredns                = {}
    eks-pod-identity-agent = {
      before_compute = true
    }
    kube-proxy             = {}
    vpc-cni                = {
      before_compute = true
    }
  }

  endpoint_public_access           = true
  enable_cluster_creator_admin_permissions = true

  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets

  eks_managed_node_groups = {
    az1 = {
      name     = "node-group-az1"
      ami_type = "AL2023_x86_64_STANDARD"
      subnet_ids = [module.vpc.private_subnets[0]]

      instance_types = ["t3.small"]

      min_size     = 1
      max_size     = 1
      desired_size = 1
    }

    az2 = {
      name     = "node-group-az2"
      ami_type = "AL2023_x86_64_STANDARD"
      subnet_ids = [module.vpc.private_subnets[1]]

      instance_types = ["t3.small"]

      min_size     = 1
      max_size     = 1
      desired_size = 1
    }

    az3 = {
      name     = "node-group-az3"
      ami_type = "AL2023_x86_64_STANDARD"
      subnet_ids = [module.vpc.private_subnets[2]]

      instance_types = ["t3.small"]

      min_size     = 2
      max_size     = 2
      desired_size = 2
    }

  }
}