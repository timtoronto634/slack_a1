resource "aws_iam_role" "bedrock_call_role" {
  name = "bedrock_call_role"

  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role_policy.json

  tags = {
    project = "ai-a2"
  }
}

data "aws_iam_policy_document" "lambda_assume_role_policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role_policy" "bedrock_call_policy" {
  name = "bedrock_call_policy"
  role = aws_iam_role.bedrock_call_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "logs:CreateLogGroup*",
        ]
        Effect   = "Allow"
        Resource = "arn:aws:logs:ap-northeast-1:*"
      },
      {
        Action = [
          "logs:CreateLogStream",
          "logs:PutLogEvents",
        ]
        Effect   = "Allow"
        Resource = ["arn:aws:logs:ap-northeast-1:*"]
      },
      {
        Action = [
          "bedrock:InvokeModel"
        ]
        Effect   = "Allow"
        Resource = ["*"]
      }
    ]
  })
}
