variable "myip" {}
variable "account_id" {}
variable "forbidden_account_ids" {type="list"}
variable "ca_domain" {}
variable "ca_host" {}
// where you want the notification emails to go
variable "primary_email" {}
data "aws_caller_identity" "current" {}
// if you want to deploy the thing all over in a new region pass
// pass -var region=<new region> to terraform plan
variable "region" { default = "us-east-1"}

// This file is in one place as I figure out an end-to-end way to do this
// in a secure and understandable way.. break it out further

provider "aws" {
    region = "${var.region}"
    assume_role {
       role_arn = "arn:aws:iam::${var.account_id}:role/Admin"
    }
    # Don't muck with the production account even accidentally
    forbidden_account_ids = "${var.forbidden_account_ids}"
}



data "aws_ami" "debian_stretch" {
   most_recent = true
   filter {
     name = "name"
     values = ["debian-stretch-hvm-x86_64-gp2*"]
   }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = [ "379101102735" ]
}


resource "aws_security_group" "egress_internet"{
  name = "egress_internet"
  description = "Allow egress to internet -- for package updates, etc."
  egress {
    from_port = 0
    to_port = 0
    # So that pings work too
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}



// TODO: figure out some way to make tasks that allows removing ingress IPs
resource "aws_security_group" "office_ssh" {
  name = "office_ssh"
  description = "Allow SSH From Office"

  // Your ip -- comment out when not running from home
  // this is populated by the make targets usually
  ingress {
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = [ "${var.myip}/32" ]
  }

  tags {
    Name = "allow_22"
  }
}

resource "aws_security_group" "cert_server" {
  name        = "cert_server"
  description = "Only allow 443 and 80"

  // TODO: once this is well tested all 443 traffic should go through ELB too
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    security_groups = ["${aws_security_group.cert_elb.id}"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    security_groups = ["${aws_security_group.cert_elb.id}"]
  }

  // health check port
  ingress {
    from_port   = 9110
    to_port     = 9110
    protocol    = "tcp"
    security_groups = ["${aws_security_group.cert_elb.id}"]
  }

  tags {
    Name = "allow_443"
  }
}

resource "aws_security_group" "cert_elb" {
  name = "cert_server_elb"
  description = "Allow elb to access"
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port = 0
    to_port = 0
    # So that pings work too
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name = "allow_443"
  }
}


/* What's going on:
- We don't want people or services to login and be able to view secrets without leaving audit trail
- AWS provides parameter store to manage that
- Create a caserver key with KMS and only allow one role to read it
- The instance has to be able to assume this new IAM role that can read the secrets
- but we don't want to give the instance that right forever
- so the service needs to know what ARN it should assume, and use that to read the passwords for
- the encrypted root keys
*/
resource "aws_kms_key" "caserver_key" {
  lifecycle {
    // it's kind of a bad idea to destroy this
    prevent_destroy = true
  }

}

resource "aws_kms_alias" "caserver_key" {
  name = "alias/caserver_key"
  target_key_id = "${aws_kms_key.caserver_key.key_id}"
}

// this is a sanity check
resource "aws_ssm_parameter" "caserver_params" {
  name  = "caserverTestString"
  type  = "SecureString"
  value = "Should be some kind of nontrivial value"
  key_id = "${aws_kms_key.caserver_key.key_id}"
}


# See: https://www.terraform.io/docs/providers/aws/d/iam_policy_document.html
data "aws_iam_policy_document" "caserver_readkeys_policy" {
  statement {
    sid = "AbilityToReadParams"
    actions = [
      "ssm:DescribeParameters"
    ]
    //effect = "Allow"
    # this is perhaps too permissive, although the action is limited to ssm
    resources = ["*"]
  }
  statement {
    sid = "DecryptTheKeys"
    actions = [
      "kms:Decrypt"
    ]
    //effect = "Allow"
    resources = [
      "${aws_kms_key.caserver_key.arn}"
    ]
  }
  //
  statement {
    sid = "GetAllCAServerParams"
    actions = [
      "ssm:GetParameters"
    ]
    resources = [
      "arn:aws:ssm:${var.region}:${var.account_id}:parameter/caserver*"
    ]
  }
}
// this is probably overkill.. think it's better to just read json and interpolate ids
resource "aws_iam_policy" "caserver_readkeys_policy" {
  name = "caserver_readkeys_policy"
  path = "/"
  policy = "${data.aws_iam_policy_document.caserver_readkeys_policy.json}"
}

// If we want people to be able to read the role -- need to make explicit trust
// to the user's ARN
// TODO: remove the access for current user once testing is complete
resource "aws_iam_role" "caserver_readkeys_role" {
  name = "caserver_readkeys"
  // Add trust so that ec2 servers can assume the role
  assume_role_policy =<<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": "1"
    },
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "AWS": "arn:aws:iam::${var.account_id}:user/pgautam"
      },
      "Effect": "Allow",
      "Sid": "2"
    },
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "AWS": "${aws_iam_role.caserver_role.arn}"
      },
      "Effect": "Allow",
      "Sid": "3"
    }
  ]
}
EOF
}

resource "aws_iam_policy_attachment" "caserver_readkeys"{
  name = "caserver_readkeys"
  roles = ["${aws_iam_role.caserver_readkeys_role.name}"]
  policy_arn = "${aws_iam_policy.caserver_readkeys_policy.arn}"
}

// Now make a role for our Cert server and allow it to assume the caserver_readkeys role
resource "aws_iam_role" "caserver_role" {
  name = "caserver"
  assume_role_policy =<<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": "1"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "caserver_policy" {
  name = "caserver_policy"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Resource": "${aws_iam_role.caserver_readkeys_role.arn}",
      "Effect": "Allow",
      "Sid": "1"
    }
  ]
}
EOF
}

resource "aws_iam_policy_attachment" "caserver_policy_attach" {
  name = "caserver_policy_attachment"
  roles = ["${aws_iam_role.caserver_role.name}"]
  policy_arn = "${aws_iam_policy.caserver_policy.arn}"
}

// Now make the instance profile to start servers with
// Understanding IAM profiles: https://github.com/hashicorp/terraform/issues/3851#issuecomment-171444541
resource "aws_iam_instance_profile" "caserver_profile" {
  name  = "caserver_profile"
  roles = ["${aws_iam_role.caserver_role.name}"]
}

# Need to make more than one machine and set them up behind the ELB
resource "aws_instance" "cert_server" {
  ami = "${data.aws_ami.debian_stretch.id}"
  instance_type = "t2.micro"

  # This is not on VPC yet, if it were it would have to be vpc_security_group_ids
  security_groups = [ "${aws_security_group.office_ssh.name}",
    "${aws_security_group.cert_server.name}",
    "${aws_security_group.egress_internet.name}"
  ]

  # ca-setup is already expected to have been generated and saved safely with
  #  aws-vault exec ca-admin -- aws ec2 create-key-pair --key-name ca-setup
  # it probably can be automated but I don't see a huge value there
  key_name = "ca-setup"

  # This is where we configure the instance with ansible-playbook
  # expect the key to be already loaded in ssh-agent before running
  provisioner "local-exec" {
    # Keep the remote files if we need to debug
    # Does it reasonably verbosely
    command = "sleep 120; ANSIBLE_HOST_KEY_CHECKING=False ANSIBLE_KEEP_REMOTE_FILES=True ansible-playbook -vv -u admin -i '${self.public_ip},' playbooks/caserver.yml"
  }

  iam_instance_profile = "${aws_iam_instance_profile.caserver_profile.id}"
}

resource "aws_sns_topic" "ssh_ca_cloudwatch_notifications" {
  name = "ssh_ca_cloudwatch_notifications"
  // shelling out because terraform model requires arn to be returned and email subscriber doesn't
  // until it validates externally
  provisioner "local-exec" {
    command = "aws sns subscribe --topic-arn ${self.arn} --region us-east-1 --protocol email --notification-endpoint ${var.primary_email}"
  }
}



resource "aws_cloudwatch_metric_alarm" "high_cpu_utilization" {
  alarm_name = "CAServerHighCPU"
  evaluation_periods = "2"
  statistic = "Average"
  metric_name = "CPUUtilization"
  namespace = "CA"
  dimensions = {
    InstanceId = "${aws_instance.cert_server.id}"
  }
  period = "300"
  threshold = 90
  unit = "Percent"
  comparison_operator = "GreaterThanThreshold"

  alarm_actions = ["${aws_sns_topic.ssh_ca_cloudwatch_notifications.arn}"]
}


resource "aws_elb" "cert_server_lb" {
  name = "cert-server-elb"
  availability_zones = [ "us-east-1a", "us-east-1b", "us-east-1c", "us-east-1d"]

  listener {
    instance_port = 443
    instance_protocol = "tcp"
    lb_port = 443
    lb_protocol = "tcp"
  }

  listener {
    instance_port = 80
    instance_protocol = "http"
    lb_port = 80
    lb_protocol = "http"
  }

  security_groups = ["${aws_security_group.cert_elb.id}"]

  health_check {
    healthy_threshold = 2
    unhealthy_threshold = 2
    timeout = 3
    target = "HTTP:9110/about"
    interval = 10
  }

  instances                   = ["${aws_instance.cert_server.id}"]
  cross_zone_load_balancing   = true
  idle_timeout                = 400

  tags {
    Type = "tcp_elb"
  }
}

resource "aws_route53_zone" "primary" {
  name = "${var.ca_domain}"
  lifecycle{
    // The name-servers are set on the main record in official account
    // so destroying requires setting it up all over again
    prevent_destroy = true
  }
}

// this makes certserver.domain
resource "aws_route53_record" "certserver" {
  zone_id = "${aws_route53_zone.primary.zone_id}"
  name    = "${var.ca_host}"
  type    = "A"
  // TTL is 60 by default for alias

  alias {
    name                   = "${aws_elb.cert_server_lb.dns_name}"
    zone_id                = "${aws_elb.cert_server_lb.zone_id}"
    evaluate_target_health = true
  }
}



output "caserver_ips" {
  value = "${aws_instance.cert_server.*.public_ip}"
}

output "caserver_readkeys_arn" {
  value = "${aws_iam_role.caserver_readkeys_role.arn}"
}

// these are the parameters where the secure strings for the root and user CA keys are
output "caserver_params_prefix" {
  value = "/caserver"
}

output "caserver_kms_keyid" {
  value = "${aws_kms_key.caserver_key.arn}"
}

output "caserver_hostname" {
  value = "${var.ca_host}.${var.ca_domain}"
}
