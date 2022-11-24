resource "aws_iam_policy" "GetAuthorizationToken" {
    name        = "GetAuthorizationToken-${var.username}"
    path       = "/"
    description = "Allows the user to get an authorization token for the ECR registry" 
    policy = jsonencode({
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "ecr-public:GetAuthorizationToken",
                "sts:GetServiceBearerToken"
            ],
            "Resource": "*"
        }
    ]
})
}

resource "aws_iam_policy" "AllowPush" {
  name = "AllowPush-${var.username}"
  path = "/"
    description = "Allows the user to push images to the ECR registry"
    policy = jsonencode({
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "ecr-public:InitiateLayerUpload",
                "ecr-public:UploadLayerPart",
                "ecr-public:PutImage",
                "ecr-public:CompleteLayerUpload",
                "ecr-public:BatchCheckLayerAvailability"
            ],
            "Resource": "arn:aws:ecr-public::297880250375:repository/immersive-go-course/docker-cloud-rita"
        }
    ]
})
}