# SRE Challenge

## Instructions

For all exercises documentation should be provided - this can be in the form of a README for each of them - describing details, like reasoning about choices made for configurations etc.

### The Challenge

#### 1) Docker Container

Present a simple Dockerfile for building a small web application container image, with application responding with “Hello world” to HTTP requests on port 8080.

- Own code (Typescript / Python) - AI tools may be used to generate this.
- Otherwise, use a ready container of choice, doing the same - like nginx, apache, lighttpd as a web server with a simple html document added.
- Provide also a docker command to launch the container.
- The launched web application/server should be accessible ideally from local computer as well as other computers in the same network (we can disregard here the potential firewall configuration issues).

#### 2) Kubernetes Deployment

Present an example of manifests to arrange a simple Kubernetes deployment with the following:

- The simple web application built in previous task, configured to launch at least 2 replicas, and say maximum 4 in case of higher load.
- In front of the application use nginx as a load balancer.
- Please provide yaml manifests with all properties that seem necessary for a proper, secure deployment, with resources scaled as we feel necessary.
- If not nginx load balancer - what other options would one consider to allow traffic distributed across multiple replicas of our web applications?

#### 3) Infrastructure as Code (Cloud)

A simple, startup infrastructure configuration as a code for a cloud of choice (preferably GCP or AWS) launching the following:

- VPC with subnets.
- A small, 4 nodes Kubernetes cluster (again, GKE or EKS, depending on the choice of cloud) that would be able to accommodate the kubernetes deployment prepared in previous exercise.
- Provide the proper network security configuration allowing only necessary traffic to our service.
- The IaC code should ideally be Terraform - however, usage of other IaC languages of choice are accepted to see the reasoning.
- On our side, though, Terraform or Config Connector knowledge are fairly important, so the candidate must be anyway willing to switch to these solutions.

#### 4) Ansible Playbook

Write an Ansible playbook doing the following:

- Connect to the list of Linux servers where both RedHat and Ubuntu distributions can installed
- Gather facts about systems
- Update repositories for all systems
- Upgrade servers with the latest packages available for the system version
- Make sure that Apache webserver is present on Ubuntu servers only
- If system is Ubuntu and Apache is installed, launch a simple configuration showing html document with, again, “Hello world”
- Once the Apache’s configuration file has been updated - Apache process should be restarted/reloaded to pick up the changes in configuration
- Make sure that MariaDB is installed on RedHat servers only

