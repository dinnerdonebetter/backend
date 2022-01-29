locals {
  public_url      = "api.prixfixe.dev"
  api_server_port = 8000
}

resource "aws_ecr_repository" "api_server" {
  name = "api_server"
  # do not set image_tag_mutability to "IMMUTABLE", or else we cannot use :latest tags.

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_security_group" "api_server" {
  name        = "dev_api"
  description = "HTTP traffic"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port        = 80
    to_port          = 80
    protocol         = "TCP"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  ingress {
    from_port        = 443
    to_port          = 443
    protocol         = "TCP"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  ingress {
    from_port        = local.api_server_port
    to_port          = local.api_server_port
    protocol         = "TCP"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }
}

resource "aws_cloudwatch_log_group" "api_server" {
  name              = "/ecs/api_server"
  retention_in_days = local.log_retention_period_in_days
}

resource "aws_cloudwatch_log_group" "api_sidecar" {
  name              = "/ecs/dev-api-telemetry-collector-sidecar"
  retention_in_days = local.log_retention_period_in_days
}

resource "aws_iam_role" "api_task_execution_role" {
  name               = "api-task-execution-role"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_execution_assume_role.json
}

resource "aws_iam_role" "api_task_role" {
  name = "api-task-role"

  assume_role_policy = data.aws_iam_policy_document.ecs_task_assume_role.json

  inline_policy {
    name   = "allow_sqs_queue_access"
    policy = data.aws_iam_policy_document.allow_to_manipulate_queues.json
  }

  inline_policy {
    name   = "allow_ssm_access"
    policy = data.aws_iam_policy_document.allow_parameter_store_access.json
  }

  inline_policy {
    name   = "allow_decrypt_ssm_parameters"
    policy = data.aws_iam_policy_document.allow_to_decrypt_parameters.json
  }
}

resource "aws_ecs_task_definition" "dev_api" {
  family = "dev_api"

  container_definitions = jsonencode([
    {
      name : "otel-collector",
      image : format("%s:latest", aws_ecr_repository.otel_collector.repository_url)
      essential : true,
      logConfiguration : {
        logDriver : "awslogs",
        options : {
          awslogs-region : local.aws_region,
          awslogs-group : aws_cloudwatch_log_group.api_sidecar.name,
          awslogs-create-group : "true",
          awslogs-stream-prefix : "otel-collector"
        }
      }
    },
    #    {
    #      name  = "meal_plan_finalizer",
    #      image = format("%s:latest", aws_ecr_repository.meal_plan_finalizer.repository_url),
    #      essential : true,
    #      logConfiguration : {
    #        logDriver : "awslogs",
    #        options : {
    #          awslogs-region : local.aws_region,
    #          awslogs-group : aws_cloudwatch_log_group.meal_plan_finalizer.name,
    #          awslogs-stream-prefix : "ecs",
    #        },
    #      },
    #    },
    {
      name  = "api_server",
      image = format("%s:latest", aws_ecr_repository.api_server.repository_url),
      portMappings : [
        {
          containerPort : local.api_server_port,
          protocol : "tcp",
        },
      ],
      essential : true,
      logConfiguration : {
        logDriver : "awslogs",
        options : {
          awslogs-region : local.aws_region,
          awslogs-group : aws_cloudwatch_log_group.api_server.name,
          awslogs-stream-prefix : "ecs",
        },
      },
    },
  ])

  execution_role_arn = aws_iam_role.api_task_execution_role.arn
  task_role_arn      = aws_iam_role.api_task_role.arn

  # These are the minimum values for Fargate containers.
  cpu                      = 256
  memory                   = 512
  requires_compatibilities = ["FARGATE"]

  network_mode = "awsvpc"
}

resource "aws_ecs_service" "api_server" {
  name                               = "api_server"
  task_definition                    = aws_ecs_task_definition.dev_api.arn
  cluster                            = aws_ecs_cluster.dev.id
  launch_type                        = "FARGATE"
  deployment_maximum_percent         = 200
  deployment_minimum_healthy_percent = 100
  desired_count                      = 1

  deployment_controller {
    type = "ECS"
  }

  deployment_circuit_breaker {
    enable   = true
    rollback = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.api.arn
    container_name   = "api_server"
    container_port   = local.api_server_port
  }

  network_configuration {
    assign_public_ip = true

    security_groups = [
      aws_security_group.api_server.id,
    ]

    subnets = concat(
      [for x in aws_subnet.public_subnets : x.id],
      [for x in aws_subnet.private_subnets : x.id],
    )
  }

  depends_on = [
    aws_lb_listener.api_http,
  ]
}

resource "aws_acm_certificate" "api_dot" {
  domain_name       = local.public_url
  validation_method = "DNS"
}

resource "aws_security_group" "load_balancer" {
  name        = "load_balancer"
  description = "public internet traffic"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port        = 80
    to_port          = 80
    protocol         = "TCP"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  ingress {
    from_port        = 443
    to_port          = 443
    protocol         = "TCP"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  ingress {
    from_port        = local.api_server_port
    to_port          = local.api_server_port
    protocol         = "TCP"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }
}

resource "aws_lb" "api" {
  name               = "api-lb"
  internal           = false
  load_balancer_type = "application"

  subnets = [for x in aws_subnet.public_subnets : x.id]

  security_groups = [
    aws_security_group.load_balancer.id,
  ]

  depends_on = [aws_internet_gateway.main]
}

resource "aws_lb_target_group" "api" {
  name        = "api"
  port        = local.api_server_port
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = aws_vpc.main.id

  health_check {
    enabled  = true
    path     = "/_meta_/ready"
    port     = "traffic-port"
    matcher  = "200"
    protocol = "HTTP"
    timeout  = 15
  }

  depends_on = [aws_lb.api]
}

resource "aws_lb_listener" "api_http" {
  load_balancer_arn = aws_lb.api.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.api.arn
  }
}

resource "aws_lb_listener" "api_https" {
  load_balancer_arn = aws_lb.api.arn
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = aws_acm_certificate.api_dot.arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.api.arn
  }
}

resource "aws_lb_listener_certificate" "api_dot" {
  listener_arn    = aws_lb_listener.api_https.arn
  certificate_arn = aws_acm_certificate.api_dot.arn
}

resource "cloudflare_record" "api_dot_prixfixe_dot_dev" {
  zone_id         = var.CLOUDFLARE_ZONE_ID
  name            = local.public_url
  value           = aws_lb.api.dns_name
  type            = "CNAME"
  proxied         = true
  allow_overwrite = true
  ttl             = 1
}

resource "cloudflare_record" "api_dot_prixfixe_dot_dev_ssl_validation" {
  zone_id         = var.CLOUDFLARE_ZONE_ID
  name            = one(aws_acm_certificate.api_dot.domain_validation_options).resource_record_name
  value           = one(aws_acm_certificate.api_dot.domain_validation_options).resource_record_value
  type            = one(aws_acm_certificate.api_dot.domain_validation_options).resource_record_type
  proxied         = false
  allow_overwrite = true
  ttl             = 60
}
