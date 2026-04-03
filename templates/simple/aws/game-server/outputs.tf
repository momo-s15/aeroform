output "public_ip" {
  description = "Elastic IP address of the game server"
  value       = aws_eip.game.public_ip
}

output "game_address" {
  description = "Connect to your game server at this address"
  value       = "${aws_eip.game.public_ip}:${var.game_port}"
}

output "ssh_command" {
  description = "Connect to manage your server"
  value       = "ssh ec2-user@${aws_eip.game.public_ip}"
}

output "instance_id" {
  description = "EC2 instance ID"
  value       = aws_instance.game.id
}
