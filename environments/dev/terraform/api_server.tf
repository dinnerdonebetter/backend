locals {
  public_url = "api.prixfixe.dev"
}

resource "aws_ecr_repository" "api_server" {
  name = "api_server"
  # do not set image_tag_mutability to "IMMUTABLE", or else we cannot use :latest tags.

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_security_group" "api_service" {
  name        = "prixfixe_api"
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
    from_port        = 8000
    to_port          = 8000
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
    from_port        = 8000
    to_port          = 8000
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

resource "aws_cloudwatch_log_group" "api_server_firelens_logs" {
  name              = "/ecs/api_server/firelens"
  retention_in_days = local.log_retention_period_in_days
}

resource "aws_ecs_task_definition" "api_server" {
  family = "api_server"

  container_definitions = jsonencode([
    {
      essential : true,
      image : "grafana/fluent-bit-plugin-loki:2.0.0-amd64",
      name : "log_router",
      firelensConfiguration : {
        type : "fluentbit",
        options : {
          enable-ecs-log-metadata : "true"
        }
      },
      logConfiguration : {
        logDriver : "awslogs",
        options : {
          awslogs-group : aws_cloudwatch_log_group.api_server_firelens_logs.name,
          awslogs-region : "us-east-1",
          awslogs-create-group : "true",
          awslogs-stream-prefix : "firelens"
        }
      },
      memoryReservation : 50
    },
    {
      name  = "api_server",
      image = format("%s:latest", aws_ecr_repository.api_server.repository_url),
      portMappings : [
        {
          containerPort : 8000,
          protocol : "tcp",
        },
      ],
      logConfiguration : {
        logDriver : "awsfirelens",
        options : {
          Name : "loki",
          Url : "https://${var.GRAFANA_CLOUD_USER_ID}:${var.GRAFANA_CLOUD_API_KEY}@logs-prod-us-central1.grafana.net/loki/api/v1/push",
          Labels : "{job=\"firelens\"}",
          RemoveKeys : "container_id,ecs_task_arn",
          LabelKeys : "container_name,ecs_task_definition,source,ecs_cluster",
          LineFormat : "key_value",
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

resource "aws_ecs_cluster" "api" {
  name = "api_servers"
}

resource "aws_ecs_service" "api_server" {
  name                               = "api_server"
  task_definition                    = aws_ecs_task_definition.api_server.arn
  cluster                            = aws_ecs_cluster.api.id
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
    container_port   = 8000
  }

  network_configuration {
    assign_public_ip = true

    security_groups = [
      aws_security_group.api_service.id,
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

data "aws_iam_policy_document" "ecs_task_execution_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "api_task_execution_role" {
  name               = "api-task-execution-role"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_execution_assume_role.json
}

# Normally we'd prefer not to hardcode an ARN in our Terraform, but since this is an AWS-managed policy, it's okay.
resource "aws_iam_role_policy_attachment" "ecs_task_execution_role" {
  role       = aws_iam_role.api_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

data "aws_iam_policy_document" "ecs_task_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
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

resource "aws_acm_certificate" "api_dot" {
  domain_name       = local.public_url
  validation_method = "DNS"
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
  port        = 8000
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
