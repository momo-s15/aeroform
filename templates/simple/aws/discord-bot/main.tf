terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

provider "aws" {
  region = var.region
}

# --- Look up the latest Amazon Linux 2023 AMI ---

data "aws_ami" "al2023" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-2023.*-x86_64"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

# --- Security group: SSH only ---

resource "aws_security_group" "bot" {
  name        = "${var.project_name}-bot-sg"
  description = "SSH access only for bot management"

  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "Allow all outbound"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name      = "${var.project_name}-bot-sg"
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}

# --- EC2 instance (t3.micro — free tier eligible for 12 months) ---

resource "aws_instance" "bot" {
  ami                    = data.aws_ami.al2023.id
  instance_type          = "t3.micro"
  vpc_security_group_ids = [aws_security_group.bot.id]
  key_name               = var.key_pair_name != "" ? var.key_pair_name : null

  root_block_device {
    volume_size = 8
    volume_type = "gp3"
    encrypted   = true
  }

  metadata_options {
    http_tokens = "required"
  }

  user_data = <<-SCRIPT
    #!/bin/bash
    set -e
    dnf update -y
    dnf install -y nodejs20 git

    # Create a systemd service so the bot auto-starts on reboot
    cat > /etc/systemd/system/${var.project_name}.service <<'SVC'
    [Unit]
    Description=${var.project_name} discord bot
    After=network.target

    [Service]
    Type=simple
    User=ec2-user
    WorkingDirectory=/home/ec2-user/bot
    ExecStart=/usr/bin/node /home/ec2-user/bot/index.js
    Restart=on-failure
    RestartSec=5

    [Install]
    WantedBy=multi-user.target
    SVC

    mkdir -p /home/ec2-user/bot
    echo 'console.log("Replace this with your bot code");' > /home/ec2-user/bot/index.js
    chown -R ec2-user:ec2-user /home/ec2-user/bot
    systemctl daemon-reload
    systemctl enable ${var.project_name}.service
  SCRIPT

  tags = {
    Name      = "${var.project_name}-bot"
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}

# --- Elastic IP for stable address ---

resource "aws_eip" "bot" {
  instance = aws_instance.bot.id
  domain   = "vpc"

  tags = {
    Name      = "${var.project_name}-bot-eip"
    ManagedBy = "aeroform"
  }
}
