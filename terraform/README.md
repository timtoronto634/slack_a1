# Infrastructure

Currently, it just creates an IAM role for lambda function, run

```
terraform apply
```

Lambda function is NOT controlled under terraform, because it needs frequent update in development process.

For creating a lambda function, see ../lambda_go_handler/README.md