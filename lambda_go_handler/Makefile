FUNCTION_NAME=myFunction
ZIP_FILE_NAME=myFunction.zip

build:
	GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap main.go

zip: build
	zip $(ZIP_FILE_NAME) bootstrap

deploy: zip
	aws lambda update-function-code --function-name $(FUNCTION_NAME) \
--zip-file fileb://$(ZIP_FILE_NAME)
