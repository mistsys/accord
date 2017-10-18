## Using terraform for end to end

First setup aws-vault

```
[profile ssh-ca]
region=us-east-1

[profile ca-admin]
mfa_serial = arn:aws:iam::<account id>:mfa/pgautam
source_profile = ssh-ca
role_arn=arn:aws:iam::<account id>:role/Admin
```

terraform obviously requires the admin role so it can make changes.

## Planning and applying

```
aws-vault exec ssh-ca -- terraform plan
```

This will return what terraform is going to do

```
aws-vault exec ssh-ca -- terraform apply
```

This will execute the changes (diff'd from the plan and create a state file)

### Destroying

For everything in terraform, always run a plan, and this more importantly applies to destroy

```
aws-vault exec ssh-ca -- terraform plan -destroy
```

### Target destroying

Often you don't want to recreate the whole infrastructure, but destroy/recreate a particular resource, server. Terraform can let you target that individual server or resource.

```
aws-vault exec ssh-ca -- terraform plan -destroy -target aws_instance.cert_server
```

If the plan looks good, destroy all corresponding records. In our case, the primary record is protected from destory, so outright destroy wont work.

## Provisioning

This uses ansible playbooks to provision the servers, and expects future updates to go via ansible too. There are lot of benefits to ansible vs chef, which is supported out of box by Terraform.

1. all the state is saved directly in files, so versioning in git repo is straightforward
2. no chef server required for the setup
3. python -- this has a wider reach and while possible to monkey patch code, the community has less tendency to do so
4. we can build and maintain servers with very simple scripts and not rely on community cookbooks, playbooks should only be added after thorough review and testing


## Using Make targets

1. Creating servers

```
make plan
```

This makes a .plan file that's used for the next step

```
make apply
```

This will only bootstrap the servers, ie it's not going to make the servers functional yet

2. Reprovisioning servers

```
make reprovision
```

This will destroy the servers, and corresponding tcp elb + Route53 resource, and then recreate the servers again. This will lead to downtime of a few minutes as the new server is brought up. We can do something different in future but for now this will let us keep the servers upto date with latest patches.

3. Uploading the server and basic configuration

```
make upload
```

This runs the `caserver-upload.yml` on the servers that were created and upload the latest files, restart the service, etc
