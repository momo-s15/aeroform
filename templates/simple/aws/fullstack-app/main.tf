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

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

# --- S3 bucket for static assets ---

resource "aws_s3_bucket" "assets" {
  bucket        = "${var.project_name}-assets"
  force_destroy = true

  tags = {
    Name      = "${var.project_name}-assets"
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}

resource "aws_s3_bucket_public_access_block" "assets" {
  bucket = aws_s3_bucket.assets.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_server_side_encryption_configuration" "assets" {
  bucket = aws_s3_bucket.assets.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

# --- Default VPC lookup for the database ---

data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

# --- Security group: PostgreSQL from VPC only ---

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

# --- RDS PostgreSQL (private, encrypted, backed up) ---

resource "aws_db_subnet_group" "db" {
  name       = "${var.project_name}-db-subnet"
  subnet_ids = data.aws_subnets.default.ids

  tags = {
    Name      = "${var.project_name}-db-subnet"
    ManagedBy = "aeroform"
  }
}

resource "aws_db_instance" "db" {
  identifier     = "${var.project_name}-db"
  engine         = "postgres"
  engine_version = "16.4"
  instance_class = "db.t3.micro"

  allocated_storage     = 20
  max_allocated_storage = 50
  storage_type          = "gp3"
  storage_encrypted     = true

  db_name  = "app"
  username = "aeroform"
  password = var.db_password

  db_subnet_group_name   = aws_db_subnet_group.db.name
  vpc_security_group_ids = [aws_security_group.db.id]
  publicly_accessible    = false

  backup_retention_period = 7
  skip_final_snapshot     = true
  deletion_protection     = false

  auto_minor_version_upgrade = true

  tags = {
    Name      = "${var.project_name}-db"
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}

# --- App Runner service ---

resource "aws_iam_role" "apprunner" {
  name = "${var.project_name}-apprunner-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect    = "Allow"
        Principal = { Service = "tasks.apprunner.amazonaws.com" }
        Action    = "sts:AssumeRole"
      }
    ]
  })

  tags = {
    Name      = "${var.project_name}-apprunner-role"
    ManagedBy = "aeroform"
  }
}

resource "aws_iam_role_policy" "apprunner" {
  name = "${var.project_name}-apprunner-policy"
  role = aws_iam_role.apprunner.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = ["s3:GetObject", "s3:PutObject", "s3:ListBucket"]
        Resource = [aws_s3_bucket.assets.arn, "${aws_s3_bucket.assets.arn}/*"]
      }
    ]
  })
}

resource "aws_apprunner_service" "app" {
  service_name = "${var.project_name}-app"

  source_configuration {
    image_repository {
      image_identifier      = var.container_image
      image_repository_type = "ECR_PUBLIC"

      image_configuration {
        port = var.app_port

        runtime_environment_variables = {
          DATABASE_URL = "postgresql://aeroform:${var.db_password}@${aws_db_instance.db.endpoint}/app"
          S3_BUCKET    = aws_s3_bucket.assets.id
          AWS_REGION   = data.aws_region.current.name
        }
      }
    }

    auto_deployments_enabled = false
  }

  instance_configuration {
    cpu               = "1024"
    memory            = "2048"
    instance_role_arn = aws_iam_role.apprunner.arn
  }

  health_check_configuration {
    protocol            = "HTTP"
    path                = "/health"
    healthy_threshold   = 2
    unhealthy_threshold = 3
    interval            = 10
  }

  tags = {
    Name      = "${var.project_name}-app"
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}
