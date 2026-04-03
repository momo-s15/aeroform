output "public_ip" {
  description = "Elastic IP address of the bot server"
  value       = aws_eip.bot.public_ip
}

output "ssh_command" {
  description = "Connect to your bot server"
  value       = "ssh ec2-user@${aws_eip.bot.public_ip}"
}

output "instance_id" {
  description = "EC2 instance ID"
  value       = aws_instance.bot.id
}

output "deploy_hint" {
  description = "Upload your bot code and restart the service"
  value       = "scp -r ./bot/* ec2-user@${aws_eip.bot.public_ip}:~/bot/ && ssh ec2-user@${aws_eip.bot.public_ip} 'sudo systemctl restart ${var.project_name}'"
}
