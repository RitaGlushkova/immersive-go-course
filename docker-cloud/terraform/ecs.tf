# resource "aws_ecs_service" "docker_cloud" {
#   name            = "docker-cloud-${var.username}"
#   cluster         = aws_ecs_cluster.docker_cloud.id
#   task_definition = aws_ecs_task_definition.docker_cloud.arn
#   desired_count   = 1
#   iam_role        = aws_iam_role.GitHubActionECRPublicPushImage.name
#   depends_on      = [aws_iam_role_policy.Attachment]

#   ordered_placement_strategy {
#     type  = "binpack"
#     field = "cpu"
#   }

#   load_balancer {
#     target_group_arn = aws_lb_target_group.foo.arn
#     container_name   = "mongo"
#     container_port   = 8080
#   }

#   placement_constraints {
#     type       = "memberOf"
#     expression = "attribute:ecs.availability-zone in [us-west-2a, us-west-2b]"
#   }
# }