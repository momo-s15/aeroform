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

# --- Look up the default VPC and its subnets ---

data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

# --- Security group: PostgreSQL only from within the VPC ---

resource "aws_security_group" "db" {
  name        = "${var.project_name}-db-sg"
  description = "Allow PostgreSQL from within the VPC only"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    description = "PostgreSQL from VPC"
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [data.aws_vpc.default.cidr_block]
  }

  egress {
    description = "Allow all outbound"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name      = "${var.project_name}-db-sg"
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}

# --- DB subnet group using default VPC subnets ---

resource "aws_db_subnet_group" "db" {
  name       = "${var.project_name}-db-subnet"
  subnet_ids = data.aws_subnets.default.ids

  tags = {
    Name      = "${var.project_name}-db-subnet"
    ManagedBy = "aeroform"
  }
}

# --- RDS PostgreSQL instance (private, encrypted, backed up) ---

resource "aws_db_instance" "db" {
  identifier     = "${var.project_name}-db"
  engine         = "postgres"
  engine_version = "16.4"
  instance_class = "db.t3.micro"

  allocated_storage     = 20
  max_allocated_storage = 50
  storage_type          = "gp3"
  storage_encrypted     = true

  db_name  = var.db_name
  username = var.db_username
  password = var.db_password

  db_subnet_group_name   = aws_db_subnet_group.db.name
  vpc_security_group_ids = [aws_security_group.db.id]
  publicly_accessible    = false

  backup_retention_period = 7
  backup_window           = "03:00-04:00"
  maintenance_window      = "sun:05:00-sun:06:00"

  skip_final_snapshot       = true
  final_snapshot_identifier = "${var.project_name}-final-snapshot"
  deletion_protection       = false

  auto_minor_version_upgrade = true

  tags = {
    Name      = "${var.project_name}-db"
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}
