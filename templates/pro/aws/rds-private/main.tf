terraform {
  required_version = ">= 1.0"
}

resource "aws_db_subnet_group" "this" {
  name       = var.name
  subnet_ids = var.subnet_ids
}

resource "aws_db_instance" "this" {
  identifier              = var.name
  instance_class          = var.instance_class
  engine                  = var.engine
  allocated_storage       = var.allocated_storage
  db_subnet_group_name    = aws_db_subnet_group.this.name
  publicly_accessible     = false
  multi_az                = true
  storage_encrypted       = true
  backup_retention_period = 7
}
