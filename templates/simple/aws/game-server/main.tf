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

# --- Security group: game port + SSH ---

resource "aws_security_group" "game" {
  name        = "${var.project_name}-game-sg"
  description = "Game port and SSH access"

  ingress {
    description = "Game port (TCP)"
    from_port   = var.game_port
    to_port     = var.game_port
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "Game port (UDP)"
    from_port   = var.game_port
    to_port     = var.game_port
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0"]
  }

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
    Name      = "${var.project_name}-game-sg"
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}

# --- EC2 instance (t3.medium for game workloads) ---

resource "aws_instance" "game" {
  ami                    = data.aws_ami.al2023.id
  instance_type          = "t3.medium"
  vpc_security_group_ids = [aws_security_group.game.id]
  key_name               = var.key_pair_name != "" ? var.key_pair_name : null

  root_block_device {
    volume_size = 30
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
    dnf install -y java-21-amazon-corretto-headless screen
    echo "Server ready — upload your game server files to /home/ec2-user/server/"
    mkdir -p /home/ec2-user/server
    chown -R ec2-user:ec2-user /home/ec2-user/server
  SCRIPT

  tags = {
    Name      = "${var.project_name}-game"
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}

# --- Elastic IP for stable address ---

resource "aws_eip" "game" {
  instance = aws_instance.game.id
  domain   = "vpc"

  tags = {
    Name      = "${var.project_name}-game-eip"
    ManagedBy = "aeroform"
  }
}
