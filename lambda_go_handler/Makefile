FUNCTION_NAME=slackA1
ZIP_FILE_NAME=slackA1.zip
MAX_SIZE=50000000  # 50MB in bytes

build:
	CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -tags lambda.norpc -o bootstrap .

zip: build
	zip $(ZIP_FILE_NAME) bootstrap

check:
	@file_size=$$(ls -l $(ZIP_FILE_NAME) | awk '{print $$5}'); \
	if [ $$file_size -gt $(MAX_SIZE) ]; then \
		echo "error: $(ZIP_FILE_NAME) is over 5MB ($$file_size)."; \
		exit 1; \
	fi

deploy: zip check
	aws lambda update-function-code --function-name $(FUNCTION_NAME) \
	--zip-file fileb://$(ZIP_FILE_NAME)
