.PHONY: all build clean release

VERSION=latest

default: lambda.zip

clean:
	rm lambda.zip main

lambda.zip: index.js main
	zip -r lambda.zip main index.js

main: *.go
	GOOS=linux go build -o main

release: lambda.zip
	aws s3 cp lambda.zip s3://convox/lambda/syslog.zip  --acl public-read
	for region in ap-northeast-1 ap-southeast-1 ap-southeast-2 eu-central-1 \
		eu-west-1 eu-west-2 us-east-1 us-east-2 us-west-1 us-west-2; do \
		aws s3 cp lambda.zip s3://convox-$$region/lambda/syslog.zip --acl public-read; \
	done
