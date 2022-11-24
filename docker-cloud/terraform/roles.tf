resource "aws_iam_role" "GitHubActionECRPublicPushImage" {
  name = "GitHubActionECRPublicPushImage"
  assume_role_policy = jsonencode({
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Federated": "arn:aws:iam::297880250375:oidc-provider/token.actions.githubusercontent.com"
            },
            "Action": "sts:AssumeRoleWithWebIdentity",
            "Condition": {
                "StringEquals": {
                    "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"
                },
                "StringLike": {
                    "token.actions.githubusercontent.com:sub": "repo:RitaGlushkova/immersive-go-course:*"
                }
            }
        }
    ]
})
}

resource "aws_iam_policy_attachment" "GetAuthorizationToken" {
  roles      = aws_iam_role.GitHubActionECRPublicPushImage.name
  policy_arn = aws_iam_policy.GetAuthorizationToken.arn
}

resource "aws_iam_policy_attachment" "AllowPush" {
  roles      = aws_iam_role.GitHubActionECRPublicPushImage.name
  policy_arn = aws_iam_policy.AllowPush.arn
}
