# initialization
first, create a lambda function by following the setup.

https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html

At the last step, you should run something like below

```
aws lambda create-function --function-name slackA1 \
--runtime provided.al2023 --handler bootstrap \
--architectures arm64 \
--role arn:aws:iam::111122223333:role/lambda-ex \
--zip-file fileb://slackA1.zip
--timeout 15
```

In addition to the default setting, following configuration is required

- configure Function URL
- add environment variable
- setting timeout longer (if needed, but 3 second will not be enough in most case)