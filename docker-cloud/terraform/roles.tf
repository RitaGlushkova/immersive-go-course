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

resource "aws_iam_policy_attachment" "GitHubActionECRPublicPushImage" {
  roles      = aws_iam_role.GitHubActionECRPublicPushImage.name
  policy_arn = aws_iam_policy.GetAuthorizationToken.arn
}

resource "aws_iam_group_policy_attachment" "AllowPush" {
  group      = aws_iam_group.AllowPush.name
  policy_arn = aws_iam_policy.AllowPush.arn
}