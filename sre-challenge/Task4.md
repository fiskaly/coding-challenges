# Task 4
## Comments
- using "~/workspaces/fiskaly_q4" as workdir
- working on VPC deployed in my personal AWS account as part of this excercise (Task 3)

## Steps
1. Deploy 3 Ubuntu and 3 RedHat EC2 instances
    - deployed in VPC public subnets
    - deployed using SG allowing access to 22 and 80 from internet
    - deployed with correct personal SSH key
    - deployed with tag "Project:Task4"
    ```
    tf apply
    ```
1. Prepare simple ansible config and enable "aws_ec2" plugin for registry (~/.ansible.cfg)
    ```
    [defaults]
    host_key_checking = False
    inventory =  ~/workspaces/fiskaly_q4/ansible/inventory/aws_ec2.yaml
    interpreter_python = auto_silent
    deprecation_warnings = False
    [inventory]
    enable_plugins = aws_ec2
    ```

1. Prepare "aws_ec2" plugin config to filter out EKS instances (~/workspaces/fiskaly_q4/ansible/inventory/aws_ec2.yaml)
    ```
    plugin: aws_ec2

    filters:
    tag:Project:
        - "Task4"

    keyed_groups:
    - key: "tags.NodeGroup"
        prefix: node_group
    ```

1. Prepare Ansible group_vars to set correct SSH user for Ubuntu/Redhat based on keyed_group:
    - ~/workspaces/fiskaly_q4/ansible/inventory/group_vars/node_group_redhat
        ```
        ansible_user: ec2-user
        ```
    - ~/workspaces/fiskaly_q4/ansible/inventory/group_vars/node_group_ubuntu
        ```
        ansible_user: ubuntu
        ```

1. Run playbook that:
    - gathers facts about all hosts
    - updates repos and packages on both Ubuntu and RHEL instances
    - installs Apache on Ubuntu instances and make sure Apache is not present on RHEL instances
    - if Ubuntu OS + Apache installed -> configure index.html and restart Apache
    - installs MariaDB on RHEL and uninstalls it on Debian