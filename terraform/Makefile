AWS_VAULT ?= ssh-ca
terraform ?= terraform
packer ?= packer
ansible-playbook ?= ansible-playbook
TOP ?= .
server_type ?= caserver

# Get your external IP to ssh into the servers
.myip:
	$(shell dig +short myip.opendns.com @resolver1.opendns.com > $(TOP)/.myip)


.plan: .myip terraform.tf .destroy_plan
	aws-vault exec $(AWS_VAULT) -- $(terraform) plan -var myip=$(shell cat $(TOP)/.myip) $(TERRAFORM_VARS) -out .plan

plan: .plan

list_servers:
	$(terraform) output $(server_type)_ips

bootstrap:
	ANSIBLE_HOST_KEY_CHECKING=False ANSIBLE_KEEP_REMOTE_FILES=True $(ansible-playbook) -u admin \
	-i '$(shell $(terraform) output $(server_type)_ips),' -vv $(TOP)/playbooks/$(server_type).yml

ansible-write-vars: terraform.tfstate
	echo "role_arn: $(shell $(terraform) output $(server_type)_readkeys_arn)" > $(TOP)/playbooks/vars/caserver.yml
	echo "params_prefix: $(shell $(terraform) output $(server_type)_params_prefix)" >> $(TOP)/playbooks/vars/caserver.yml
	echo "hostname: $(shell $(terraform) output $(server_type)_hostname)" >> $(TOP)/playbooks/vars/caserver.yml

upload: ansible-write-vars
	ANSIBLE_HOST_KEY_CHECKING=False ANSIBLE_KEEP_REMOTE_FILES=True $(ansible-playbook) -u admin \
	-i '$(shell $(terraform) output $(server_type)_ips),' -vv $(TOP)/playbooks/$(server_type)-upload.yml

apply: terraform.tf .myip .plan
	aws-vault exec $(AWS_VAULT) -- $(terraform) apply .plan

.destroy_plan:
	aws-vault exec $(AWS_VAULT) -- $(terraform) destroy -target aws_instance.cert_server -var myip=$(shell cat $(TOP)/.myip) $(TERRAFORM_VARS) -out .destroy_plan

terraform-plan.png: terraform.tf
	$(terraform) graph -draw-cycles -type=plan | dot -Tpng > terraform-plan.png

terraform.png: terraform.tf
	$(terraform) graph -draw-cycles | dot -Tpng > terraform.png

graph-plan: terraform-plan.png
	open terraform-plan.png

graph: terraform.png
	open terraform.png

destroy_servers: .destroy_plan
	aws-vault exec $(AWS_VAULT) -- $(terraform) destroy -target aws_instance.cert_server -var myip=$(shell cat $(TOP)/.myip) $(TERRAFORM_VARS)

reprovision: destroy_servers plan apply

packer-base: base-packer.json
	aws-vault exec $(AWS_VAULT) -- $(packer) build base-packer.json

PHONY: plan clean

clean:
	rm .plan .myip
