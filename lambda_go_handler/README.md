# initialization
first, create a lambda function by following the setup.

https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html

At the last step, you should run,

```
aws lambda create-function --function-name myFunction \
--runtime provided.al2023 --handler bootstrap \
--architectures arm64 \
--role arn:aws:iam::111122223333:role/lambda-ex \
--zip-file fileb://myFunction.zip
```