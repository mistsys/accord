#!/bin/bash
CURDIR="$( cd "$( dirname "$0" )" && pwd )"
AWS_VAULT=${AWS_VAULT:-"ca-admin"}
PREFIX="/caserver"
DAYS=30
AWS_REGION=${AWS_REGION:-"us-east-1"}
DRY_RUN=${DRY_RUN:-1}
START_TIME=$(date --iso-8601=seconds)
END_TIME=$(date --iso-8601=seconds -d "-$DAYS days ago")
# Get the key id from the terraform setup
KEY_ID=$(cd $CURDIR/terraform; terraform output caserver_kms_keyid)
TYPE=${1}
ID=${2}
path=$CURDIR/terraform/playbooks/files/certs
bits=64
password=$(</dev/urandom tr -dc 'A-Za-z0-9!"#$%&'\''()*+,-./:;<=>?@[\]^_`{|}~' | head -c $bits  ; echo)
comment="{\"id\": $ID, \"valid_from\": \"$START_TIME\", \"valid_until\": \"$END_TIME\"}"
mkdir -p $path || echo "directory already exists: $path"
case $TYPE in
    user|host)
        keypath=$path/ca_${TYPE}_${ID}
        ssh-keygen -t rsa -b 4096 -P $password -f $keypath -C "$comment"
        if [ "$DRY_RUN" = "1" ]; then
            echo "Dry Run: Not uploaded to param store"
        else
            if [ "x$KEY_ID" = "x" ]; then
                echo "No Key ID, setup environment with terraform"
                exit 1
            fi
            aws-vault exec ${AWS_VAULT} -- aws ssm put-parameter --name "$PREFIX/$ID" --value "$password" --type SecureString --region $AWS_REGION --key-id ${KEY_ID}
        fi
        echo "Encrypted passphrase for $keypath is in \"$PREFIX/$ID\", it wont be shown again. Do not save it anywhere other than ssm"
        ;;
    *)
        echo "Unknown cert type"
        exit 1
esac
