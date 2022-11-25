resource "aws_lb_target_group" "docker_cloud" {
  name     = "docker-cloud-target-group-${var.username}-1"
  port     = 80
  protocol = "HTTP"
  vpc_id   = data.aws_vpc.default.id

  target_type = "ip"
}

resource "aws_lb" "docker_cloud" {
  name               = "docker-cloud-${var.username}-LB-1"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.lb_sg.id]
  subnets            = data.aws_subnets.public.ids

}

resource "aws_lb_listener" "docker_cloud" {
  load_balancer_arn = aws_lb.docker_cloud.arn
  port              = "80"
  protocol          = "HTTP"

    default_action {
        type             = "forward"
        target_group_arn = aws_lb_target_group.docker_cloud.arn
    }
}

  