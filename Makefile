APP_NAME=poststone-chat
STACK_NAME=poststone-chat
TEMPLATE=template.yaml
BUILD_DIR=.aws-sam/build
FUNCTION_NAME=PostStoneChatFunction

.PHONY: build local-invoke local-api logs deploy clean

bin:
	mkdir -p bin

build: bin
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bootstrap ./cmd/lambda

build-PostStoneChatFunction: build
	mv ./bootstrap $(ARTIFACTS_DIR)

clean:
	rm -r bin/