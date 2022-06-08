# Learning Go and Pulumi for the first time.  

Set up a simple AWS infrastructure that is composed of:
* Network
    * VPC
    * Private/Public Subnets
* Kubernetes
    * Autoscaling Group
    * Spot Instances
    * Application(?) Load Balancers
    * Kubernetes Pods
    * Hardening(?) prob too much effort
    * NGINX maybe?
* IAM
    * Service roles
    * Instance roles with SSM Session Manager access
* IaC Unit tests - Never done this before but might be good practice

Since Pulumi is pretty bad at handling AWS assume-role profiles we gotta write something to support this separately.\
Package 'config' has to run first to configure AWS credentials with MFA.
It will cache the credentials for X minutes.
