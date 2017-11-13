#!/bin/bash
CURDIR="$( cd "$( dirname "$0" )" && pwd )"
function check_requirements(){
    COUNTER=0
    for cmd in $(echo $1); do
        if [ -x "$(command -v $cmd)" ]; then
           echo "$cmd exists";
        else
            echo "$cmd doesn't exist";
            (( COUNTER +=1 ))
        fi
    done

    if [ "$COUNTER" -gt 0 ]; then
        echo "Requirements not met. Please install them before continuing"
        exit 1
    fi
}


echo "This script will help you setup Accord on your own AWS account. A fresh AWS account is recommended"

requirements="terraform ansible-playbook aws-vault make"
check_requirements "$requirements"

echo "====================="
echo "What do you want to call your organization (this will name the file, so no spaces)"
read org

echo "What is your email address? We'll configure SNS to send alerts here."
read email

echo "AWS Account ID?"
read aws_account_id

echo "AWS Account IDs to avoid running terraform on?"
read avoid_aws_accounts

echo "Domain to run the server on"
read domain

echo "hostname to run the server on. So this is something like just the accord part of accord.my.domain.net"
read ca_server

echo "Now we'll setup Google Client"

echo "Google Client ID"
read google_client_id

echo "Google Client Secret"
read google_client_secret

echo "Deployment aws-vault name to use"
read vault_name

echo "Deployment s3 url for the client binaries"
read deployment_s3_url


cat > $CURDIR/${org}.make <<EOF
GOOGLE_CLIENT_ID := ${google_client_id}
GOOGLE_CLIENT_SECRET := ${google_client_secret}
DEFAULT_SERVER := ${ca_server}.${domain}:443
TERRAFORM_VARS := -var 'account_id=${aws_account_id}' -var 'forbidden_account_ids=["${avoid_aws_accounts}"]' -var 'ca_domain=${domain}' -var 'ca_host=${ca_server}' -var 'primary_email=${email}'
DEPLOYMENT_VAULT := ${vault_name}
DEPLOYMENT_S3_URL := ${deployment_s3_url}
EOF

cat > $CURDIR/terraform/${org}.make <<EOF
include ../${org}.make
EOF
