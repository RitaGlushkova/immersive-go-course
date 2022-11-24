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