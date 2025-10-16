# Task 3
## Comments
- using "~/workspaces/fiskaly_q3" as workdir
- working in my personal AWS account
- deploying EKS cluster in 3 AZs with 4 nodes
- based on examples from Terraform Registry

## Steps
1. Run Terraform to deploy VPC/EKS/ECR
    ```
    tf apply
    ```

1. Update kubeconfig with EKS credentials
    ```
    aws eks update-kubeconfig --region eu-central-1 --name fiskaly-task3-eks
    ```

1. Deploy metrics-server needed for HPA to work
    ```
    kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
    ```
